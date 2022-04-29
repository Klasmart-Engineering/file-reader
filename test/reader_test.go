package test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	avro "file_reader/avro_gencode"
	src "file_reader/src"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Skip for now
func TestReadOrgCsv(t *testing.T) {

	t.Skip("Skip for now")
	brokerAddrs := []string{"localhost:9092"}

	// Fake CSV data
	orgId1 := uuid.NewString()
	orgId2 := uuid.NewString()
	orgName1 := "org_name1"
	orgName2 := "org_name2"
	csv := fmt.Sprintf("%s,%s\n%s,%s", orgId1, orgName1, orgId2, orgName2)
	reader := bytes.NewReader([]byte(csv))

	brokerAddrs := []string{"localhost:9092"}

	var config = src.Config{
		BrokerAddrs: brokerAddrs,
		Reader:      reader,
		Context:     context.Background(),
		Logger:      *log.New(os.Stdout, "kafka writer: ", 0),
	}
	src.Organization.IngestFile(config)

	// Until we have a fresh topic for testing,
	// Not sure yet how to do assertions as consumer will have to read everything
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddrs,
		Topic:   "organization",
	})
	ctx := context.Background()
	for i := 0; i < 2; i++ {
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

func TestDummy(t *testing.T) {
	got := 1
	if got != 1 {
		t.Errorf("got = %d; want 1", got)
	}
}
