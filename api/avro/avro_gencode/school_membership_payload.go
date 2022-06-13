// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     class.avsc
 *     organization_membership.avsc
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

type SchoolMembershipPayload struct {
	Name string `json:"name"`

	School_uuid string `json:"school_uuid"`

	User_uuid string `json:"user_uuid"`
}

const SchoolMembershipPayloadAvroCRC64Fingerprint = "S\xd6i|\xd8\x19\x94$"

func NewSchoolMembershipPayload() SchoolMembershipPayload {
	r := SchoolMembershipPayload{}
	return r
}

func DeserializeSchoolMembershipPayload(r io.Reader) (SchoolMembershipPayload, error) {
	t := NewSchoolMembershipPayload()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeSchoolMembershipPayloadFromSchema(r io.Reader, schema string) (SchoolMembershipPayload, error) {
	t := NewSchoolMembershipPayload()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeSchoolMembershipPayload(r SchoolMembershipPayload, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Name, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.School_uuid, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.User_uuid, w)
	if err != nil {
		return err
	}
	return err
}

func (r SchoolMembershipPayload) Serialize(w io.Writer) error {
	return writeSchoolMembershipPayload(r, w)
}

func (r SchoolMembershipPayload) Schema() string {
	return "{\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"school_uuid\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"user_uuid\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.SchoolMembershipPayload\",\"type\":\"record\"}"
}

func (r SchoolMembershipPayload) SchemaName() string {
	return "com.kidsloop.onboarding.SchoolMembershipPayload"
}

func (_ SchoolMembershipPayload) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetInt(v int32)       { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetLong(v int64)      { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetString(v string)   { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *SchoolMembershipPayload) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Name}

		return w

	case 1:
		w := types.String{Target: &r.School_uuid}

		return w

	case 2:
		w := types.String{Target: &r.User_uuid}

		return w

	}
	panic("Unknown field index")
}

func (r *SchoolMembershipPayload) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *SchoolMembershipPayload) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ SchoolMembershipPayload) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) HintSize(int)                     { panic("Unsupported operation") }
func (_ SchoolMembershipPayload) Finalize()                        {}

func (_ SchoolMembershipPayload) AvroCRC64Fingerprint() []byte {
	return []byte(SchoolMembershipPayloadAvroCRC64Fingerprint)
}

func (r SchoolMembershipPayload) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["name"], err = json.Marshal(r.Name)
	if err != nil {
		return nil, err
	}
	output["school_uuid"], err = json.Marshal(r.School_uuid)
	if err != nil {
		return nil, err
	}
	output["user_uuid"], err = json.Marshal(r.User_uuid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *SchoolMembershipPayload) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
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
		if v, ok := fields["school_uuid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.School_uuid); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for school_uuid")
	}
	val = func() json.RawMessage {
		if v, ok := fields["user_uuid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.User_uuid); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for user_uuid")
	}
	return nil
}