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

type User struct {
	Payload UserPayload `json:"payload"`

	Metadata UserMetadata `json:"metadata"`
}

const UserAvroCRC64Fingerprint = "7\xbeɖ\x86\xebg_"

func NewUser() User {
	r := User{}
	r.Payload = NewUserPayload()

	r.Metadata = NewUserMetadata()

	return r
}

func DeserializeUser(r io.Reader) (User, error) {
	t := NewUser()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeUserFromSchema(r io.Reader, schema string) (User, error) {
	t := NewUser()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeUser(r User, w io.Writer) error {
	var err error
	err = writeUserPayload(r.Payload, w)
	if err != nil {
		return err
	}
	err = writeUserMetadata(r.Metadata, w)
	if err != nil {
		return err
	}
	return err
}

func (r User) Serialize(w io.Writer) error {
	return writeUser(r, w)
}

func (r User) Schema() string {
	return "{\"fields\":[{\"name\":\"payload\",\"type\":{\"fields\":[{\"logicalType\":\"uuid\",\"name\":\"uuid\",\"type\":\"string\"},{\"name\":\"given_name\",\"type\":\"string\"},{\"name\":\"family_name\",\"type\":\"string\"},{\"default\":null,\"name\":\"email\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"name\":\"phone_number\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"name\":\"date_of_birth\",\"type\":[\"null\",\"string\"]},{\"name\":\"gender\",\"type\":\"string\"}],\"name\":\"UserPayload\",\"type\":\"record\"}},{\"name\":\"metadata\",\"type\":{\"fields\":[{\"name\":\"origin_application\",\"type\":\"string\"},{\"name\":\"region\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"tracking_id\",\"type\":\"string\"}],\"name\":\"UserMetadata\",\"type\":\"record\"}}],\"name\":\"com.kidsloop.onboarding.user.User\",\"type\":\"record\"}"
}

func (r User) SchemaName() string {
	return "com.kidsloop.onboarding.user.User"
}

func (_ User) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ User) SetInt(v int32)       { panic("Unsupported operation") }
func (_ User) SetLong(v int64)      { panic("Unsupported operation") }
func (_ User) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ User) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ User) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ User) SetString(v string)   { panic("Unsupported operation") }
func (_ User) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *User) Get(i int) types.Field {
	switch i {
	case 0:
		r.Payload = NewUserPayload()

		w := types.Record{Target: &r.Payload}

		return w

	case 1:
		r.Metadata = NewUserMetadata()

		w := types.Record{Target: &r.Metadata}

		return w

	}
	panic("Unknown field index")
}

func (r *User) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *User) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ User) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ User) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ User) HintSize(int)                     { panic("Unsupported operation") }
func (_ User) Finalize()                        {}

func (_ User) AvroCRC64Fingerprint() []byte {
	return []byte(UserAvroCRC64Fingerprint)
}

func (r User) MarshalJSON() ([]byte, error) {
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

func (r *User) UnmarshalJSON(data []byte) error {
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
