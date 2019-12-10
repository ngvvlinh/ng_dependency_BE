// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package auth_mode

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumAuthModeName = map[int]string{
	0: "default",
	1: "manual",
}

var enumAuthModeValue = map[string]int{
	"default": 0,
	"manual":  1,
}

func ParseAuthMode(s string) (AuthMode, bool) {
	val, ok := enumAuthModeValue[s]
	return AuthMode(val), ok
}

func ParseAuthModeWithDefault(s string, d AuthMode) AuthMode {
	val, ok := enumAuthModeValue[s]
	if !ok {
		return d
	}
	return AuthMode(val)
}

func ParseAuthModeWithNull(s dot.NullString, d AuthMode) NullAuthMode {
	if !s.Valid {
		return NullAuthMode{}
	}
	val, ok := enumAuthModeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return AuthMode(val).Wrap()
}

func (e AuthMode) Enum() int {
	return int(e)
}

func (e AuthMode) Wrap() NullAuthMode {
	return WrapAuthMode(e)
}

func (e AuthMode) Name() string {
	return enumAuthModeName[e.Enum()]
}

func (e AuthMode) String() string {
	s, ok := enumAuthModeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("AuthMode(%v)", e.Enum())
}

func (e AuthMode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumAuthModeName[e.Enum()] + "\""), nil
}

func (e *AuthMode) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumAuthModeValue, data, "AuthMode")
	if err != nil {
		return err
	}
	*e = AuthMode(value)
	return nil
}

func (e AuthMode) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *AuthMode) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumAuthModeValue, src, "AuthMode")
	*e = (AuthMode)(value)
	return err
}

type NullAuthMode struct {
	Enum  AuthMode
	Valid bool
}

func WrapAuthMode(enum AuthMode) NullAuthMode {
	return NullAuthMode{Enum: enum, Valid: true}
}

func (n NullAuthMode) Apply(s AuthMode) AuthMode {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullAuthMode) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullAuthMode) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullAuthMode) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullAuthMode) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
