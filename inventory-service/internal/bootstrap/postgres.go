package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// RegistryPostgreSQL establish PostgreSQL connection pool
func RegistryPostgreSQL(ctx context.Context, dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to posgresql")
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		log.Fatal().Err(err).Msg("cannot ping posgresql")
	}

	return pool
}
