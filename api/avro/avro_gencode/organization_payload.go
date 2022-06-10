// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     class.avsc
 *     organization_membership.avsc
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

type OrganizationPayload struct {
	Name string `json:"name"`

	Uuid string `json:"uuid"`

	Owner_user_uuid string `json:"owner_user_uuid"`
}

const OrganizationPayloadAvroCRC64Fingerprint = "\xe1&\x9f\n\xf6\xa0\xb3w"

func NewOrganizationPayload() OrganizationPayload {
	r := OrganizationPayload{}
	return r
}

func DeserializeOrganizationPayload(r io.Reader) (OrganizationPayload, error) {
	t := NewOrganizationPayload()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeOrganizationPayloadFromSchema(r io.Reader, schema string) (OrganizationPayload, error) {
	t := NewOrganizationPayload()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeOrganizationPayload(r OrganizationPayload, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Name, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Uuid, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Owner_user_uuid, w)
	if err != nil {
		return err
	}
	return err
}

func (r OrganizationPayload) Serialize(w io.Writer) error {
	return writeOrganizationPayload(r, w)
}

func (r OrganizationPayload) Schema() string {
	return "{\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"owner_user_uuid\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.OrganizationPayload\",\"type\":\"record\"}"
}

func (r OrganizationPayload) SchemaName() string {
	return "com.kidsloop.onboarding.OrganizationPayload"
}

func (_ OrganizationPayload) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ OrganizationPayload) SetInt(v int32)       { panic("Unsupported operation") }
func (_ OrganizationPayload) SetLong(v int64)      { panic("Unsupported operation") }
func (_ OrganizationPayload) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ OrganizationPayload) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ OrganizationPayload) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ OrganizationPayload) SetString(v string)   { panic("Unsupported operation") }
func (_ OrganizationPayload) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *OrganizationPayload) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Name}

		return w

	case 1:
		w := types.String{Target: &r.Uuid}

		return w

	case 2:
		w := types.String{Target: &r.Owner_user_uuid}

		return w

	}
	panic("Unknown field index")
}

func (r *OrganizationPayload) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *OrganizationPayload) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ OrganizationPayload) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ OrganizationPayload) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ OrganizationPayload) HintSize(int)                     { panic("Unsupported operation") }
func (_ OrganizationPayload) Finalize()                        {}

func (_ OrganizationPayload) AvroCRC64Fingerprint() []byte {
	return []byte(OrganizationPayloadAvroCRC64Fingerprint)
}

func (r OrganizationPayload) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["name"], err = json.Marshal(r.Name)
	if err != nil {
		return nil, err
	}
	output["uuid"], err = json.Marshal(r.Uuid)
	if err != nil {
		return nil, err
	}
	output["owner_user_uuid"], err = json.Marshal(r.Owner_user_uuid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *OrganizationPayload) UnmarshalJSON(data []byte) error {
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
		return fmt.Errorf("no value specified for uuid")
	}
	val = func() json.RawMessage {
		if v, ok := fields["owner_user_uuid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Owner_user_uuid); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for owner_user_uuid")
	}
	return nil
}
