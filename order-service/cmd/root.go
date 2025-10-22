package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/bootstrap"
	"github.com/ymanshur/synasishouse/order/internal/connector"
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

	conn := bootstrap.RegistryGRPCClient(config.GRPCClient.Inventory)
	defer conn.Close()

	inventoryClient := connector.NewInventory(conn)

	orderUseCase := usecase.NewOrder(inventoryClient)

	runAPIServer(
		ctx,
		orderUseCase,
	)
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
