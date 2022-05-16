package src

import (
	"encoding/binary"
	"log"

	"github.com/riferrei/srclient"
)

type SchemaRegistry struct {
	C           *srclient.SchemaRegistryClient
	IdSchemaMap map[int]string
}

func (SchemaRegistryClient *SchemaRegistry) GetSchemaIdBytes(schemaBody string, topic string) []byte {
	// Gets byte representation of id for the schema provided
	schema, err := SchemaRegistryClient.C.LookupSchema(topic, schemaBody, "AVRO")
	if schema == nil || err != nil {
		schema, err = SchemaRegistryClient.C.CreateSchema(topic, schemaBody, "AVRO")
		if err != nil {
			log.Fatal(err)
		}
	}

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	return schemaIDBytes
}

func (SchemaRegistryClient *SchemaRegistry) GetSchema(schemaId int) string {
	// Gets schema from local cache if exists, otherwise from schema registry
	if _, ok := SchemaRegistryClient.IdSchemaMap[schemaId]; !ok {
		schema, err := SchemaRegistryClient.C.GetSchema(schemaId)
		if err != nil {
			panic("could not get consumer schema " + err.Error())
		}
		SchemaRegistryClient.IdSchemaMap[schemaId] = schema.Schema()
	}
	return SchemaRegistryClient.IdSchemaMap[schemaId]
}
