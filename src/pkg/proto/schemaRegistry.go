package proto

import (
	"context"
	"encoding/binary"
	"file_reader/src/instrument"
	"file_reader/src/protos/onboarding"
	"file_reader/src/third_party/protobuf"
	"file_reader/src/third_party/protobuf/srclient"
	"log"
)

var cachingEnabled = true

type SchemaRegistry struct {
	c           srclient.Client
	ctx         context.Context
	IdSchemaMap map[int]string
}

var schemaRegistryClient = &SchemaRegistry{
	c:           srclient.NewClient(srclient.WithURL(instrument.MustGetEnv("SCHEMA_CLIENT_ENDPOINT"))),
	ctx:         context.Background(),
	IdSchemaMap: make(map[int]string),
}

func GetSchemaRegistryClient() *SchemaRegistry {
	return schemaRegistryClient
}
func (client *SchemaRegistry) GetSchemaID(topic string) int {

	// Retrieve the lastest schema
	schema, err := schemaRegistryClient.c.GetLatestSchema(schemaRegistryClient.ctx, topic)

	// If it does not exist then register a new one and cache it

	if schema == nil || err != nil {

		registrator := protobuf.NewSchemaRegistrator(schemaRegistryClient.c)

		schema, err = registrator.RegisterValue(schemaRegistryClient.ctx, topic, &onboarding.Organization{})

		if err != nil {
			log.Fatal(err)
		}
	}
	return schema.ID

}

func (client *SchemaRegistry) GetSchemaIdBytes(schemaID int) []byte {

	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schemaID))
	return schemaIDBytes
}

func (client *SchemaRegistry) GetSchema(schemaId int) string {
	// Gets schema from local cache if exists, otherwise from schema registry
	if _, ok := client.IdSchemaMap[schemaId]; !ok {
		schema, err := client.c.GetSchemaByID(schemaRegistryClient.ctx, schemaId)
		if err != nil {
			panic("could not get consumer schema " + err.Error())
		}
		if cachingEnabled {

			client.IdSchemaMap[schemaId] = schema.Schema
		}
	}
	return client.IdSchemaMap[schemaId]
}
