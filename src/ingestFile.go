package src

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"log"

	"github.com/segmentio/kafka-go"
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
	topic         string
	key           string
	schemaIDBytes []byte
	rowToSchema   func(row []string) avroCodec
}

func (op Operation) IngestFile(config Config) {
	csvReader := csv.NewReader(config.Reader)
	w := kafka.Writer{
		Addr:   kafka.TCP(config.BrokerAddrs...),
		Topic:  op.topic,
		Logger: &config.Logger,
	}
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

		//Combine row bytes with schema id to make a record
		var recordValue []byte
		recordValue = append(recordValue, byte(0))
		recordValue = append(recordValue, op.schemaIDBytes...)
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
