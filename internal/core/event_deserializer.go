package core

import (
	"bytes"
	"encoding/binary"

	avrogen "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
)

func deserializeS3Event(schemaRegistry *SchemaRegistry, msgValue []byte) (avrogen.S3FileCreatedUpdated, error) {
	// Deserializes message bytes for S3 event to S3FileCreated
	schemaIdBytes := msgValue[1:5]
	schemaId := int(binary.BigEndian.Uint32(schemaIdBytes))
	schema, err := schemaRegistry.GetSchema(schemaId)
	if err != nil {
		return avrogen.S3FileCreatedUpdated{}, err
	}
	r := bytes.NewReader(msgValue[5:])
	s3FileCreated, err := avrogen.DeserializeS3FileCreatedUpdatedFromSchema(r, schema)
	if err != nil {
		return avrogen.S3FileCreatedUpdated{}, err
	}
	return s3FileCreated, nil
}
