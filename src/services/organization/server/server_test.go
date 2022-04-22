package server_test

import (
	"file_reader/src/config"
	"file_reader/src/instrument"
	"file_reader/src/services/organization/server"
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestCsvProcessingServer(t *testing.T) {
	log, _ := zap.NewDevelopment()
	Logger := config.Logger{
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		Level:             "info",
	}
	host := instrument.MustGetEnv("GRPC_SERVER")
	port := instrument.MustGetEnv("GRPC_SERVER_PORT")
	addr := host + ":" + port
	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers: []string{"localhost:9092"},
		},
	}
	orgServer := server.NewServer(*log, cfg)

	err := orgServer.Run()

	if err != nil {
		fmt.Printf("Failed to run server %v", err)
		return
	}
	t.Logf("Starting server on port: %v", cfg.Server.Port)

}
