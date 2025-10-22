package usecase

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/presentation"
	"github.com/ymanshur/synasishouse/order/internal/typex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Orderer interface {
	Checkout(ctx context.Context, req presentation.OrderRequest) (*presentation.OrderResponse, error)
}

type orderUseCase struct {
	conn connector.Inventorier
}

func NewOrder(conn connector.Inventorier) Orderer {
	return &orderUseCase{conn: conn}
}

func (u *orderUseCase) Checkout(ctx context.Context, req presentation.OrderRequest) (*presentation.OrderResponse, error) {
	err := validation.Validate(&req)
	if err != nil {
		return nil, err
	}

	isAvailable, err := u.conn.CheckStock(ctx, connector.StockParams{
		Code:   req.Code,
		Amount: req.Amount,
	})
	if err != nil {
		errRPC := status.Convert(err)
		if errRPC.Code() == codes.NotFound {
			return nil, typex.NewNotFoundError("stock")
		}

		return nil, fmt.Errorf("check stock: %w", err)
	}

	res := &presentation.OrderResponse{IsAvailable: isAvailable}
	return res, nil
}
