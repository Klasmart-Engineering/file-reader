package src

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/segmentio/kafka-go"
)

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
	RowToSchema   func(row []string) avroCodec
}

func (op Operation) IngestFile(ctx context.Context, reader Reader, kafkaWriter kafka.Writer) {
	for {
		row, err := reader.Read()
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
		err = kafkaWriter.WriteMessages(
			ctx,
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
