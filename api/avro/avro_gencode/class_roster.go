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

type ClassRoster struct {
	Payload ClassRosterPayload `json:"payload"`

	Metadata ClassRosterMetadata `json:"metadata"`
}

const ClassRosterAvroCRC64Fingerprint = "\x8b\x05\xc9\xdb\xd9\x00\xda\xcd"

func NewClassRoster() ClassRoster {
	r := ClassRoster{}
	r.Payload = NewClassRosterPayload()

	r.Metadata = NewClassRosterMetadata()

	return r
}

func DeserializeClassRoster(r io.Reader) (ClassRoster, error) {
	t := NewClassRoster()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeClassRosterFromSchema(r io.Reader, schema string) (ClassRoster, error) {
	t := NewClassRoster()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeClassRoster(r ClassRoster, w io.Writer) error {
	var err error
	err = writeClassRosterPayload(r.Payload, w)
	if err != nil {
		return err
	}
	err = writeClassRosterMetadata(r.Metadata, w)
	if err != nil {
		return err
	}
	return err
}

func (r ClassRoster) Serialize(w io.Writer) error {
	return writeClassRoster(r, w)
}

func (r ClassRoster) Schema() string {
	return "{\"fields\":[{\"name\":\"payload\",\"type\":{\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"class_uuid\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"user_uuid\",\"type\":\"string\"},{\"name\":\"participating_as\",\"type\":\"string\"}],\"name\":\"ClassRosterPayload\",\"type\":\"record\"}},{\"name\":\"metadata\",\"type\":{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_uuid\",\"type\":\"string\"}],\"name\":\"ClassRosterMetadata\",\"type\":\"record\"}}],\"name\":\"com.kidsloop.onboarding.ClassRoster\",\"type\":\"record\"}"
}

func (r ClassRoster) SchemaName() string {
	return "com.kidsloop.onboarding.ClassRoster"
}

func (_ ClassRoster) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ ClassRoster) SetInt(v int32)       { panic("Unsupported operation") }
func (_ ClassRoster) SetLong(v int64)      { panic("Unsupported operation") }
func (_ ClassRoster) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ ClassRoster) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ ClassRoster) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ ClassRoster) SetString(v string)   { panic("Unsupported operation") }
func (_ ClassRoster) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *ClassRoster) Get(i int) types.Field {
	switch i {
	case 0:
		r.Payload = NewClassRosterPayload()

		w := types.Record{Target: &r.Payload}

		return w

	case 1:
		r.Metadata = NewClassRosterMetadata()

		w := types.Record{Target: &r.Metadata}

		return w

	}
	panic("Unknown field index")
}

func (r *ClassRoster) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *ClassRoster) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ ClassRoster) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ ClassRoster) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ ClassRoster) HintSize(int)                     { panic("Unsupported operation") }
func (_ ClassRoster) Finalize()                        {}

func (_ ClassRoster) AvroCRC64Fingerprint() []byte {
	return []byte(ClassRosterAvroCRC64Fingerprint)
}

func (r ClassRoster) MarshalJSON() ([]byte, error) {
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

func (r *ClassRoster) UnmarshalJSON(data []byte) error {
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
