package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
)

func (r *repo) CheckStock(ctx context.Context, code string, amount int32) (bool, error) {
	isAvailable := true
	err := r.execTx(ctx, func(q *db.Queries) error {
		product, err := q.UpdateStock(ctx, db.UpdateStockParams{
			Reserved: pgtype.Int4{
				Int32: amount,
				Valid: true,
			},
			Code: code,
		})
		if err != nil {
			return err
		}

		if product.Reserved > product.Total {
			isAvailable = false
			return fmt.Errorf("out of reserved")
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
		product, err := q.UpdateStock(ctx, db.UpdateStockParams{
			Total: pgtype.Int4{
				Int32: -amount,
				Valid: false,
			},
			Reserved: pgtype.Int4{
				Int32: -amount,
				Valid: true,
			},
			Code: code,
		})
		if err != nil {
			return err
		}

		if product.Total < 0 {
			isAvailable = false
			return fmt.Errorf("out of stock")
		}

		if product.Reserved > product.Total {
			isAvailable = false
			return fmt.Errorf("out of reserved")
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
		product, err := q.UpdateStock(ctx, db.UpdateStockParams{
			Reserved: pgtype.Int4{
				Int32: amount,
				Valid: true,
			},
			Code: code,
		})
		if err != nil {
			return err
		}

		if product.Reserved > product.Total {
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
