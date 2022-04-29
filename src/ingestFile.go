package src

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/csv"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

// interface fulfilled by the generated structs
type AvroRecord interface {
	Serialize(io.Writer) error
	AvroCRC64Fingerprint() []byte
}

func ComposeCsvIngester[operation AvroRecord](
	topic string,
	w *kafka.Writer,
	rowToOpConverter func([]string) operation,
	schema *srclient.Schema,
	brokerAddrs []string,
	logger *log.Logger,
) func(r io.Reader, ctx context.Context) {
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	return func(r io.Reader, ctx context.Context) {
		csvReader := csv.NewReader(r)
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
			op := rowToOpConverter(row)
			op.Serialize(&buf)
			valueBytes := buf.Bytes()
			if err != nil {
				panic("error converting to avro bytes " + err.Error())
			}

			// Combine row bytes with schema id to make a record
			var recordValue []byte
			recordValue = append(recordValue, byte(0))
			recordValue = append(recordValue, schemaIDBytes...)
			recordValue = append(recordValue, valueBytes...)

			// Put the row on the topic
			key, _ := uuid.NewUUID()
			err = w.WriteMessages(
				ctx,
				kafka.Message{
					Key:   []byte(key.String()),
					Value: recordValue,
				},
			)
			if err != nil {
				panic("could not write message " + err.Error())
			}
		}
	}
}
