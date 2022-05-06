package main

import (
	"context"
	"file_reader/src/config"
	"file_reader/src/instrument"
	"file_reader/src/log"
<<<<<<< HEAD
	"file_reader/src/protos/csvfile"
=======
	csvpb "file_reader/src/protos/csvfile"
>>>>>>> a4a7270 (resolve conflict)

	csvGrpc "file_reader/src/services/organization/delivery/grpc"

	"go.uber.org/zap"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func grpcServerInstrument() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, _ := zap.NewDevelopment()

	logger := log.Wrap(l)

	Logger := config.Logger{
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		Level:             "info",
	}

	addr := instrument.GetAddressForHealthCheck()

	// grpc Server instrument
	lis, grpcServer, err := instrument.GetInstrumentGrpcServer("Csv health check", addr, logger)

	if err != nil {
		logger.Fatalf(ctx, "Failed to start server. Error : %v", err)
	}

	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers: instrument.GetBrokers(),
		},
	}

	csvFileService := csvGrpc.NewCsvFileService(ctx, logger, cfg)
	healthServer := health.NewServer()

	csvfile.RegisterCsvFileServiceServer(grpcServer, csvFileService)

	//healthService := healthcheck.NewHealthChecker()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(csvfile.CsvFileService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	logger.Infof(ctx, "Server starting to listen on %s", addr)

	if err = grpcServer.Serve(lis); err != nil {
		logger.Error(ctx, "Error while starting the gRPC server on the %s listen address %v", lis, err.Error())
	}

}

func main() {
	// grpc Server instrument
	grpcServerInstrument()
}
