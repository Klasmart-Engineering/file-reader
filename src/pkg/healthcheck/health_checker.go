package healthcheck

import (
	"context"
	"sync"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	mu     sync.Mutex
	status healthpb.HealthCheckResponse_ServingStatus
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		status: healthpb.HealthCheckResponse_NOT_SERVING,
	}
}

func (s *HealthServer) Check(ctx context.Context) (*healthpb.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return &healthpb.HealthCheckResponse{
		Status: s.status,
	}, nil
}

// SetServingStatus is called when need to reset the serving status of a service
func (s *HealthServer) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()
}
