// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
 *     class.avsc
 *     organization_membership.avsc
 *     school_membership.avsc
 *     class_roster.avsc
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

type UserMetadata struct {
	Origin_application string `json:"origin_application"`

	Region string `json:"region"`

	Tracking_uuid string `json:"tracking_uuid"`
}

const UserMetadataAvroCRC64Fingerprint = "*mg\xce\xc9\xffX\xdb"

func NewUserMetadata() UserMetadata {
	r := UserMetadata{}
	return r
}

func DeserializeUserMetadata(r io.Reader) (UserMetadata, error) {
	t := NewUserMetadata()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeUserMetadataFromSchema(r io.Reader, schema string) (UserMetadata, error) {
	t := NewUserMetadata()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeUserMetadata(r UserMetadata, w io.Writer) error {
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

func (r UserMetadata) Serialize(w io.Writer) error {
	return writeUserMetadata(r, w)
}

func (r UserMetadata) Schema() string {
	return "{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_uuid\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.user.UserMetadata\",\"type\":\"record\"}"
}

func (r UserMetadata) SchemaName() string {
	return "com.kidsloop.onboarding.user.UserMetadata"
}

func (_ UserMetadata) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ UserMetadata) SetInt(v int32)       { panic("Unsupported operation") }
func (_ UserMetadata) SetLong(v int64)      { panic("Unsupported operation") }
func (_ UserMetadata) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ UserMetadata) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ UserMetadata) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ UserMetadata) SetString(v string)   { panic("Unsupported operation") }
func (_ UserMetadata) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *UserMetadata) Get(i int) types.Field {
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

func (r *UserMetadata) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *UserMetadata) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ UserMetadata) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ UserMetadata) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ UserMetadata) HintSize(int)                     { panic("Unsupported operation") }
func (_ UserMetadata) Finalize()                        {}

func (_ UserMetadata) AvroCRC64Fingerprint() []byte {
	return []byte(UserMetadataAvroCRC64Fingerprint)
}

func (r UserMetadata) MarshalJSON() ([]byte, error) {
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

func (r *UserMetadata) UnmarshalJSON(data []byte) error {
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
