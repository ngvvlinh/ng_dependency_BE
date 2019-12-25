package health

import (
	"net/http"

	"go.uber.org/atomic"
)

const DefaultRoute = "/healthcheck"

// New returns new health.Service
func New() *Service {
	return &Service{}
}

type Service struct {
	ready atomic.Bool
}

// RegisterHTTPHandler registers health service.
func (s *Service) RegisterHTTPHandler(mux *http.ServeMux) {
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

// ServeHTTP implements http healthcheck
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 400
	if s.ready.Load() {
		status = 200
	}
	w.WriteHeader(status)
}
