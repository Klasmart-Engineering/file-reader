package healthcheck

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	mu sync.Mutex
	// statusMap stores the serving status of the services this HealthServer monitors.
	statusMap map[string]healthpb.HealthCheckResponse_ServingStatus
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		statusMap: make(map[string]healthpb.HealthCheckResponse_ServingStatus),
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
	if status, ok := s.statusMap[in.Service]; ok {
		return &healthpb.HealthCheckResponse{
			Status: status,
		}, nil
	}
	return nil, grpc.Errorf(codes.NotFound, "Unknown Service")
}

func (s *HealthServer) Watch(in *healthpb.HealthCheckRequest, stream healthgrpc.Health_WatchServer) error {
	return stream.Send(&healthgrpc.HealthCheckResponse{
		Status: healthgrpc.HealthCheckResponse_SERVING,
	})
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *HealthServer) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.mu.Lock()
	s.statusMap[service] = status
	s.mu.Unlock()
}
