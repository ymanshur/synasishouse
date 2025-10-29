package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/inventory/internal/appctx"
	"github.com/ymanshur/synasishouse/inventory/internal/bootstrap"
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

	dbURL := postgresDSN(config.DB)
	db := bootstrap.RegistryPostgreSQL(ctx, dbURL)

	runDBMigration(config.DBMigrationURL, dbURL)

	repo := repo.NewRepo(db)

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

func runDBMigration(migrationURL string, dbURL string) {
	migration, err := migrate.New(migrationURL, dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run migrate up")
	}
}

// postgresDSN return PostgreSQL Data Source Name
func postgresDSN(config appctx.DBConfig) string {
	param := url.Values{}
	param.Add("user", url.QueryEscape(config.User))
	param.Add("password", url.QueryEscape(config.Password))
	param.Add("port", fmt.Sprint(config.Port))
	param.Add("sslmode", "disable")

	dsn := fmt.Sprintf("postgresql://%s/%s?%s",
		config.Host,
		config.Name,
		param.Encode(),
	)

	return dsn
}
