package server

import (
	"context"
	"file_reader/src/config"
	"file_reader/src/instrument"
	zaplogger "file_reader/src/log"

	filepb "file_reader/src/protos/inputfile"

	fileGrpc "file_reader/src/services/organization/delivery/grpc"

	"go.uber.org/zap"
)

type ingestFileServer struct {
	logger *zaplogger.ZapLogger
	cfg    *config.Config
}

func NewServer(logger *zaplogger.ZapLogger, cfg *config.Config) *ingestFileServer {
	return &ingestFileServer{
		logger: logger,
		cfg:    cfg,
	}
}
func (s *ingestFileServer) Run(ctx context.Context) error {
	s.logger.Infof(ctx, false, "GRPC Server is listening... at port %v\n", s.cfg.Server.Port)
	addr := instrument.GetAddressForGrpc()

	lis, grpcServer, err := instrument.GetGrpcServer("File Processing Server", addr, s.logger)

	if err != nil {

		panic(err)
	}

	defer lis.Close()

	ingestFileService := fileGrpc.NewIngestFileService(ctx, s.logger, s.cfg)

	filepb.RegisterIngestFileServiceServer(grpcServer, ingestFileService)

	s.logger.Infof(ctx, false, "GRPC Server is listening...", zap.String("port", s.cfg.Server.Port))

	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Errorf(ctx, false, "Server issue.", zap.String("error", err.Error()))
	}

	return nil

}
