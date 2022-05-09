package src

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/riferrei/srclient"
)

var cachingEnabled = true
var fileSchemaCache = registerProtoSchemas(organizationProtoTopic)

type schemaRegistry struct {
	c *srclient.SchemaRegistryClient
}

var schemaRegistryClient = &schemaRegistry{
	c: srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
}

var protoSchema *srclient.Schema = nil

func cacheKey(fileName string, topic string) string {
	return fmt.Sprintf("%s-%s", fileName, topic)
}

func registerProtoSchemas(topic string) map[string]*srclient.Schema {
	var fileSchemaCache = make(map[string]*srclient.Schema)

	/*files, err := protoSchemaDir.ReadDir(os.Getenv("PROTO_SCHEMA_DIRECTORY"))
	fmt.Printf("files = %v", files)
	if err != nil {
		fmt.Printf("error here")
		log.Fatal(err)
	}*/

	files := []string{"onboarding.proto"}
	for _, file := range files {

		schemaPath := path.Join(os.Getenv("PROTO_SCHEMA_DIRECTORY"), file)
		schemaBytes, err := protoSchemaDir.ReadFile(schemaPath)

		if err != nil {
			panic(err)
		}

		schema, err := schemaRegistryClient.c.CreateSchema(topic, string(schemaBytes), srclient.Protobuf)
		// Cache the schema
		if cachingEnabled {
			cacheKey := cacheKey(file,
				topic)
			fileSchemaCache[cacheKey] = schema
		}

	}
	return fileSchemaCache
}
func GetSchemaIdBytes(schema *srclient.Schema) []byte {
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))
	return schemaIDBytes

}
func (client *schemaRegistry) getProtoSchema(schemaFileName string, topic string) *srclient.Schema {
	// First check if the schema is already cached
	cacheKey := cacheKey(schemaFileName,
		topic)
	if cachingEnabled {

		// Retrieve the schema from cache
		cachedSchema := fileSchemaCache[cacheKey]

		if cachedSchema != nil {
			return cachedSchema
		}

	}

	// Retrieve the lastest schema
	schema, err := schemaRegistryClient.c.GetLatestSchema(topic)

	// If it does not exist then create a new one and cache it
	if schema == nil || err != nil {

		schemaPath := path.Join(os.Getenv("PROTO_SCHEMA_DIRECTORY"), schemaFileName)

		schemaBytes, err := protoSchemaDir.ReadFile(schemaPath)

		if err != nil {
			panic(err)
		}
		schema, err = schemaRegistryClient.c.CreateSchema(topic, string(schemaBytes), srclient.Protobuf)
		if err != nil {
			panic(fmt.Sprintf("Error creating the schema: %s", err))
		}
	}
	// Cache the schema
	if cachingEnabled {
		fileSchemaCache[cacheKey] = schema
	}

	return schema

}
func (*schemaRegistry) getSchemaIdBytes(schemaFileName string, topic string) []byte {

	schemaPath := path.Join(os.Getenv("AVRO_SCHEMA_DIRECTORY"), schemaFileName)
	schemaBytes, err := avrosSchemaDir.ReadFile(schemaPath)
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
