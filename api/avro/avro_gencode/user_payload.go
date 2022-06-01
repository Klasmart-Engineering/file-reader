// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCES:
 *     organization.avsc
 *     school.avsc
 *     user.avsc
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

type UserPayload struct {
	Uuid string `json:"uuid"`

	Given_name string `json:"given_name"`

	Family_name string `json:"family_name"`

	Email string `json:"email"`

	Phone_number *UnionNullString `json:"phone_number"`

	Date_of_birth string `json:"date_of_birth"`

	Gender string `json:"gender"`
}

const UserPayloadAvroCRC64Fingerprint = "\xa5L\x89EgO\xa48"

func NewUserPayload() UserPayload {
	r := UserPayload{}
	return r
}

func DeserializeUserPayload(r io.Reader) (UserPayload, error) {
	t := NewUserPayload()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeUserPayloadFromSchema(r io.Reader, schema string) (UserPayload, error) {
	t := NewUserPayload()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeUserPayload(r UserPayload, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Uuid, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Given_name, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Family_name, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Email, w)
	if err != nil {
		return err
	}
	err = writeUnionNullString(r.Phone_number, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Date_of_birth, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Gender, w)
	if err != nil {
		return err
	}
	return err
}

func (r UserPayload) Serialize(w io.Writer) error {
	return writeUserPayload(r, w)
}

func (r UserPayload) Schema() string {
	return "{\"fields\":[{\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":\"string\"},{\"name\":\"given_name\",\"type\":\"string\"},{\"name\":\"family_name\",\"type\":\"string\"},{\"name\":\"email\",\"type\":\"string\"},{\"name\":\"phone_number\",\"type\":[\"null\",\"string\"]},{\"name\":\"date_of_birth\",\"type\":\"string\"},{\"name\":\"gender\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.user.UserPayload\",\"type\":\"record\"}"
}

func (r UserPayload) SchemaName() string {
	return "com.kidsloop.onboarding.user.UserPayload"
}

func (_ UserPayload) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ UserPayload) SetInt(v int32)       { panic("Unsupported operation") }
func (_ UserPayload) SetLong(v int64)      { panic("Unsupported operation") }
func (_ UserPayload) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ UserPayload) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ UserPayload) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ UserPayload) SetString(v string)   { panic("Unsupported operation") }
func (_ UserPayload) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *UserPayload) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Uuid}

		return w

	case 1:
		w := types.String{Target: &r.Given_name}

		return w

	case 2:
		w := types.String{Target: &r.Family_name}

		return w

	case 3:
		w := types.String{Target: &r.Email}

		return w

	case 4:
		r.Phone_number = NewUnionNullString()

		return r.Phone_number
	case 5:
		w := types.String{Target: &r.Date_of_birth}

		return w

	case 6:
		w := types.String{Target: &r.Gender}

		return w

	}
	panic("Unknown field index")
}

func (r *UserPayload) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *UserPayload) NullField(i int) {
	switch i {
	case 4:
		r.Phone_number = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ UserPayload) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ UserPayload) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ UserPayload) HintSize(int)                     { panic("Unsupported operation") }
func (_ UserPayload) Finalize()                        {}

func (_ UserPayload) AvroCRC64Fingerprint() []byte {
	return []byte(UserPayloadAvroCRC64Fingerprint)
}

func (r UserPayload) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["uuid"], err = json.Marshal(r.Uuid)
	if err != nil {
		return nil, err
	}
	output["given_name"], err = json.Marshal(r.Given_name)
	if err != nil {
		return nil, err
	}
	output["family_name"], err = json.Marshal(r.Family_name)
	if err != nil {
		return nil, err
	}
	output["email"], err = json.Marshal(r.Email)
	if err != nil {
		return nil, err
	}
	output["phone_number"], err = json.Marshal(r.Phone_number)
	if err != nil {
		return nil, err
	}
	output["date_of_birth"], err = json.Marshal(r.Date_of_birth)
	if err != nil {
		return nil, err
	}
	output["gender"], err = json.Marshal(r.Gender)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *UserPayload) UnmarshalJSON(data []byte) error {
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
		return fmt.Errorf("no value specified for uuid")
	}
	val = func() json.RawMessage {
		if v, ok := fields["given_name"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Given_name); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for given_name")
	}
	val = func() json.RawMessage {
		if v, ok := fields["family_name"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Family_name); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for family_name")
	}
	val = func() json.RawMessage {
		if v, ok := fields["email"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Email); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for email")
	}
	val = func() json.RawMessage {
		if v, ok := fields["phone_number"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Phone_number); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for phone_number")
	}
	val = func() json.RawMessage {
		if v, ok := fields["date_of_birth"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Date_of_birth); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for date_of_birth")
	}
	val = func() json.RawMessage {
		if v, ok := fields["gender"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Gender); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for gender")
	}
	return nil
}