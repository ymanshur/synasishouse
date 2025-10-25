package common

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// PGErrorCode return PostgreSQL error code
func PGErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
