package proto

import (
	"context"
	"encoding/binary"
	"log"

	"github.com/KL-Engineering/file-reader/cmd/instrument"
	"github.com/KL-Engineering/file-reader/cmd/protos/onboarding"
	"github.com/KL-Engineering/file-reader/cmd/third_party/protobuf"
	"github.com/KL-Engineering/file-reader/cmd/third_party/protobuf/srclient"
)

var cachingEnabled = true

type SchemaRegistry struct {
	c           srclient.Client
	ctx         context.Context
	IdSchemaMap map[int]string
}

var SchemaRegistryClient = &SchemaRegistry{
	c:           srclient.NewClient(srclient.WithURL(instrument.MustGetEnv("SCHEMA_CLIENT_ENDPOINT"))),
	ctx:         context.Background(),
	IdSchemaMap: make(map[int]string),
}

func GetSchemaRegistryClient() *SchemaRegistry {
	return SchemaRegistryClient
}
func (client *SchemaRegistry) GetSchemaID(topic string) int {

	// Retrieve the lastest schema
	schema, err := SchemaRegistryClient.c.GetLatestSchema(SchemaRegistryClient.ctx, topic)

	// If it does not exist then register a new one and cache it

	if schema == nil || err != nil {

		registrator := protobuf.NewSchemaRegistrator(SchemaRegistryClient.c)

		schema, err = registrator.RegisterValue(SchemaRegistryClient.ctx, topic, &onboarding.Organization{})

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
		schema, err := client.c.GetSchemaByID(SchemaRegistryClient.ctx, schemaId)
		if err != nil {
			panic("could not get consumer schema " + err.Error())
		}
		if cachingEnabled {

			client.IdSchemaMap[schemaId] = schema.Schema
		}
	}
	return client.IdSchemaMap[schemaId]
}
