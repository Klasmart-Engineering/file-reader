package test

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	avro "file_reader/avro_gencode"
	src "file_reader/src"

	"github.com/segmentio/kafka-go"
)

func TestReadOrgCsv(t *testing.T) {
	const csvPath = "../data.csv"
	brokerAddrs := []string{"localhost:9092"}

	//Setup
	f, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	ctx := context.Background()

	// Compose File reader for organization
	csvIngester := src.CreateCsvIngester(brokerAddrs)

	// Actually put the csv rows onto the topic
	csvIngester(f, ctx)

	// Consume from the topic to prove it worked
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddrs,
		Topic:   "organisation",
	})
	for i := 1; i < 5; i++ {
		msg, err := r.ReadMessage(ctx)
		if err == nil {
			val, e := avro.DeserializeOrganization(bytes.NewReader(msg.Value[5:]))
			if e == nil {
				t.Logf("Here is the message %s\n", val)
			} else {
				t.Logf("Error deserializing: %e", e)
			}
		} else {
			t.Logf("Error consuming the message: %v (%v)\n", err, msg)
			break
		}
	}
}
