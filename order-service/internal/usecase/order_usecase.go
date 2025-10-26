package usecase

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	db "github.com/ymanshur/synasishouse/order/db/sqlc"
	"github.com/ymanshur/synasishouse/order/internal/common"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/consts"
	"github.com/ymanshur/synasishouse/order/internal/presentation"
	"github.com/ymanshur/synasishouse/order/internal/repo"
	"github.com/ymanshur/synasishouse/order/internal/typex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Orderer interface {
	Create(ctx context.Context, req presentation.OrderRequest) (*presentation.OrderResponse, error)
	Settle(ctx context.Context, req presentation.SettleRequest) (*presentation.SettleResponse, error)
}

type orderUseCase struct {
	repo repo.Repo
	conn connector.Inventorier
}

func NewOrder(repo repo.Repo, conn connector.Inventorier) Orderer {
	return &orderUseCase{
		repo: repo,
		conn: conn,
	}
}

func (u *orderUseCase) Create(ctx context.Context, req presentation.OrderRequest) (*presentation.OrderResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	res := &presentation.OrderResponse{}
	err = u.repo.Transact(ctx, pgx.ReadCommitted, func(r repo.Repo) error {
		var err error

		newOrder, err := r.CreateOrder(ctx, db.CreateOrderParams{
			OrderNo: req.OrderNo,
			UserID:  userID,
			Status:  consts.OrderStatusPending,
		})
		if err != nil {
			if common.PGErrorCode(err) == consts.UniqueViolation {
				return typex.NewUnprocessableEntityError("order unique constraint violated")
			}
			return fmt.Errorf("create order: %w", err)
		}

		res.OrderNo = newOrder.OrderNo
		res.UserID = newOrder.UserID.String()
		res.Status = newOrder.Status

		for _, detail := range req.Details {
			newOrderDetail, err := r.CreateOrderDetail(ctx, db.CreateOrderDetailParams{
				OrderID:     newOrder.ID,
				ProductCode: detail.ProductCode,
				Amount:      detail.Amount,
			})
			if err != nil {
				return fmt.Errorf("create order detail: %w", err)
			}

			res.Details = append(res.Details, presentation.OrderDetailResponse{
				ProductCode: newOrderDetail.ProductCode,
				Amount:      newOrderDetail.Amount,
			})
		}

		params := connector.CheckStockParams{
			Stocks: []connector.StockParams{},
		}
		for _, detail := range req.Details {
			params.Stocks = append(params.Stocks, connector.StockParams{
				ProductCode: detail.ProductCode,
				Amount:      detail.Amount,
			})
		}

		isAvailable, err := u.conn.CheckStock(ctx, params)
		if err != nil {
			errRPC := status.Convert(err)
			if errRPC.Code() == codes.NotFound {
				return typex.NewNotFoundError("product")
			}

			return fmt.Errorf("check stock: %w", err)
		}

		if !isAvailable {
			return typex.NewUnprocessableEntityError("stocks is not available")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *orderUseCase) Settle(ctx context.Context, req presentation.SettleRequest) (*presentation.SettleResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	userOrder, err := u.repo.GetUserOrder(ctx, db.GetUserOrderParams{
		OrderNo: req.OrderNo,
		UserID:  userID,
	})
	if err != nil {
		if errors.Is(err, consts.ErrRecordNotFound) {
			return nil, typex.NewNotFoundError("user order")
		}
		return nil, err
	}

	if userOrder.Status == consts.OrderStatusSettled {
		return nil, typex.NewConflictError("user order already settled")
	}

	res := &presentation.SettleResponse{}
	err = u.repo.Transact(ctx, pgx.ReadCommitted, func(r repo.Repo) error {
		var err error

		settledOrder, err := r.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
			ID:     userOrder.ID,
			Status: consts.OrderStatusSettled,
		})
		if err != nil {
			return fmt.Errorf("update order status: %w", err)
		}

		res.OrderNo = settledOrder.OrderNo
		res.UserID = settledOrder.UserID.String()
		res.Status = settledOrder.Status

		orderDetails, err := u.repo.ListOrderDetails(ctx, userOrder.ID)
		if err != nil {
			return fmt.Errorf("list order details: %w", err)
		}

		params := connector.ReserveStockParams{
			Stocks: []connector.StockParams{},
		}
		for _, detail := range orderDetails {
			params.Stocks = append(params.Stocks, connector.StockParams{
				ProductCode: detail.ProductCode,
				Amount:      detail.Amount,
			})
		}

		isAvailable, err := u.conn.ReserveStock(ctx, params)
		if err != nil {
			errRPC := status.Convert(err)
			if errRPC.Code() == codes.NotFound {
				return typex.NewNotFoundError("product")
			}

			return fmt.Errorf("reserve stock: %w", err)
		}

		if !isAvailable {
			return fmt.Errorf("unmatch order stocks")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
