package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"etop.vn/common/jsonx"
)

type Foo struct {
	Number int            `protobuf:"varint,1,opt,name=number" json:"number"`
	Object *RawJSONObject `protobuf:"bytes,2,opt,name=object,json=object" json:"object"`
}

func TestRawJSONObject(t *testing.T) {
	v := &Foo{
		Number: 10,
		Object: &RawJSONObject{
			Data: []byte(`{"foo":"bar"}`),
		},
	}
	t.Run("Marshal object", func(t *testing.T) {
		data, err := jsonx.MarshalToString(v)
		assert.NoError(t, err)
		assert.Equal(t, `{"number":10,"object":{"foo":"bar"}}`, data)
	})

	t.Run("Marshal nil", func(t *testing.T) {
		v.Object = nil
		data, err := jsonx.MarshalToString(v)
		assert.NoError(t, err)
		assert.Equal(t, `{"number":10,"object":null}`, data)
	})

	t.Run("Marshal empty", func(t *testing.T) {
		v.Object = &RawJSONObject{Data: []byte(`{}`)}
		data, err := jsonx.MarshalToString(v)
		assert.NoError(t, err)
		assert.Equal(t, `{"number":10,"object":{}}`, data)
	})

	t.Run("Unmarshal object", func(t *testing.T) {
		data := `{"number":10,"object":{"foo":123}}`

		v.Object = nil
		err := jsonx.UnmarshalString(data, v)
		assert.NoError(t, err)
		assert.Equal(t, `{"foo":123}`, string(v.Object.Data))
	})

	t.Run("Unmarshal empty", func(t *testing.T) {
		data := `{"number":10,"object":{}}`

		v.Object = nil
		err := jsonx.UnmarshalString(data, v)
		assert.NoError(t, err)
		assert.Equal(t, `{}`, string(v.Object.Data))
	})

	t.Run("Unmarshal null", func(t *testing.T) {
		data := `{"number":10,"object":null}`

		v.Object = nil
		err := jsonx.UnmarshalString(data, v)
		assert.NoError(t, err)
	})

	t.Run("Unmarshal undefined", func(t *testing.T) {
		data := `{"number":10}`

		v.Object = nil
		err := jsonx.UnmarshalString(data, v)
		assert.NoError(t, err)
		assert.Nil(t, v.Object)
	})
}
