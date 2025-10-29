package bootstrap

import (
	"math/rand"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

// RegistryGRPCClient establish gRPC client channel
func RegistryGRPCClient(config config.GRPCClientConfig) *grpc.ClientConn {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithMax(config.MaxRetry),                              // Max 3 retry attempts
		grpc_retry.WithPerRetryTimeout(config.PerRetryTimeout),           // Timeout for each individual retry
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
	conn, err := grpc.NewClient(config.GetAddr(), opts...)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC channel")
	}

	return conn
}
