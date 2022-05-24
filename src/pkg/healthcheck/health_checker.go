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

func (s *HealthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if in.Service == "" {
		// check the server overall health status.
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_SERVING,
		}, nil
	}
	return &healthpb.HealthCheckResponse{
		Status: s.status,
	}, nil
}

func (s *HealthServer) Watch(in *healthpb.HealthCheckRequest, stream healthpb.Health_WatchServer) error {
	return stream.Send(&healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	})
}

// SetServingStatus is called when need to reset the serving status of a service
func (s *HealthServer) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()
}
