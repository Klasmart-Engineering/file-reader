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

type ClassDetailsMetadata struct {
	Origin_application string `json:"origin_application"`

	Region string `json:"region"`

	Tracking_uuid string `json:"tracking_uuid"`
}

const ClassDetailsMetadataAvroCRC64Fingerprint = "\xc4\x1c\x01E\x86\x02[\xb8"

func NewClassDetailsMetadata() ClassDetailsMetadata {
	r := ClassDetailsMetadata{}
	return r
}

func DeserializeClassDetailsMetadata(r io.Reader) (ClassDetailsMetadata, error) {
	t := NewClassDetailsMetadata()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeClassDetailsMetadataFromSchema(r io.Reader, schema string) (ClassDetailsMetadata, error) {
	t := NewClassDetailsMetadata()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeClassDetailsMetadata(r ClassDetailsMetadata, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Origin_application, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Region, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Tracking_uuid, w)
	if err != nil {
		return err
	}
	return err
}

func (r ClassDetailsMetadata) Serialize(w io.Writer) error {
	return writeClassDetailsMetadata(r, w)
}

func (r ClassDetailsMetadata) Schema() string {
	return "{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_uuid\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.ClassDetailsMetadata\",\"type\":\"record\"}"
}

func (r ClassDetailsMetadata) SchemaName() string {
	return "com.kidsloop.onboarding.ClassDetailsMetadata"
}

func (_ ClassDetailsMetadata) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetInt(v int32)       { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetLong(v int64)      { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetString(v string)   { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *ClassDetailsMetadata) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Origin_application}

		return w

	case 1:
		w := types.String{Target: &r.Region}

		return w

	case 2:
		w := types.String{Target: &r.Tracking_uuid}

		return w

	}
	panic("Unknown field index")
}

func (r *ClassDetailsMetadata) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *ClassDetailsMetadata) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ ClassDetailsMetadata) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) HintSize(int)                     { panic("Unsupported operation") }
func (_ ClassDetailsMetadata) Finalize()                        {}

func (_ ClassDetailsMetadata) AvroCRC64Fingerprint() []byte {
	return []byte(ClassDetailsMetadataAvroCRC64Fingerprint)
}

func (r ClassDetailsMetadata) MarshalJSON() ([]byte, error) {
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
	output["tracking_uuid"], err = json.Marshal(r.Tracking_uuid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *ClassDetailsMetadata) UnmarshalJSON(data []byte) error {
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
		if v, ok := fields["tracking_uuid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Tracking_uuid); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for tracking_uuid")
	}
	return nil
}
