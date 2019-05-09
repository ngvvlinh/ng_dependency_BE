package common

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/assert"
)

type Foo struct {
	Number int32          `protobuf:"varint,1,opt,name=number" json:"number,omitempty"`
	Object *RawJSONObject `protobuf:"bytes,2,opt,name=object,json=object" json:"object,omitempty"`
}

func (f *Foo) Reset()         {}
func (f *Foo) String() string { return "" }
func (f *Foo) ProtoMessage()  {}

// Implement proto.Message
var _ proto.Message = &Foo{}

func TestRawJSONObject(t *testing.T) {
	var jSON = runtime.JSONPb{
		OrigName:     true,
		EmitDefaults: true,
	}

	v := &Foo{
		Number: 10,
		Object: &RawJSONObject{
			Data: []byte(`{"foo":"bar"}`),
		},
	}
	t.Run("Marshal object", func(t *testing.T) {
		data, err := jSON.Marshal(v)
		assert.NoError(t, err)
		assert.Equal(t, `{"number":10,"object":{"foo":"bar"}}`, string(data))
	})

	t.Run("Marshal nil", func(t *testing.T) {
		v.Object = nil
		data, err := jSON.Marshal(v)
		assert.NoError(t, err)
		assert.Equal(t, `{"number":10,"object":null}`, string(data))
	})

	t.Run("Marshal empty", func(t *testing.T) {
		v.Object = &RawJSONObject{Data: []byte(`{}`)}
		data, err := jSON.Marshal(v)
		assert.NoError(t, err)
		assert.Equal(t, `{"number":10,"object":{}}`, string(data))
	})

	t.Run("Unmarshal object", func(t *testing.T) {
		data := `{"number":10,"object":{"foo":123}}`

		v.Object = nil
		err := jSON.Unmarshal([]byte(data), &v)
		assert.NoError(t, err)
		assert.Equal(t, `{"foo":123}`, string(v.Object.Data))
	})

	t.Run("Unmarshal empty", func(t *testing.T) {
		data := `{"number":10,"object":{}}`

		v.Object = nil
		err := jSON.Unmarshal([]byte(data), &v)
		assert.NoError(t, err)
		assert.Equal(t, `{}`, string(v.Object.Data))
	})

	t.Run("Unmarshal null", func(t *testing.T) {
		data := `{"number":10,"object":null}`

		v.Object = nil
		err := jSON.Unmarshal([]byte(data), &v)
		assert.EqualError(t, err, "Expect JSON object")
	})

	t.Run("Unmarshal undefined", func(t *testing.T) {
		data := `{"number":10}`

		v.Object = nil
		err := jSON.Unmarshal([]byte(data), &v)
		assert.NoError(t, err)
		assert.Nil(t, v.Object)
	})
}
