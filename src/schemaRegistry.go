package src

import (
	"encoding/binary"
	"log"
	"os"
	"path"

	"github.com/riferrei/srclient"
)

type schemaRegistry struct {
	c *srclient.SchemaRegistryClient
}

var schemaRegistryClient = &schemaRegistry{
	c: srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
}

func (*schemaRegistry) getSchemaIdBytes(schemaFileName string, topic string) []byte {
	schemaPath := path.Join(os.Getenv("AVRO_SCHEMA_DIRECTORY"), schemaFileName)
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	schema, err := schemaRegistryClient.c.CreateSchema(topic, string(schemaBytes), "AVRO")
	if err != nil {
		log.Fatal(err)
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	return schemaIDBytes
}
