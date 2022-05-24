package main

import (
	"context"

	"github.com/KL-Engineering/file-reader/internal/filereader"
	"github.com/KL-Engineering/file-reader/src/config"
	"github.com/KL-Engineering/file-reader/src/instrument"
	zaplogger "github.com/KL-Engineering/file-reader/src/log"
	filepb "github.com/KL-Engineering/file-reader/src/protos/inputfile"

	fileGrpc "github.com/KL-Engineering/file-reader/src/services/organization/delivery/grpc"

	"go.uber.org/zap"

	health "github.com/KL-Engineering/file-reader/src/pkg/healthcheck"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var IngestFileService *fileGrpc.IngestFileService

func grpcServerInstrument(ctx context.Context, logger *zaplogger.ZapLogger) {
	Logger := config.Logger{
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		Level:             "info",
	}

	addr := instrument.GetAddressForHealthCheck()

	// grpc Server instrument
	lis, grpcServer, err := instrument.GetGrpcServer("File service health check", addr, logger)

	if err != nil {
		logger.Fatalf(ctx, "Failed to start server. Error : %v", err)
	}

	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers:                instrument.GetBrokers(),
			AllowAutoTopicCreation: true,
		},
	}

	IngestFileService = fileGrpc.NewIngestFileService(ctx, logger, cfg)

	filepb.RegisterIngestFileServiceServer(grpcServer, IngestFileService)

	healthServer := health.NewHealthServer()
	healthServer.SetServingStatus("File reader service", 1)

	healthpb.RegisterHealthServer(grpcServer, healthServer)

	logger.Infof(ctx, "Server starting to listen on %s", addr)

	if err = grpcServer.Serve(lis); err != nil {
		logger.Error(ctx, "Error while starting the gRPC server on the %s listen address %v", lis, err.Error())
	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var l *zap.Logger
	mode := instrument.MustGetEnv("MODE")
	switch mode {
	case "debug":
		l = zap.NewNop()
	default:
		l, _ = zap.NewDevelopment()
	}

	logger := zaplogger.Wrap(l)

	// Initialise Consumer for file create event
	filereader.StartFileCreateConsumer(ctx, logger)

	// grpc Server instrument
	grpcServerInstrument(ctx, logger)
}
