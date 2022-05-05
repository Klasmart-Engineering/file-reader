package server

import (
	"context"
	"file_reader/src/config"
	"file_reader/src/instrument"
	"file_reader/src/log"

	"file_reader/src/protos"
	"file_reader/src/services/organization/delivery/kafka"

	orgGrpc "file_reader/src/services/organization/delivery/grpc"

	"go.uber.org/zap"
)

type csvFileServer struct {
	logger *log.ZapLogger
	cfg    *config.Config
}

func NewServer(logger *log.ZapLogger, cfg *config.Config) *csvFileServer {
	return &csvFileServer{
		logger: logger,
		cfg:    cfg,
	}
}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	validate := validator.New()
	organizationProducer := kafka.NewOrganizationProducer(s.log, s.cfg)
	organizationProducer.Run()

	defer organizationProducer.Close()

	organizationUC := usecase.NewOrganizationUC(s.log, organizationProducer)

	s.logger.Infof(ctx, "GRPC Server is listening... at port %v\n", s.cfg.Server.Port)
	addr := instrument.GetAddressForGrpc()

	lis, grpcServer, err := instrument.GetInstrumentGrpcServer("Csv Processing Server", addr, s.logger)

	if err != nil {

		panic(err)
	}

	defer lis.Close()

	csvFileService := csvGrpc.NewCsvFileService(ctx, s.logger, s.cfg)

	orgConsumeGroup := kafka.NewOrganizationsConsumerGroup(s.cfg.Kafka.Brokers, kafkaGroupID, s.log, s.cfg, organizationUC, validate)

	s.logger.Infof(ctx, "GRPC Server is listening...", zap.String("port", s.cfg.Server.Port))

	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Errorf(ctx, "Server issue.", zap.String("error", err.Error()))
	}

	return nil

}
