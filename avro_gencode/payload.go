// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     organisation.avsc
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

type Payload struct {
	Organization_name string `json:"organization_name"`

	Guid string `json:"guid"`
}

const PayloadAvroCRC64Fingerprint = "ߒ\xe2u6\x9c\xea\xe1"

func NewPayload() Payload {
	r := Payload{}
	return r
}

func DeserializePayload(r io.Reader) (Payload, error) {
	t := NewPayload()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializePayloadFromSchema(r io.Reader, schema string) (Payload, error) {
	t := NewPayload()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writePayload(r Payload, w io.Writer) error {
	var err error
	err = vm.WriteString(r.Organization_name, w)
	if err != nil {
		return err
	}
	err = vm.WriteString(r.Guid, w)
	if err != nil {
		return err
	}
	return err
}

func (r Payload) Serialize(w io.Writer) error {
	return writePayload(r, w)
}

func (r Payload) Schema() string {
	return "{\"fields\":[{\"name\":\"organization_name\",\"type\":\"string\"},{\"logicalType\":\"uuid\",\"name\":\"guid\",\"type\":\"string\"}],\"name\":\"com.kidsloop.onboarding.payload\",\"type\":\"record\"}"
}

func (r Payload) SchemaName() string {
	return "com.kidsloop.onboarding.payload"
}

func (_ Payload) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ Payload) SetInt(v int32)       { panic("Unsupported operation") }
func (_ Payload) SetLong(v int64)      { panic("Unsupported operation") }
func (_ Payload) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ Payload) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ Payload) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ Payload) SetString(v string)   { panic("Unsupported operation") }
func (_ Payload) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *Payload) Get(i int) types.Field {
	switch i {
	case 0:
		w := types.String{Target: &r.Organization_name}

		return w

	case 1:
		w := types.String{Target: &r.Guid}

		return w

	}
	panic("Unknown field index")
}

func (r *Payload) SetDefault(i int) {
	switch i {
	}
	panic("Unknown field index")
}

func (r *Payload) NullField(i int) {
	switch i {
	}
	panic("Not a nullable field index")
}

func (_ Payload) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ Payload) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ Payload) HintSize(int)                     { panic("Unsupported operation") }
func (_ Payload) Finalize()                        {}

func (_ Payload) AvroCRC64Fingerprint() []byte {
	return []byte(PayloadAvroCRC64Fingerprint)
}

func (r Payload) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["organization_name"], err = json.Marshal(r.Organization_name)
	if err != nil {
		return nil, err
	}
	output["guid"], err = json.Marshal(r.Guid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *Payload) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["organization_name"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Organization_name); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for organization_name")
	}
	val = func() json.RawMessage {
		if v, ok := fields["guid"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.Guid); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no value specified for guid")
	}
	return nil
}
