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

	csvError "file_reader/src/pkg"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type Config struct {
	Topic       string
	BrokerAddrs []string
	Reader      io.Reader
	Context     context.Context
	Logger      *log.ZapLogger
}

type Operation struct {
	rowToProtoSchema func(row []string, trackingId string) (*orgPb.Organization, error)
}

func (op Operation) GetNewKafkaWriter(config Config) *kafka.Writer {

	writerRequiredAcks, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_REQUIRED_ACKS"))
	writerMaxAttempts, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_MAX_ATTEMPTS"))
	writerReadTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_READ_TIMEOUT"))

	writerWriteTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_WRITE_TIMEOUT"))
	w := &kafka.Writer{
		Addr:         kafka.TCP(config.BrokerAddrs...),
		Topic:        config.Topic,
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

func (op Operation) IngestFilePROTO(config Config, fileTypeName string, trackingId string) (errorStr string) {

	var errors []error
	switch fileTypeName {

	case "CSV":
		csvReader := csv.NewReader(config.Reader)
		w := op.GetNewKafkaWriter(config)

		schemaID, err := schemaRegistryClient.GetProtoSchemaID(organizationSchemaName, config.Topic)
		if err != nil {
			return fmt.Sprintf("%s", err)
		}

		serde := protobuf.NewProtoSerDe()
		fails := 0 // Keep track of how many rows failed to process
		total := 0 // Total of row
		for {
			row, err := csvReader.Read()
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
			orgSchema, err := op.rowToProtoSchema(row, trackingId)
			if err != nil {
				config.Logger.Errorf(config.Context, err.Error())
				errors = append(errors, csvError.CSVRowError{RowNum: total, Message: "Fail to process the row: Invalid UUID format."})
				fails += 1
				continue
			}

			valueBytes, err := serde.Serialize(schemaID, orgSchema)

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
			config.Logger.Errorf(config.Context, err.Error())
			err = csvError.InvalidRowError{Message: "All rows are invalid"}
			return fmt.Sprintf("%s", err)
		}
	}
	return fmt.Sprintf("%s", errors)
}
