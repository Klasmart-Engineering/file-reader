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

type OrganizationMetadata struct {
	Origin_application string `json:"origin_application"`

	Region string `json:"region"`

	Tracking_id string `json:"tracking_id"`
}

const OrganizationMetadataAvroCRC64Fingerprint = "9R\xbf\x96Į\xff\xf7"

func NewOrganizationMetadata() OrganizationMetadata {
	r := OrganizationMetadata{}
	return r
}

func DeserializeOrganizationMetadata(r io.Reader) (OrganizationMetadata, error) {
	t := NewOrganizationMetadata()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeOrganizationMetadataFromSchema(r io.Reader, schema string) (OrganizationMetadata, error) {
	t := NewOrganizationMetadata()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeOrganizationMetadata(r OrganizationMetadata, w io.Writer) error {
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

func (r OrganizationMetadata) Serialize(w io.Writer) error {
	return writeOrganizationMetadata(r, w)
}

func (r OrganizationMetadata) Schema() string {
	return "{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_id\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.OrganizationMetadata\",\"type\":\"record\"}"
}

func (r OrganizationMetadata) SchemaName() string {
	return "com.kidsloop.onboarding.OrganizationMetadata"
}

func (_ OrganizationMetadata) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetInt(v int32)       { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetLong(v int64)      { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetString(v string)   { panic("Unsupported operation") }
func (_ OrganizationMetadata) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *OrganizationMetadata) Get(i int) types.Field {
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

func (r *OrganizationMetadata) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *OrganizationMetadata) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ OrganizationMetadata) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ OrganizationMetadata) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ OrganizationMetadata) HintSize(int)                     { panic("Unsupported operation") }
func (_ OrganizationMetadata) Finalize()                        {}

func (_ OrganizationMetadata) AvroCRC64Fingerprint() []byte {
	return []byte(OrganizationMetadataAvroCRC64Fingerprint)
}

func (r OrganizationMetadata) MarshalJSON() ([]byte, error) {
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

func (r *OrganizationMetadata) UnmarshalJSON(data []byte) error {
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
