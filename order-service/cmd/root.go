package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/server/api"
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

	runAPIServer(
		ctx,
		&config,
	)
}

func runAPIServer(
	ctx context.Context,
	config *appctx.Config,
) {
	server := api.NewServer(config)

	err := server.Run(ctx)
	if err != nil {
		log.Fatal().Err(err)
	}
}
