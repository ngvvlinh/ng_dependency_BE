package common

import (
	"errors"
)

func (m *Error) Error() string {
	return m.Msg
}

var _jsonEmptyObject = []byte(`{}`)

// MarshalJSON implements JSONMarshaler
func (m *RawJSONObject) MarshalJSON() ([]byte, error) {
	if len(m.Data) == 0 {
		return _jsonEmptyObject, nil
	}
	return m.Data, nil
}

// UnmarshalJSON implements JSONUnmarshaler
func (m *RawJSONObject) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '{' || data[len(data)-1] != '}' {
		return errors.New("expect JSON object")
	}
	m.Data = data
	return nil
}
