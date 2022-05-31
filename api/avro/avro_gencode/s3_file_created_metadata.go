// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     s3filecreated.avsc
 */
package avro

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

var _ = fmt.Printf

type S3FileCreatedMetadata struct {
	Origin_application string `json:"origin_application"`

	Region string `json:"region"`

	Tracking_id string `json:"tracking_id"`
}

const S3FileCreatedMetadataAvroCRC64Fingerprint = "\an\x84p1\xdd\xd5\""

func NewS3FileCreatedMetadata() S3FileCreatedMetadata {
	r := S3FileCreatedMetadata{}
	return r
}

func DeserializeS3FileCreatedMetadata(r io.Reader) (S3FileCreatedMetadata, error) {
	t := NewS3FileCreatedMetadata()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeS3FileCreatedMetadataFromSchema(r io.Reader, schema string) (S3FileCreatedMetadata, error) {
	t := NewS3FileCreatedMetadata()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeS3FileCreatedMetadata(r S3FileCreatedMetadata, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Origin_application, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Region, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Tracking_id, w)
	if err != nil {
		return err
	}
	return err
}

func (r S3FileCreatedMetadata) Serialize(w io.Writer) error {
	return writeS3FileCreatedMetadata(r, w)
}

func (r S3FileCreatedMetadata) Schema() string {
	return "{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_id\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.S3FileCreatedMetadata\",\"type\":\"record\"}"
}

func (r S3FileCreatedMetadata) SchemaName() string {
	return "com.kidsloop.onboarding.S3FileCreatedMetadata"
}

func (_ S3FileCreatedMetadata) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetInt(v int32)       { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetLong(v int64)      { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetString(v string)   { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *S3FileCreatedMetadata) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Origin_application}

		return w

	case 1:
		w := types.String{Target: &r.Region}

		return w

	case 2:
		w := types.String{Target: &r.Tracking_id}

		return w

	}
	panic("Unknown field index")
}

func (r *S3FileCreatedMetadata) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *S3FileCreatedMetadata) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ S3FileCreatedMetadata) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) HintSize(int)                     { panic("Unsupported operation") }
func (_ S3FileCreatedMetadata) Finalize()                        {}

func (_ S3FileCreatedMetadata) AvroCRC64Fingerprint() []byte {
	return []byte(S3FileCreatedMetadataAvroCRC64Fingerprint)
}

func (r S3FileCreatedMetadata) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["origin_application"], err = json.Marshal(r.Origin_application)
	if err != nil {
		return nil, err
	}
	output["region"], err = json.Marshal(r.Region)
	if err != nil {
		return nil, err
	}
	output["tracking_id"], err = json.Marshal(r.Tracking_id)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *S3FileCreatedMetadata) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["origin_application"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Origin_application); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for origin_application")
	}
	val = func() json.RawMessage {
		if v, ok := fields["region"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Region); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for region")
	}
	val = func() json.RawMessage {
		if v, ok := fields["tracking_id"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Tracking_id); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for tracking_id")
	}
	return nil
}
