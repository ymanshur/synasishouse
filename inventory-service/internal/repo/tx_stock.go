package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
)

func (r *repo) CheckStock(ctx context.Context, code string, amount int32) error {
	return r.execTx(ctx, func(q *db.Queries) error {
		_, err := q.UpdateStock(ctx, db.UpdateStockParams{
			Reserved: pgtype.Int4{
				Int32: amount,
				Valid: true,
			},
			Code: code,
		})
		return err
	})
}

func (r *repo) ReserveStock(ctx context.Context, code string, amount int32) error {
	return r.execTx(ctx, func(q *db.Queries) error {
		_, err := q.UpdateStock(ctx, db.UpdateStockParams{
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
		return err
	})
}

func (r *repo) ReleaseStock(ctx context.Context, code string, amount int32) error {
	return r.execTx(ctx, func(q *db.Queries) error {
		_, err := q.UpdateStock(ctx, db.UpdateStockParams{
			Reserved: pgtype.Int4{
				Int32: amount,
				Valid: true,
			},
			Code: code,
		})
		return err
	})
}
