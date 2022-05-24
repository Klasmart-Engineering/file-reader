package util

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/KL-Engineering/file-reader/src/config"
	"github.com/KL-Engineering/file-reader/src/instrument"
	"github.com/KL-Engineering/file-reader/src/log"

	filepb "github.com/KL-Engineering/file-reader/src/protos/inputfile"
	fileGrpc "github.com/KL-Engineering/file-reader/src/services/organization/delivery/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func dialer(server *grpc.Server, service *fileGrpc.IngestFileService) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	filepb.RegisterIngestFileServiceServer(server, service)

	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func StartGrpc(logger *log.ZapLogger, cfg *config.Config, addr string) (context.Context, filepb.IngestFileServiceClient) {
	timeout, _ := strconv.Atoi(os.Getenv("DEFAULT_SERVER_TIMEOUT_MS"))
	defaultTimeOut := time.Duration(timeout * int(time.Millisecond))
	ctx, _ := context.WithTimeout(context.Background(), defaultTimeOut)

	ingestFileService := fileGrpc.NewIngestFileService(ctx, logger, cfg)

	_, grpcServer, _ := instrument.GetGrpcServer("Mock service", addr, logger)

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(grpcServer, ingestFileService)))

	if err != nil {
		logger.Errorf(ctx, err.Error())
	}

	client := filepb.NewIngestFileServiceClient(conn)

	return ctx, client
}
