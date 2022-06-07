// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     class.avsc
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

type SchoolPayload struct {
	Uuid *UnionNullString `json:"uuid"`

	Organization_id string `json:"organization_id"`

	Name string `json:"name"`

	Program_ids *UnionNullArrayString `json:"program_ids"`
}

const SchoolPayloadAvroCRC64Fingerprint = "u8x8V\xb4\xd5V"

func NewSchoolPayload() SchoolPayload {
	r := SchoolPayload{}
	r.Uuid = nil
	r.Program_ids = nil
	return r
}

func DeserializeSchoolPayload(r io.Reader) (SchoolPayload, error) {
	t := NewSchoolPayload()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeSchoolPayloadFromSchema(r io.Reader, schema string) (SchoolPayload, error) {
	t := NewSchoolPayload()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeSchoolPayload(r SchoolPayload, w io.Writer) error {
	var err error
	err = writeUnionNullString(r.Uuid, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Organization_id, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Name, w)
	if err != nil {
		return err
	}
	err = writeUnionNullArrayString(r.Program_ids, w)
	if err != nil {
		return err
	}
	return err
}

func (r SchoolPayload) Serialize(w io.Writer) error {
	return writeSchoolPayload(r, w)
}

func (r SchoolPayload) Schema() string {
	return "{\"fields\":[{\"default\":null,\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":[\"null\",\"string\"]},{\"logicalType\":\"uuid\",\"name\":\"organization_id\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"},{\"default\":null,\"name\":\"program_ids\",\"type\":[\"null\",{\"items\":\"string\",\"type\":\"array\"}]}],\"name\":\"com.kidsloop.onboarding.SchoolPayload\",\"type\":\"record\"}"
}

func (r SchoolPayload) SchemaName() string {
	return "com.kidsloop.onboarding.SchoolPayload"
}

func (_ SchoolPayload) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ SchoolPayload) SetInt(v int32)       { panic("Unsupported operation") }
func (_ SchoolPayload) SetLong(v int64)      { panic("Unsupported operation") }
func (_ SchoolPayload) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ SchoolPayload) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ SchoolPayload) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ SchoolPayload) SetString(v string)   { panic("Unsupported operation") }
func (_ SchoolPayload) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *SchoolPayload) Get(i int) types.Field {
	switch i {
	case 0:
		r.Uuid = NewUnionNullString()

		return r.Uuid
	case 1:
		w := types.String{Target: &r.Organization_id}

		return w

	case 2:
		w := types.String{Target: &r.Name}

		return w

	case 3:
		r.Program_ids = NewUnionNullArrayString()

		return r.Program_ids
	}
	panic("Unknown field index")
}

func (r *SchoolPayload) SetDefault(i int) {
	switch i {
	case 0:
		r.Uuid = nil
		return
	case 3:
		r.Program_ids = nil
		return
	}
	panic("Unknown field index")
}

func (r *SchoolPayload) NullField(i int) {
	switch i {
	case 0:
		r.Uuid = nil
		return
	case 3:
		r.Program_ids = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ SchoolPayload) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ SchoolPayload) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ SchoolPayload) HintSize(int)                     { panic("Unsupported operation") }
func (_ SchoolPayload) Finalize()                        {}

func (_ SchoolPayload) AvroCRC64Fingerprint() []byte {
	return []byte(SchoolPayloadAvroCRC64Fingerprint)
}

func (r SchoolPayload) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["uuid"], err = json.Marshal(r.Uuid)
	if err != nil {
		return nil, err
	}
	output["organization_id"], err = json.Marshal(r.Organization_id)
	if err != nil {
		return nil, err
	}
	output["name"], err = json.Marshal(r.Name)
	if err != nil {
		return nil, err
	}
	output["program_ids"], err = json.Marshal(r.Program_ids)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *SchoolPayload) UnmarshalJSON(data []byte) error {
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
		if v, ok := fields["organization_id"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Organization_id); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for organization_id")
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
	val = func() json.RawMessage {
		if v, ok := fields["program_ids"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Program_ids); err != nil {
			return err
		}
	} else {
		r.Program_ids = NewUnionNullArrayString()

		r.Program_ids = nil
	}
	return nil
}
