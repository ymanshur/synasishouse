package gapi

import (
	"context"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

// Check is a method returns health response.
func (s *Server) Check(context.Context, *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is a method returns health response.
func (s *Server) Watch(r *health.HealthCheckRequest, server health.Health_WatchServer) error {
	return server.Send(&health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	})
}
