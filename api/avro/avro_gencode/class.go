// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     class.avsc
 *     organization_membership.avsc
 *     class_details.avsc
 *     class_roster.avsc
 *     school_membership.avsc
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

type Class struct {
	Payload ClassPayload `json:"payload"`

	Metadata ClassMetadata `json:"metadata"`
}

const ClassAvroCRC64Fingerprint = "\x19E\x9d\x16H\x83\xa6\xd4"

func NewClass() Class {
	r := Class{}
	r.Payload = NewClassPayload()

	r.Metadata = NewClassMetadata()

	return r
}

func DeserializeClass(r io.Reader) (Class, error) {
	t := NewClass()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeClassFromSchema(r io.Reader, schema string) (Class, error) {
	t := NewClass()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeClass(r Class, w io.Writer) error {
	var err error
	err = writeClassPayload(r.Payload, w)
	if err != nil {
		return err
	}
	err = writeClassMetadata(r.Metadata, w)
	if err != nil {
		return err
	}
	return err
}

func (r Class) Serialize(w io.Writer) error {
	return writeClass(r, w)
}

func (r Class) Schema() string {
	return "{\"fields\":[{\"name\":\"payload\",\"type\":{\"fields\":[{\"default\":null,\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":[\"null\",\"string\"]},{\"logicalType\":\"uuid\",\"name\":\"organization_uuid\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"ClassPayload\",\"type\":\"record\"}},{\"name\":\"metadata\",\"type\":{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_uuid\",\"type\":\"string\"}],\"name\":\"ClassMetadata\",\"type\":\"record\"}}],\"name\":\"com.kidsloop.onboarding.Class\",\"type\":\"record\"}"
}

func (r Class) SchemaName() string {
	return "com.kidsloop.onboarding.Class"
}

func (_ Class) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ Class) SetInt(v int32)       { panic("Unsupported operation") }
func (_ Class) SetLong(v int64)      { panic("Unsupported operation") }
func (_ Class) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ Class) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ Class) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ Class) SetString(v string)   { panic("Unsupported operation") }
func (_ Class) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Class) Get(i int) types.Field {
	switch i {
	case 0:
		r.Payload = NewClassPayload()

		w := types.Record{Target: &r.Payload}

		return w

	case 1:
		r.Metadata = NewClassMetadata()

		w := types.Record{Target: &r.Metadata}

		return w

	}
	panic("Unknown field index")
}

func (r *Class) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Class) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ Class) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ Class) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ Class) HintSize(int)                     { panic("Unsupported operation") }
func (_ Class) Finalize()                        {}

func (_ Class) AvroCRC64Fingerprint() []byte {
	return []byte(ClassAvroCRC64Fingerprint)
}

func (r Class) MarshalJSON() ([]byte, error) {
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

func (r *Class) UnmarshalJSON(data []byte) error {
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
