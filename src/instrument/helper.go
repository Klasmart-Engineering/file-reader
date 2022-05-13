package instrument

import (
	"file_reader/src/log"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/newrelic/go-agent/_integrations/nrgrpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func MustGetEnv(key string) string {
	if val := os.Getenv(key); "" != val {
		return val
	}
	panic(fmt.Sprintf("environment variable %s unset", key))
}

func GetAddressForHealthCheck() string {
	host := MustGetEnv("GRPC_SERVER")
	port := MustGetEnv("GRPC_HEALTH_PORT")
	addr := host + ":" + port
	return addr

}

func GetAddressForGrpc() string {
	host := MustGetEnv("GRPC_SERVER")
	port := MustGetEnv("GRPC_SERVER_PORT")
	addr := host + ":" + port
	return addr

}

func GetBrokers() []string {
	return strings.Split(MustGetEnv("BROKERS"), ",")
}

func GetGrpcServer(serviceName string, address string, logger *log.ZapLogger) (net.Listener, *grpc.Server, error) {
	l, err := net.Listen("tcp", address)

	if err != nil {
		return nil, nil, errors.Wrap(err, "net.Listen")
	}
	newRelicEnabled := MustGetEnv("NEW_RELIC_ENABLED")
	var grpcServer *grpc.Server
	if newRelicEnabled == "true" {
		// grpc Server instrument
		nr, _ := GetNewRelic(serviceName, logger)

		grpcServer = grpc.NewServer(
			// Add the New Relic gRPC server instrumentation
			grpc.UnaryInterceptor(nrgrpc.UnaryServerInterceptor(nr.app)),
			grpc.StreamInterceptor(nrgrpc.StreamServerInterceptor(nr.app)),
		)
	} else {
		grpcServer = grpc.NewServer()
	}
	return l, grpcServer, nil

}
