package src

import (
	"bytes"
	"context"
	"encoding/csv"
	"file_reader/src/instrument"
	orgPb "file_reader/src/protos/onboarding"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"

	"google.golang.org/protobuf/proto"
)

type avroCodec interface {
	Serialize(io.Writer) error
}

type Config struct {
	BrokerAddrs []string
	Reader      io.Reader
	Context     context.Context
	Logger      log.Logger
}

type Operation struct {
	topic            string
	key              string
	schema           *srclient.Schema
	rowToSchema      func(row []string) avroCodec
	rowToProtoSchema func(row []string) *orgPb.Organization
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

func (op Operation) IngestFilePROTO(config Config, fileTypeName string) {

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

			// Put the row on the topic
			err = w.WriteMessages(
				config.Context,
				kafka.Message{
					Key:   []byte(op.key),
					Value: recordValue,
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}
}

func (op Operation) IngestFileAVROS(config Config, fileTypeName string) {
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
				log.Fatal(err)
			}

			// Serialise row using schema
			var buf bytes.Buffer
			schemaCodec := op.rowToSchema(row)
			schemaCodec.Serialize(&buf)
			valueBytes := buf.Bytes()

			schemaIDBytes := GetSchemaIdBytes(op.schema)
			//Combine row bytes with schema id to make a record
			var recordValue []byte
			recordValue = append(recordValue, byte(0))
			recordValue = append(recordValue, schemaIDBytes...)
			recordValue = append(recordValue, valueBytes...)

			// Put the row on the topic
			err = w.WriteMessages(
				config.Context,
				kafka.Message{
					Key:   []byte(op.key),
					Value: recordValue,
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}

}
