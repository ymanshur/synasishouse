package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/ymanshur/synasishouse/order/db/sqlc"
)

// Repo defines all functions to execute db queries and transactions
type Repo interface {
	db.Querier

	Transact(ctx context.Context, iso pgx.TxIsoLevel, txFunc func(Repo) error) (err error)
}

// repo provides all functions to execute SQL queries and transactions
type repo struct {
	pool *pgxpool.Pool
	tx   *pgx.Tx

	// composition
	*db.Queries
}

// NewRepo creates a new Repo
func NewRepo(pool *pgxpool.Pool) Repo {
	return &repo{
		pool:    pool,
		Queries: db.New(pool),
	}
}

// newRepoWithTx creates a new Repo
func newRepoWithTx(pool *pgxpool.Pool, tx pgx.Tx) Repo {
	return &repo{
		pool:    pool,
		Queries: db.New(tx),
	}
}

// InTransaction check if the connection already in transaction
func (r *repo) InTransaction() bool {
	return r.tx != nil
}

// Transact executes the given function in the context of a SQL transaction at
// the given isolation level.
// For the levels which require retry, see: https://www.postgresql.org/docs/11/transaction-iso.html.
func (r *repo) Transact(ctx context.Context, iso pgx.TxIsoLevel, fn func(Repo) error) (err error) {
	opts := pgx.TxOptions{IsoLevel: iso}
	return r.execTx(ctx, opts, fn)
}

// execTx executes a function within a database transaction
func (r *repo) execTx(ctx context.Context, opts pgx.TxOptions, fn func(Repo) error) error {
	if r.InTransaction() {
		return fmt.Errorf("repository already in a transaction")
	}
	defer func() {
		r.tx = nil
	}()

	tx, err := r.pool.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	repo := newRepoWithTx(r.pool, tx)
	err = fn(repo)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
