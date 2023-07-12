package health

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ABCD struct {
}

type Server struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (h *Server) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (h *Server) Watch(request *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	return nil
}
