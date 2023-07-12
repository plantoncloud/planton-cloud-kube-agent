package main

import (
	"context"
	"fmt"
	healthcontroller "github.com/plantoncloud-inc/planton-cloud-kube-agent/controller/health"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/auth/token"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/config"
	"github.com/plantoncloud-inc/planton-cloud-kube-agent/internal/opencost/scheduler"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const (
	Port = 8080
	//EnvVarLogLevel can be set to info, debug or error
	EnvVarLogLevel = "LOG_LEVEL"
)

func main() {
	setupLogLevel()
	ctx := context.Background()
	c, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config from environment with error %v", err)
	}

	if err := token.Initialize(ctx, c); err != nil {
		log.Fatalf("failed to intialize token for planton-cloud-service api client conn with error %v", err)
	}

	log.Infof("access token has been successfully initialized")

	go token.StartRotator(ctx, c)

	log.Infof("access token rotator has been successfully started")

	log.Infof("starting planton-cloud-kube-agent on %s hosting-env",
		c.PlantonCloudKubeAgentHostingClusterId)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go scheduler.Start(ctx, c)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	registerServices(grpcServer)

	log.Printf("grpc server listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupLogLevel() {
	logLevel, ok := os.LookupEnv(EnvVarLogLevel)
	if !ok {
		log.SetLevel(log.InfoLevel)
	}
	switch logLevel {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func registerServices(grpcServer *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(grpcServer, &healthcontroller.Server{})
}
