package usecase

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/presentation"
)

type Orderer interface {
	Checkout(ctx context.Context, req presentation.OrderRequest) (bool, error)
}

type orderUseCase struct {
	conn connector.Inventorier
}

func NewOrder(conn connector.Inventorier) Orderer {
	return &orderUseCase{conn: conn}
}

func (u *orderUseCase) Checkout(ctx context.Context, req presentation.OrderRequest) (bool, error) {
	err := validation.Validate(&req)
	if err != nil {
		return false, err
	}

	ok, err := u.conn.CheckStock(ctx, connector.StockParams{
		Code:   req.Code,
		Amount: req.Amount,
	})
	if err != nil {
		return false, err
	}

	return ok, nil
}
