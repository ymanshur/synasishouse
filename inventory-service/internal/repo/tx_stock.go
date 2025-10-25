package repo

import (
	"context"
	"fmt"

	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
)

func (r *repo) CheckStock(ctx context.Context, code string, amount int32) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		product, err := q.AddStock(ctx, db.AddStockParams{
			Hold: amount,
			Code: code,
		})
		if err != nil {
			return err
		}

		if product.Hold > product.Total {
			isAvailable = false
			return fmt.Errorf("out of hold")
		}

		return nil
	})
	if err != nil && isAvailable {
		return false, err
	}

	return isAvailable, nil
}

func (r *repo) ReserveStock(ctx context.Context, code string, amount int32) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		product, err := q.AddStock(ctx, db.AddStockParams{
			Total: -amount,
			Hold:  -amount,
			Code:  code,
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

		return nil
	})
	if err != nil && isAvailable {
		return false, err
	}

	return isAvailable, nil
}

func (r *repo) ReleaseStock(ctx context.Context, code string, amount int32) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		product, err := q.AddStock(ctx, db.AddStockParams{
			Hold: -amount,
			Code: code,
		})
		if err != nil {
			return err
		}

		if product.Hold < 0 {
			isAvailable = false
			return fmt.Errorf("out of release")
		}

		return nil
	})
	if err != nil && isAvailable {
		return false, err
	}

	return isAvailable, nil
}
