package healthcheck

import (
	"context"
	"sync/atomic"

	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	ops    int64
	status healthpb.HealthCheckResponse_ServingStatus
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		status: healthpb.HealthCheckResponse_NOT_SERVING,
	}
}

func (s *HealthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	atomic.AddInt64(&s.ops, 1)
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

func (s *HealthServer) Watch(in *healthpb.HealthCheckRequest, stream healthgrpc.Health_WatchServer) error {
	atomic.AddInt64(&s.ops, 1)
	return stream.Send(&healthgrpc.HealthCheckResponse{
		Status: healthgrpc.HealthCheckResponse_SERVING,
	})
}

// SetServingStatus is called when need to reset the serving status of a service
// or insert a new service entry into the statusMap.
func (s *HealthServer) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	atomic.AddInt64(&s.ops, 1)
	s.status = status

}
