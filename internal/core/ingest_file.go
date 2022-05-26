package core

import (
	"context"
	"encoding/csv"
	"errors"
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
	Topic         string
	Key           string
	SchemaID      int
	SerializeRow  func(row []string, trackingId string, schemaId int, headerIndexes map[string]int) ([]byte, error)
	HeaderIndexes map[string]int
}

func UpdateHeaderIndexes(headerIndexes map[string]int, row []string) (map[string]int, error) {
	for i, col := range row {
		_, exists := headerIndexes[col]
		if !exists {
			return nil, errors.New(fmt.Sprint("unrecognised header in file", row))
		}
		headerIndexes[col] = i
	}
	return headerIndexes, nil
}

func (op Operation) IngestFile(ctx context.Context, fileRows chan []string, config IngestFileConfig) {
	logger := config.Logger
	for row := range fileRows {
		// Serialise row using schema
		recordValue, err := op.SerializeRow(row, config.TrackingId, op.SchemaID, op.HeaderIndexes)
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
