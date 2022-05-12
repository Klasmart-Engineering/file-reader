package proto

import (
	"context"
	"encoding/csv"
	"file_reader/src/instrument"
	"file_reader/src/log"
	orgPb "file_reader/src/protos/onboarding"
	"fmt"
	"io"
	"strconv"
	"time"

	"file_reader/src/third_party/protobuf"

	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type Config struct {
	BrokerAddrs []string
	Reader      io.Reader
	Context     context.Context
	Logger      *log.ZapLogger
}

type Operation struct {
	topic            string
	schema           *srclient.Schema
	rowToProtoSchema func(row []string) (*orgPb.Organization, error)
}

func (op Operation) GetNewKafkaWriter(config Config) *kafka.Writer {

	writerRequiredAcks, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_REQUIRED_ACKS"))
	writerMaxAttempts, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_MAX_ATTEMPTS"))
	writerReadTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_READ_TIMEOUT"))

	writerWriteTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_WRITE_TIMEOUT"))
	w := &kafka.Writer{
		Addr:         kafka.TCP(config.BrokerAddrs...),
		Topic:        op.topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(writerRequiredAcks),
		MaxAttempts:  writerMaxAttempts,
		//Logger:       kafka.LoggerFunc(config.Logger.),
		Compression:  compress.Snappy,
		ReadTimeout:  time.Duration(writerReadTimeout * int(time.Second)),
		WriteTimeout: time.Duration(writerWriteTimeout * int(time.Second)),
	}
	return w
}

func (op Operation) IngestFilePROTO(config Config, fileTypeName string) error {

	switch fileTypeName {

	case "CSV":
		csvReader := csv.NewReader(config.Reader)
		w := op.GetNewKafkaWriter(config)
		for {
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				config.Logger.Fatalf(config.Context, err.Error())
				return err
			}
			// Serialise row using schema
			orgSchema, err := op.rowToProtoSchema(row)
			if err != nil {
				config.Logger.Fatalf(config.Context, err.Error())
			}
			schemaID := schemaRegistryClient.GetProtoSchemaID(organizationSchemaName, organizationProtoTopic)
			serde := protobuf.NewProtoSerDe()

			valueBytes, err := serde.Serialize(schemaID, orgSchema)

			if err != nil {
				config.Logger.Fatalf(config.Context, fmt.Sprintf("error serializing message: %w", err))
			}

			// Put the row on the topic
			err = w.WriteMessages(
				config.Context,
				kafka.Message{
					Key:   []byte(""),
					Value: valueBytes,
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}
	return nil
}
