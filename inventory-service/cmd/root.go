package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/inventory/internal/appctx"
	"github.com/ymanshur/synasishouse/inventory/internal/repo"
	"github.com/ymanshur/synasishouse/inventory/internal/server/gapi"
	"github.com/ymanshur/synasishouse/inventory/internal/usecase"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

// Start servers
func Start() {
	config := appctx.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	runDBMigration(config.DBMigrationURL, config.DBSource)

	repo := repo.NewRepo(conn)

	productUseCase := usecase.NewProduct(repo)
	stockUseCase := usecase.NewStock(repo)

	runGRPCServer(
		ctx,
		productUseCase,
		stockUseCase,
	)
}

func runGRPCServer(
	ctx context.Context,
	productUseCase usecase.Producter,
	stockUseCase usecase.Stocker,
) {
	server, err := gapi.NewServer(productUseCase, stockUseCase)

	err = server.Run(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run gRPC server")
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migrate up")
	}
}
