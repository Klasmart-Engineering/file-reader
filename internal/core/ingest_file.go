package core

import (
	"context"
	"fmt"
	"io"

	zaplogger "github.com/KL-Engineering/file-reader/internal/log"

	"github.com/segmentio/kafka-go"
)

type IngestFileConfig struct {
	Reader      Reader
	KafkaWriter kafka.Writer
	TrackingId  string
	Logger      *zaplogger.ZapLogger
}

type Reader interface {
	Read() ([]string, error)
}

type Operation struct {
	Topic        string
	Key          string
	SchemaID     int
	SerializeRow func(row []string, trackingId string, schemaId int) ([]byte, error)
}

func (op Operation) IngestFile(ctx context.Context, config IngestFileConfig) {
	logger := config.Logger
	for {
		row, err := config.Reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(ctx, "Error reading from file", err)
			continue
		}

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
