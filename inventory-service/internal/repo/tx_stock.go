package repo

import (
	"context"
	"fmt"

	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
)

type CreateStockParams struct {
	Stocks []StockParams
}

type StockParams struct {
	ProductCode string
	Amount      int32
}

func (r *repo) CheckStock(ctx context.Context, arg CreateStockParams) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		for _, stock := range arg.Stocks {
			product, err := q.AddStock(ctx, db.AddStockParams{
				Hold: stock.Amount,
				Code: stock.ProductCode,
			})
			if err != nil {
				return err
			}

			if product.Hold > product.Total {
				isAvailable = false
				return fmt.Errorf("out of hold")
			}
		}

		return nil
	})
	if err != nil && isAvailable {
		return false, err
	}

	return isAvailable, nil
}

func (r *repo) ReserveStock(ctx context.Context, arg CreateStockParams) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		for _, stock := range arg.Stocks {
			product, err := q.AddStock(ctx, db.AddStockParams{
				Total: -stock.Amount,
				Hold:  -stock.Amount,
				Code:  stock.ProductCode,
			})
			if err != nil {
				return err
			}

			if product.Total < 0 {
				isAvailable = false
				return fmt.Errorf("out of stock")
			}

			if product.Hold < 0 {
				isAvailable = false
				return fmt.Errorf("out of hold")
			}
		}

		return nil
	})
	if err != nil && isAvailable {
		return false, err
	}

	return isAvailable, nil
}

func (r *repo) ReleaseStock(ctx context.Context, arg CreateStockParams) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		for _, stock := range arg.Stocks {
			product, err := q.AddStock(ctx, db.AddStockParams{
				Hold: -stock.Amount,
				Code: stock.ProductCode,
			})
			if err != nil {
				return err
			}

			if product.Hold < 0 {
				isAvailable = false
				return fmt.Errorf("out of release")
			}
		}

		return nil
	})
	if err != nil && isAvailable {
		return false, err
	}

	return isAvailable, nil
}
