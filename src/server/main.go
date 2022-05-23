package main

import (
	"context"
	"file_reader/internal/filereader"
	"file_reader/src/config"
	"file_reader/src/instrument"
	zaplogger "file_reader/src/log"
	filepb "file_reader/src/protos/inputfile"

	fileGrpc "file_reader/src/services/organization/delivery/grpc"

	"go.uber.org/zap"

	"google.golang.org/grpc/health"
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
	healthServer := health.NewServer()

	filepb.RegisterIngestFileServiceServer(grpcServer, IngestFileService)

	//healthService := healthcheck.NewHealthChecker()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(filepb.IngestFileService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

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
