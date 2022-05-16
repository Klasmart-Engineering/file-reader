package src

import (
	"bytes"
	"context"
	zaplogger "file_reader/src/log"
	"io"

	"github.com/segmentio/kafka-go"
)

type IngestFileConfig struct {
	Reader      Reader
	KafkaWriter kafka.Writer
	Tracking_id string
	Logger      *zaplogger.ZapLogger
}

type avroCodec interface {
	Serialize(io.Writer) error
}

type Reader interface {
	Read() ([]string, error)
}

type Operation struct {
	Topic         string
	Key           string
	SchemaIDBytes []byte
	RowToSchema   func(row []string, tracking_id string) avroCodec
}

func (op Operation) IngestFile(ctx context.Context, config IngestFileConfig) {
	logger := config.Logger
	for {
		row, err := config.Reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatal(ctx, err)
		}

		// Serialise row using schema
		var buf bytes.Buffer
		schemaCodec := op.RowToSchema(row, config.Tracking_id)
		schemaCodec.Serialize(&buf)
		valueBytes := buf.Bytes()

		//Combine row bytes with schema id to make a record
		var recordValue []byte
		recordValue = append(recordValue, byte(0))
		recordValue = append(recordValue, op.SchemaIDBytes...)
		recordValue = append(recordValue, valueBytes...)

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
