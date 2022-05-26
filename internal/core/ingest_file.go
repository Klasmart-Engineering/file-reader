package core

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	zaplogger "github.com/KL-Engineering/file-reader/internal/log"

	"github.com/segmentio/kafka-go"
)

type IngestFileConfig struct {
	KafkaWriter kafka.Writer
	TrackingId  string
	Logger      *zaplogger.ZapLogger
}
type Operation struct {
	Topic        string
	Key          string
	SchemaID     int
	SerializeRow func(row []string, trackingId string, schemaId int) ([]byte, error)
}

func (op Operation) IngestFile(ctx context.Context, fileRows chan []string, config IngestFileConfig) {
	logger := config.Logger
	for row := range fileRows {
		// Serialise row using schema
		fmt.Println("row = ", row)
		recordValue, err := op.SerializeRow(row, config.TrackingId, op.SchemaID)
		fmt.Println("recordValue = ", recordValue)
		if err != nil {
			logger.Error(ctx, "Error serialising record to bytes", err)
			continue
		}
		// Put the row on the topic
		err = config.KafkaWriter.WriteMessages(
			ctx,
			kafka.Message{
				Key:   []byte(op.Key),
				Value: recordValue,
			},
		)
		if err != nil {
			logger.Info(ctx, "could not write message "+err.Error())
			continue
		}
	}
}

func ReadRows(ctx context.Context, logger *zaplogger.ZapLogger, f *os.File, fileType string, fileRows chan []string) {
	defer f.Close()
	// Use different reader depending on filetype
	var reader Reader
	switch fileType {
	default:
		reader = csv.NewReader(f)
	}
	// Read rows and pass to fileRows channel
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(ctx, "Error reading from file", err)
			continue
		}
		fileRows <- row
	}
	close(fileRows)
}
