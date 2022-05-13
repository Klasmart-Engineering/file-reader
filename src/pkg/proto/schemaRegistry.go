package proto

import (
	"context"
	"encoding/binary"
	"file_reader/src/instrument"
	"file_reader/src/protos/onboarding"
	"file_reader/src/third_party/protobuf"
	"file_reader/src/third_party/protobuf/srclient"
	"fmt"

	"github.com/pkg/errors"
)

var cachingEnabled = true
var fileSchemaCache = registerProtoSchemas(organizationProtoTopic)

type schemaRegistry struct {
	c   srclient.Client
	ctx context.Context
}

var schemaRegistryClient = &schemaRegistry{
	c:   srclient.NewClient(srclient.WithURL(instrument.MustGetEnv("SCHEMA_CLIENT_ENDPOINT"))),
	ctx: context.Background(),
}

func cacheKey(schemaName string, topic string) string {
	return fmt.Sprintf("%s-%s", schemaName, topic)
}

func registerProtoSchemas(topic string) map[string]int {

	registrator := protobuf.NewSchemaRegistrator(schemaRegistryClient.c)

	schemaID, err := registrator.RegisterValue(schemaRegistryClient.ctx, topic, &onboarding.Organization{})

	if err != nil {
		panic(fmt.Errorf("error registering schema: %w", err))
	}

	var fileSchemaCache = make(map[string]int)

	// Cache the schema
	if cachingEnabled {
		cacheKey := cacheKey("organization",
			topic)
		fileSchemaCache[cacheKey] = schemaID
	}
	return fileSchemaCache
}
func GetSchemaIdBytes(schemaID int) []byte {
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schemaID))
	return schemaIDBytes

}
func (client *schemaRegistry) GetProtoSchemaID(schemaName string, topic string) (int, error) {
	// First check if the schema is already cached
	cacheKey := cacheKey(schemaName,
		topic)
	if cachingEnabled {

		// Retrieve the schema from cache
		if schemaID, ok := fileSchemaCache[cacheKey]; ok {
			return schemaID, nil
		}

	}

	// Retrieve the lastest schema
	schema, err := schemaRegistryClient.c.GetLatestSchema(schemaRegistryClient.ctx, topic)

	// If it does not exist then register a new one and cache it

	if schema == nil || err != nil {

		registrator := protobuf.NewSchemaRegistrator(schemaRegistryClient.c)

		schemaID, err := registrator.RegisterValue(schemaRegistryClient.ctx, topic, &onboarding.Organization{})

		if err != nil {
			err = errors.Wrap(err, "error registering schema")
			return -1, err
		}
		// Cache the schema
		if cachingEnabled {
			fileSchemaCache[cacheKey] = schemaID
		}
		return schemaID, nil
	}
	// Cache the schema
	if cachingEnabled {
		fileSchemaCache[cacheKey] = schema.ID
	}

	return schema.ID, nil

}
