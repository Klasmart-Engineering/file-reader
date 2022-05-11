package src

import (
	"encoding/binary"
	"log"

	"github.com/riferrei/srclient"
)

type SchemaRegistry struct {
	C *srclient.SchemaRegistryClient
}

func (SchemaRegistryClient *SchemaRegistry) GetSchemaIdBytes(schemaBody string, topic string) []byte {
	schema, err := SchemaRegistryClient.C.CreateSchema(topic, schemaBody, "AVRO")
	if err != nil {
		log.Fatal(err)
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	return schemaIDBytes
}
