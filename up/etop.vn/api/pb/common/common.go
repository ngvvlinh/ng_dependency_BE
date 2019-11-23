package common

import (
	"errors"

	"github.com/golang/protobuf/jsonpb"
)

func (m *Error) Error() string {
	return m.Msg
}

var _jsonEmptyObject = []byte(`{}`)

// MarshalJSONPB implements JSONPBMarshaler
func (m *RawJSONObject) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	if len(m.Data) == 0 {
		return _jsonEmptyObject, nil
	}
	return m.Data, nil
}

// UnmarshalJSONPB implements JSONPBUnmarshaler
func (m *RawJSONObject) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	if len(data) < 2 || data[0] != '{' || data[len(data)-1] != '}' {
		return errors.New("expect JSON object")
	}
	m.Data = data
	return nil
}
