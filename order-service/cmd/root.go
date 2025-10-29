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
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/bootstrap"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/messaging"
	"github.com/ymanshur/synasishouse/order/internal/repo"
	"github.com/ymanshur/synasishouse/order/internal/server/api"
	"github.com/ymanshur/synasishouse/order/internal/server/mq"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
	"golang.org/x/sync/errgroup"
)

var interuptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

// Start servers
//
// TODO: Implement timeout propagation
func Start() {
	config := appctx.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), interuptSignals...)
	defer stop()

	dbURL := postgresDSN(config.DB)
	db := bootstrap.RegistryPostgreSQL(ctx, dbURL)

	runDBMigration(config.DBMigrationURL, dbURL)

	repo := repo.NewRepo(db)

	clientConn := bootstrap.RegistryGRPCClient(config.GRPCClient.Inventory)
	defer clientConn.Close()

	inventoryClient := connector.NewInventory(clientConn)

	mqConn := bootstrap.RegistryRabbitMQ(config.RabbitMQ)
	defer mqConn.Close()

	cancelOrderConfig := config.RabbitMQClient.CancelOrder
	autoCancelOrderProducer, err := messaging.NewAutoCancelOrderProducer(mqConn, cancelOrderConfig.Exchange, cancelOrderConfig.Key)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create cancel order producer")
	}

	orderUseCase := usecase.NewOrder(repo, inventoryClient, autoCancelOrderProducer)

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return runMQServer(
			ctx,
			mqConn,
			orderUseCase,
		)
	})

	wg.Go(func() error {
		return runAPIServer(
			ctx,
			orderUseCase,
		)
	})

	err = wg.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start")
	}
}

func runMQServer(
	ctx context.Context,
	conn *amqp091.Connection,
	orderUseCase usecase.Orderer,
) error {
	server := mq.NewServer(conn, orderUseCase)
	err := server.Run(ctx)
	if err != nil {
		return fmt.Errorf("run mq server: %w", err)
	}
	return nil
}

func runAPIServer(
	ctx context.Context,
	orderUseCase usecase.Orderer,
) error {
	server := api.NewServer(orderUseCase)
	err := server.Run(ctx)
	if err != nil {
		return fmt.Errorf("run http server: %w", err)
	}
	return nil
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
