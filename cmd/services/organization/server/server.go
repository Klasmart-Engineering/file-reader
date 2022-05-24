package server

import (
	"context"

	"github.com/KL-Engineering/file-reader/cmd/config"
	"github.com/KL-Engineering/file-reader/cmd/instrument"
	zaplogger "github.com/KL-Engineering/file-reader/cmd/log"

	filepb "github.com/KL-Engineering/file-reader/cmd/protos/inputfile"

	fileGrpc "github.com/KL-Engineering/file-reader/cmd/services/organization/delivery/grpc"

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
	s.logger.Infof(ctx, "GRPC Server is listening... at port %v\n", s.cfg.Server.Port)
	addr := instrument.GetAddressForGrpc()

	lis, grpcServer, err := instrument.GetGrpcServer("File Processing Server", addr, s.logger)

	if err != nil {

		panic(err)
	}

	defer lis.Close()

	ingestFileService := fileGrpc.NewIngestFileService(ctx, s.logger, s.cfg)

	filepb.RegisterIngestFileServiceServer(grpcServer, ingestFileService)

	s.logger.Infof(ctx, "GRPC Server is listening...", zap.String("port", s.cfg.Server.Port))

	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Errorf(ctx, "Server issue.", zap.String("error", err.Error()))
	}

	return nil

}
