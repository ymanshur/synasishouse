package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/ymanshur/synasishouse/inventory/db/sqlc"
)

// Repo defines all functions to execute db queries and transactions
type Repositorer interface {
	db.Querier

	CheckStock(ctx context.Context, arg CreateStockParams) (bool, error)
	ReserveStock(ctx context.Context, arg CreateStockParams) (bool, error)
	ReleaseStock(ctx context.Context, arg CreateStockParams) (bool, error)
}

// repo provides all functions to execute SQL queries and transactions
type repo struct {
	pool *pgxpool.Pool

	// composition
	*db.Queries
}

// NewRepo creates a new Repo
func NewRepo(pool *pgxpool.Pool) Repositorer {
	return &repo{
		pool:    pool,
		Queries: db.New(pool),
	}
}

// execTx executes a function within a database transaction
func (r *repo) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
