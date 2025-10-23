package gapi

import (
	"context"
	"errors"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/inventory/internal/appctx"
	"github.com/ymanshur/synasishouse/inventory/internal/server/gapi/interceptor"
	"github.com/ymanshur/synasishouse/inventory/internal/usecase"
	pb "github.com/ymanshur/synasishouse/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	// UnimplementedInventoryServer enable forward compatibility
	pb.UnimplementedInventoryServer

	rpcServerAddr string
	rpcServer     *grpc.Server

	productUseCase usecase.Producter
	stockUseCase   usecase.Stocker
}

// NewServer creates a new gRPC server.
func NewServer(
	productUseCase usecase.Producter,
	stockUseCase usecase.Stocker,
) (*Server, error) {
	config := appctx.LoadConfig()

	server := &Server{
		rpcServerAddr:  config.GRPCServer.GetAddr(),
		productUseCase: productUseCase,
		stockUseCase:   stockUseCase,
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		log.Error().
			Any("panic", p).
			Bytes("stack", debug.Stack()).
			Msg("recovered from panic")

		return status.Errorf(codes.Internal, "%s", p)
	}
	grpcRecovery := recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler))

	server.rpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.Logger,
		grpcRecovery,
	))
	pb.RegisterInventoryServer(server.rpcServer, server)

	// Allows the gRPC client to explore available RPCs on the server
	// as some kind of self server documentation.
	reflection.Register(server.rpcServer)

	return server, nil
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.rpcServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	errServe := make(chan error)

	// Start the service listening for requests.
	go func() {
		errServe <- s.rpcServer.Serve(listener)
	}()

	// Blocking and waiting for shutdown or error from server.
	select {
	case err := <-errServe:
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			return err
		}
	case <-ctx.Done():
		s.rpcServer.GracefulStop()
	}

	return nil
}
