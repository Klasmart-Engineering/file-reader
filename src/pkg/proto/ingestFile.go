package proto

import (
	"context"
	"file_reader/src/instrument"
	"file_reader/src/log"
	"fmt"
	"io"
	"strconv"
	"time"

	csvError "file_reader/src/pkg"
	"file_reader/src/pkg/reader"
	"file_reader/src/third_party/protobuf"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type Config struct {
	BrokerAddrs []string
	Reader      reader.Reader
	Context     context.Context
	Logger      *log.ZapLogger
}

type Operation struct {
	Topic       string
	Key         string
	SchemaID    int
	RowToSchema func(row []string, trackingId string) (interface{}, error)
}

var ProtoOperationMap map[string]Operation

// A map that store proto topic for each operation
var topicMap = func() map[string]string {
	return map[string]string{
		"ORGANIZATION": instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC"),
	}
}

func CreateOperationMapProto() map[string]Operation {
	// creates a map of key to Operation struct

	return map[string]Operation{
		"ORGANIZATION": {
			Topic:       topicMap()["ORGANIZATION"],
			Key:         "",
			SchemaID:    schemaRegistryClient.GetSchemaID(topicMap()["ORGANIZATION"]),
			RowToSchema: RowToOrganization,
		},
	}
}
func (op Operation) GetNewKafkaWriter(config Config) *kafka.Writer {

	writerRequiredAcks, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_REQUIRED_ACKS"))
	writerMaxAttempts, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_MAX_ATTEMPTS"))
	writerReadTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_READ_TIMEOUT"))

	writerWriteTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_WRITE_TIMEOUT"))
	w := &kafka.Writer{
		Addr:         kafka.TCP(config.BrokerAddrs...),
		Topic:        op.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(writerRequiredAcks),
		MaxAttempts:  writerMaxAttempts,
		//Logger:       kafka.LoggerFunc(config.Logger.),
		Compression:            compress.Snappy,
		ReadTimeout:            time.Duration(writerReadTimeout * int(time.Second)),
		WriteTimeout:           time.Duration(writerWriteTimeout * int(time.Second)),
		AllowAutoTopicCreation: instrument.IsEnv("TEST"),
	}
	return w
}

func (op Operation) IngestFilePROTO(config Config, trackingId string) (errorStr string) {

	var errors []error

	w := op.GetNewKafkaWriter(config)
	serde := protobuf.NewProtoSerDe()
	fails := 0 // Keep track of how many rows failed to process
	total := 0 // Total of row
	for {
		row, err := config.Reader.Read()
		if err == io.EOF {
			break
		}
		total += 1
		if err != nil {
			config.Logger.Errorf(config.Context, err.Error())
			// If there is an error when reading the current row then skip it
			errors = append(errors, csvError.CSVRowError{RowNum: total, Message: "Can't read the data"})
			fails += 1
			continue
		}
		// Process the row, if it fails then skip that row
		schema, err := op.RowToSchema(row, trackingId)
		if err != nil {
			config.Logger.Errorf(config.Context, err.Error())
			errors = append(errors, csvError.CSVRowError{RowNum: total, Message: "Fail to process the row: Invalid UUID format."})
			fails += 1
			continue
		}

		valueBytes, err := serde.Serialize(op.SchemaID, schema)

		if err != nil {
			config.Logger.Errorf(config.Context, fmt.Sprintf("error serializing message: %s", err))
			errors = append(errors, err)
			return fmt.Sprintf("%s", errors)
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
			config.Logger.Errorf(config.Context, fmt.Sprintf("could not write message: %s", err))
			errors = append(errors, err)
			return fmt.Sprintf("%s", errors)
		}
	}
	// If all the rows are failed to process then
	if total == fails {
		config.Logger.Errorf(config.Context, "All rows are invalid")
		err := csvError.InvalidRowError{Message: "All rows are invalid"}
		return fmt.Sprintf("%s", err)

	}
	return fmt.Sprintf("%s", errors)
}
