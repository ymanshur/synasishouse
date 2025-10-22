package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/server/api"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	target := fmt.Sprintf("%s:%d", config.GRPCClientHostInventory, config.GRPCClientPortInventory)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc channel")
	}
	defer conn.Close()

	inventoryClient := connector.NewInventory(conn)

	orderUseCase := usecase.NewOrder(inventoryClient)

	runAPIServer(
		ctx,
		&config,
		orderUseCase,
	)
}

func runAPIServer(
	ctx context.Context,
	config *appctx.Config,
	orderUseCase usecase.Orderer,
) {
	server := api.NewServer(config, orderUseCase)

	err := server.Run(ctx)
	if err != nil {
		log.Fatal().Err(err)
	}
}
