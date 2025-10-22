package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/ymanshur/synasishouse/order/internal/appctx"
	"github.com/ymanshur/synasishouse/order/internal/consts"
	"github.com/ymanshur/synasishouse/order/internal/server"
	"github.com/ymanshur/synasishouse/order/internal/server/api/router"
	"github.com/ymanshur/synasishouse/order/internal/usecase"
)

var _ server.Server = (*Server)(nil) // check in runtime [Server] implement [server.Server]

// Server serves HTTP requests
type Server struct {
	httpServerAddr string
	router         router.Router
}

// NewServer creates a new HTTP server and set up routing
func NewServer(
	config *appctx.Config,
	orderUseCase usecase.Orderer,
) *Server {
	return &Server{
		httpServerAddr: config.HTTPServerAddress,
		router: router.NewRouter(
			config,
			orderUseCase,
		),
	}
}

// Run the HTTP server
func (s *Server) Run(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    s.httpServerAddr,
		Handler: s.Route(),
	}

	log.Info().Msgf("start HTTP server at %s", httpServer.Addr)

	errServe := make(chan error)

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			errServe <- fmt.Errorf("start HTTP server: %w", err)
		}
	}()

	for {
		select {
		case err := <-errServe:
			return err

		// waiting for a interupt signal
		case <-ctx.Done():
			ctxShutdown, cancel := context.WithTimeout(context.Background(), consts.DefaultServerTimeout)
			defer cancel()

			// shutdown stop the HTTP server gratefully
			if err := httpServer.Shutdown(ctxShutdown); err != nil {
				return fmt.Errorf("shutdown HTTP server: %w", err)
			}

			return nil
		}
	}
}

// Route return the HTTP handler
func (s *Server) Route() http.Handler {
	return s.router.Route()
}
