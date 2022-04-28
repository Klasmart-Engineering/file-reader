package src

import (
	"context"
	avro "file_reader/avro_gencode"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

const (
	topic      = "organisation"
	schemaPath = "../src/avros/organisation.avsc"
)

func OrgCsvIngester(brokerAddrs []string) func(r io.Reader, ctx context.Context) {
	// schema logic will be part of build process
	// and this will be replaced with a cache lookup
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	schemaRegistryClient := srclient.CreateSchemaRegistryClient("http://localhost:8081")
	schema, err := schemaRegistryClient.CreateSchema(topic, string(schemaBytes), "AVRO")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Schema '%d' retrieved successfully!\n", schema.ID())
	// Prepare kafka producer for the provided topic
	logger := log.New(os.Stdout, "kafka writer: ", 0)
	w := kafka.Writer{
		Addr:   kafka.TCP(brokerAddrs...),
		Topic:  topic,
		Logger: logger,
	}
	return ComposeCsvIngester(
		topic,
		&w,
		rowToOrganisation,
		schema,
		brokerAddrs,
		logger,
	)
}

func rowToOrganisation(row []string) avro.Organization {
	md := avro.Metadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        uuid.NewString(),
	}
	pl := avro.Payload{
		Guid:              row[0],
		Organization_name: row[1],
	}
	return avro.Organization{Payload: pl, Metadata: md}
}
