package health

import (
	"context"
	"net/http"

	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const DefaultRoute = "/healthcheck"

// New returns new grpchealth.Service
func New() *Service {
	return &Service{}
}

// Service implements grpc_health_v1.HealthServer
type Service struct {
	ready atomic.Bool
}

// Register registers health service.
func (s *Service) Register(grpcServer *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(grpcServer, s)
}

// RegisterHTTP registers health service.
func (s *Service) RegisterHTTP(mux *http.ServeMux) {
	mux.Handle(DefaultRoute, s)
}

// MarkReady marks the service ready
func (s *Service) MarkReady() {
	s.ready.Store(true)
}

// Shutdown marks the service not ready
func (s *Service) Shutdown() {
	s.ready.Store(false)
}

// Check reports whether the service is ready
func (s *Service) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	resp := &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	}
	if s.ready.Load() {
		resp.Status = grpc_health_v1.HealthCheckResponse_SERVING
	}
	return resp, nil
}

// Watch is unimplemented
func (s *Service) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

// ServeHTTP implements http healthcheck
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 400
	if s.ready.Load() {
		status = 200
	}
	w.WriteHeader(status)
}
