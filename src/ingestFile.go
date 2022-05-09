package src

import (
	"bytes"
	"context"
	"encoding/csv"
	"file_reader/src/instrument"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type avroCodec interface {
	Serialize(io.Writer) error
}

type IngestConfig struct {
	BrokerAddrs []string
	Reader      io.Reader
	Context     context.Context
	Logger      log.Logger
}

type Operation struct {
	Topic         string
	Key           string
	SchemaIDBytes []byte
	RowToSchema   func(row []string) avroCodec
}

func (op Operation) GetNewKafkaWriter(config IngestConfig) *kafka.Writer {

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
		Compression:  compress.Snappy,
		ReadTimeout:  time.Duration(writerReadTimeout * int(time.Second)),
		WriteTimeout: time.Duration(writerWriteTimeout * int(time.Second)),
	}
	return w
}

func (op Operation) IngestFile(config IngestConfig, fileType string) {
	switch fileType {

	case "text/csv":
		csvReader := csv.NewReader(config.Reader)
		w := kafka.Writer{
			Addr:   kafka.TCP(config.BrokerAddrs...),
			Topic:  op.Topic,
			Logger: &config.Logger,
		}
		fmt.Println("INGEST TOPIC", op.Topic, config.BrokerAddrs)
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
			schemaCodec := op.RowToSchema(row)
			schemaCodec.Serialize(&buf)
			valueBytes := buf.Bytes()

			//Combine row bytes with schema id to make a record
			var recordValue []byte
			recordValue = append(recordValue, byte(0))
			recordValue = append(recordValue, op.SchemaIDBytes...)
			recordValue = append(recordValue, valueBytes...)

			// Put the row on the topic
			err = w.WriteMessages(
				config.Context,
				kafka.Message{
					Key:   []byte(op.Key),
					Value: recordValue,
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}
}
