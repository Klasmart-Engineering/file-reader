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

type Organization struct {
	Payload OrganizationPayload `json:"payload"`

	Metadata OrganizationMetadata `json:"metadata"`
}

const OrganizationAvroCRC64Fingerprint = " !\xac\x1ciנi"

func NewOrganization() Organization {
	r := Organization{}
	r.Payload = NewOrganizationPayload()

	r.Metadata = NewOrganizationMetadata()

	return r
}

func DeserializeOrganization(r io.Reader) (Organization, error) {
	t := NewOrganization()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeOrganizationFromSchema(r io.Reader, schema string) (Organization, error) {
	t := NewOrganization()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeOrganization(r Organization, w io.Writer) error {
	var err error
	err = writeOrganizationPayload(r.Payload, w)
	if err != nil {
		return err
	}
	err = writeOrganizationMetadata(r.Metadata, w)
	if err != nil {
		return err
	}
	return err
}

func (r Organization) Serialize(w io.Writer) error {
	return writeOrganization(r, w)
}

func (r Organization) Schema() string {
	return "{\"fields\":[{\"name\":\"payload\",\"type\":{\"fields\":[{\"name\":\"name\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"owner_user_uuid\",\"type\":\"string\"}],\"name\":\"OrganizationPayload\",\"type\":\"record\"}},{\"name\":\"metadata\",\"type\":{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_uuid\",\"type\":\"string\"}],\"name\":\"OrganizationMetadata\",\"type\":\"record\"}}],\"name\":\"com.kidsloop.onboarding.Organization\",\"type\":\"record\"}"
}

func (r Organization) SchemaName() string {
	return "com.kidsloop.onboarding.Organization"
}

func (_ Organization) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ Organization) SetInt(v int32)       { panic("Unsupported operation") }
func (_ Organization) SetLong(v int64)      { panic("Unsupported operation") }
func (_ Organization) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ Organization) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ Organization) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ Organization) SetString(v string)   { panic("Unsupported operation") }
func (_ Organization) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Organization) Get(i int) types.Field {
	switch i {
	case 0:
		r.Payload = NewOrganizationPayload()

		w := types.Record{Target: &r.Payload}

		return w

	case 1:
		r.Metadata = NewOrganizationMetadata()

		w := types.Record{Target: &r.Metadata}

		return w

	}
	panic("Unknown field index")
}

func (r *Organization) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Organization) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ Organization) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ Organization) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ Organization) HintSize(int)                     { panic("Unsupported operation") }
func (_ Organization) Finalize()                        {}

func (_ Organization) AvroCRC64Fingerprint() []byte {
	return []byte(OrganizationAvroCRC64Fingerprint)
}

func (r Organization) MarshalJSON() ([]byte, error) {
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

func (r *Organization) UnmarshalJSON(data []byte) error {
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
