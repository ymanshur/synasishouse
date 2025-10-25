package usecase

import (
	"context"
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
