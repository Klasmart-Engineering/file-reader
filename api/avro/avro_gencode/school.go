// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     class.avsc
 *     organization_membership.avsc
 *     class_details.avsc
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

type School struct {
	Payload SchoolPayload `json:"payload"`

	Metadata SchoolMetadata `json:"metadata"`
}

const SchoolAvroCRC64Fingerprint = "\xbd\xdc\xd7\xeb\xfd \xa1p"

func NewSchool() School {
	r := School{}
	r.Payload = NewSchoolPayload()

	r.Metadata = NewSchoolMetadata()

	return r
}

func DeserializeSchool(r io.Reader) (School, error) {
	t := NewSchool()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeSchoolFromSchema(r io.Reader, schema string) (School, error) {
	t := NewSchool()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeSchool(r School, w io.Writer) error {
	var err error
	err = writeSchoolPayload(r.Payload, w)
	if err != nil {
		return err
	}
	err = writeSchoolMetadata(r.Metadata, w)
	if err != nil {
		return err
	}
	return err
}

func (r School) Serialize(w io.Writer) error {
	return writeSchool(r, w)
}

func (r School) Schema() string {
	return "{\"fields\":[{\"name\":\"payload\",\"type\":{\"fields\":[{\"default\":null,\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":[\"null\",\"string\"]},{\"logicalType\":\"uuid\",\"name\":\"organization_uuid\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"},{\"default\":null,\"name\":\"program_uuids\",\"type\":[\"null\",{\"items\":\"string\",\"type\":\"array\"}]}],\"name\":\"SchoolPayload\",\"type\":\"record\"}},{\"name\":\"metadata\",\"type\":{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_uuid\",\"type\":\"string\"}],\"name\":\"SchoolMetadata\",\"type\":\"record\"}}],\"name\":\"com.kidsloop.onboarding.School\",\"type\":\"record\"}"
}

func (r School) SchemaName() string {
	return "com.kidsloop.onboarding.School"
}

func (_ School) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ School) SetInt(v int32)       { panic("Unsupported operation") }
func (_ School) SetLong(v int64)      { panic("Unsupported operation") }
func (_ School) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ School) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ School) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ School) SetString(v string)   { panic("Unsupported operation") }
func (_ School) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *School) Get(i int) types.Field {
	switch i {
	case 0:
		r.Payload = NewSchoolPayload()

		w := types.Record{Target: &r.Payload}

		return w

	case 1:
		r.Metadata = NewSchoolMetadata()

		w := types.Record{Target: &r.Metadata}

		return w

	}
	panic("Unknown field index")
}

func (r *School) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *School) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ School) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ School) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ School) HintSize(int)                     { panic("Unsupported operation") }
func (_ School) Finalize()                        {}

func (_ School) AvroCRC64Fingerprint() []byte {
	return []byte(SchoolAvroCRC64Fingerprint)
}

func (r School) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["payload"], err = json.Marshal(r.Payload)
	if err != nil {
		return nil, err
	}
	output["metadata"], err = json.Marshal(r.Metadata)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *School) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["payload"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Payload); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for payload")
	}
	val = func() json.RawMessage {
		if v, ok := fields["metadata"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Metadata); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for metadata")
	}
	return nil
}
