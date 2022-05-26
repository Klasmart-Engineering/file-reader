package core

import (
	"bytes"
	"encoding/binary"

	avrogen "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
)

func deserializeS3Event(schemaRegistry *SchemaRegistry, msgValue []byte) (avrogen.S3FileCreated, error) {
	// Deserializes message bytes for S3 event to S3FileCreated
	schemaIdBytes := msgValue[1:5]
	schemaId := int(binary.BigEndian.Uint32(schemaIdBytes))
	schema, err := schemaRegistry.GetSchema(schemaId)
	if err != nil {
		return avrogen.S3FileCreated{}, err
	}
	r := bytes.NewReader(msgValue[5:])
	s3FileCreated, err := avrogen.DeserializeS3FileCreatedFromSchema(r, schema)
	if err != nil {
		return avrogen.S3FileCreated{}, err
	}
	return s3FileCreated, nil
}
