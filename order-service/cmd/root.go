package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/connector"
	"github.com/ymanshur/synasishouse/order/internal/server/api"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	conn := bootstrapGRPCConnection(target)
	defer conn.Close()

	inventoryClient := connector.NewInventory(conn)

	orderUseCase := usecase.NewOrder(inventoryClient)

	runAPIServer(
		ctx,
		&config,
		orderUseCase,
	)
}

func bootstrapGRPCConnection(target string) *grpc.ClientConn {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(3),                                            // Max 3 retry attempts
		grpc_retry.WithPerRetryTimeout(1 * time.Second),                  // Timeout for each individual retry
		grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted), // Retry on these gRPC codes
		grpc_retry.WithBackoff(func(_ uint) time.Duration {
			// Custom backoff function with jitter
			baseDelay := 100 * time.Millisecond
			jitter := time.Duration(rand.Intn(int(baseDelay.Milliseconds()/5))) * time.Millisecond // 20% jitter
			return baseDelay + jitter
		}),
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	}
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create grpc channel")
	}
	return conn
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
