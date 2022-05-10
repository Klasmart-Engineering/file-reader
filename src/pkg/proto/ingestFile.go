package proto

import (
	"context"
	"encoding/csv"
	"file_reader/src/instrument"
<<<<<<< HEAD
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
=======
	orgPb "file_reader/src/protos/onboarding"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

	"google.golang.org/protobuf/proto"
>>>>>>> 989d433 (Create a separate package for proto logic)
)

type Config struct {
	BrokerAddrs []string
	Reader      io.Reader
	Context     context.Context
<<<<<<< HEAD
	Logger      *log.ZapLogger
=======
	Logger      log.Logger
>>>>>>> 989d433 (Create a separate package for proto logic)
}

type Operation struct {
	topic            string
<<<<<<< HEAD
	schema           *srclient.Schema
	rowToProtoSchema func(row []string) (*orgPb.Organization, error)
=======
	key              string
	schema           *srclient.Schema
	rowToProtoSchema func(row []string) *orgPb.Organization
>>>>>>> 989d433 (Create a separate package for proto logic)
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

<<<<<<< HEAD
func (op Operation) IngestFilePROTO(config Config, fileTypeName string) error {
=======
func (op Operation) IngestFilePROTO(config Config, fileTypeName string) {
>>>>>>> 989d433 (Create a separate package for proto logic)

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
<<<<<<< HEAD
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
=======
				log.Fatal(err)
			}
			// Serialise row using schema
			orgSchema := op.rowToProtoSchema(row)
			value, err := proto.Marshal(orgSchema)
			valueBytes, _ := op.schema.Codec().BinaryFromNative(nil, value)

			//Combine row bytes with schema id to make a record
			var recordValue []byte
			schemaIDBytes := GetSchemaIdBytes(op.schema)
			recordValue = append(recordValue, byte(0))
			recordValue = append(recordValue, schemaIDBytes...)
			recordValue = append(recordValue, valueBytes...)
>>>>>>> 989d433 (Create a separate package for proto logic)

			// Put the row on the topic
			err = w.WriteMessages(
				config.Context,
				kafka.Message{
<<<<<<< HEAD
					Key:   []byte(""),
					Value: valueBytes,
=======
					Key:   []byte(op.key),
					Value: recordValue,
>>>>>>> 989d433 (Create a separate package for proto logic)
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}
<<<<<<< HEAD
	return nil
=======
>>>>>>> 989d433 (Create a separate package for proto logic)
}
