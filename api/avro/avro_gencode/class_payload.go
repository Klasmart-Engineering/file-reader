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

type ClassPayload struct {
	Uuid *UnionNullString `json:"uuid"`

	Organization_uuid string `json:"organization_uuid"`

	Name string `json:"name"`
}

const ClassPayloadAvroCRC64Fingerprint = "\xf3\xf2\xa0Cz\x16\x03\xa7"

func NewClassPayload() ClassPayload {
	r := ClassPayload{}
	r.Uuid = nil
	return r
}

func DeserializeClassPayload(r io.Reader) (ClassPayload, error) {
	t := NewClassPayload()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeClassPayloadFromSchema(r io.Reader, schema string) (ClassPayload, error) {
	t := NewClassPayload()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeClassPayload(r ClassPayload, w io.Writer) error {
	var err error
	err = writeUnionNullString(r.Uuid, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Organization_uuid, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Name, w)
	if err != nil {
		return err
	}
	return err
}

func (r ClassPayload) Serialize(w io.Writer) error {
	return writeClassPayload(r, w)
}

func (r ClassPayload) Schema() string {
	return "{\"fields\":[{\"default\":null,\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":[\"null\",\"string\"]},{\"logicalType\":\"uuid\",\"name\":\"organization_uuid\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.ClassPayload\",\"type\":\"record\"}"
}

func (r ClassPayload) SchemaName() string {
	return "com.kidsloop.onboarding.ClassPayload"
}

func (_ ClassPayload) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ ClassPayload) SetInt(v int32)       { panic("Unsupported operation") }
func (_ ClassPayload) SetLong(v int64)      { panic("Unsupported operation") }
func (_ ClassPayload) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ ClassPayload) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ ClassPayload) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ ClassPayload) SetString(v string)   { panic("Unsupported operation") }
func (_ ClassPayload) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *ClassPayload) Get(i int) types.Field {
	switch i {
	case 0:
		r.Uuid = NewUnionNullString()

		return r.Uuid
	case 1:
		w := types.String{Target: &r.Organization_uuid}

		return w

	case 2:
		w := types.String{Target: &r.Name}

		return w

	}
	panic("Unknown field index")
}

func (r *ClassPayload) SetDefault(i int) {
	switch i {
	case 0:
		r.Uuid = nil
		return
	}
	panic("Unknown field index")
}

func (r *ClassPayload) NullField(i int) {
	switch i {
	case 0:
		r.Uuid = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ ClassPayload) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ ClassPayload) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ ClassPayload) HintSize(int)                     { panic("Unsupported operation") }
func (_ ClassPayload) Finalize()                        {}

func (_ ClassPayload) AvroCRC64Fingerprint() []byte {
	return []byte(ClassPayloadAvroCRC64Fingerprint)
}

func (r ClassPayload) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["uuid"], err = json.Marshal(r.Uuid)
	if err != nil {
		return nil, err
	}
	output["organization_uuid"], err = json.Marshal(r.Organization_uuid)
	if err != nil {
		return nil, err
	}
	output["name"], err = json.Marshal(r.Name)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *ClassPayload) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["uuid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Uuid); err != nil {
			return err
		}
	} else {
		r.Uuid = NewUnionNullString()

		r.Uuid = nil
	}
	val = func() json.RawMessage {
		if v, ok := fields["organization_uuid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Organization_uuid); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for organization_uuid")
	}
	val = func() json.RawMessage {
		if v, ok := fields["name"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Name); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for name")
	}
	return nil
}
