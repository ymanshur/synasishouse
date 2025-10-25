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
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/bootstrap"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/repo"
	"github.com/ymanshur/synasishouse/order/internal/server/api"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
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

	db, err := pgxpool.New(ctx, config.DB.GetURL())
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	runDBMigration(config.DBMigrationURL, config.DB.GetURL())

	repo := repo.NewRepo(db)

	conn := bootstrap.RegistryGRPCClient(config.GRPCClient.Inventory)
	defer conn.Close()

	inventoryClient := connector.NewInventory(conn)

	orderUseCase := usecase.NewOrder(repo, inventoryClient)

	runAPIServer(
		ctx,
		orderUseCase,
	)
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

func runAPIServer(
	ctx context.Context,
	orderUseCase usecase.Orderer,
) {
	server := api.NewServer(orderUseCase)

	err := server.Run(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run HTTP server")
	}
}
