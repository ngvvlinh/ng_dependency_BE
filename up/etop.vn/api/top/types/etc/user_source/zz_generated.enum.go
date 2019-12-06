// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package user_source

import (
	fmt "fmt"

	encode "etop.vn/capi/encode"
)

var __jsonNull = []byte("null")

var enumUserSourceName = map[int]string{
	0: "unknown",
	1: "psx",
	2: "etop",
	3: "topship",
	4: "ts_app_android",
	5: "ts_app_ios",
	6: "ts_app_web",
	7: "partner",
}

var enumUserSourceValue = map[string]int{
	"unknown":        0,
	"psx":            1,
	"etop":           2,
	"topship":        3,
	"ts_app_android": 4,
	"ts_app_ios":     5,
	"ts_app_web":     6,
	"partner":        7,
}

func ParseUserSource(s string) (UserSource, bool) {
	val, ok := enumUserSourceValue[s]
	return UserSource(val), ok
}

func (e UserSource) Enum() int {
	return int(e)
}

func (e UserSource) Wrap() NullUserSource {
	return WrapUserSource(e)
}

func (e UserSource) Name() string {
	return enumUserSourceName[e.Enum()]
}

func (e UserSource) String() string {
	s, ok := enumUserSourceName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("UserSource(%v)", e.Enum())
}

func (e UserSource) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumUserSourceName[e.Enum()] + "\""), nil
}

func (e *UserSource) UnmarshalJSON(data []byte) error {
	value, err := encode.UnmarshalJSONEnumInt(enumUserSourceValue, data, "UserSource")
	if err != nil {
		return err
	}
	*e = UserSource(value)
	return nil
}

func (e UserSource) Value() (interface{}, error) {
	return e.String(), nil
}

func (e *UserSource) Scan(src interface{}) error {
	value, err := encode.ScanEnumInt(enumUserSourceValue, src, "UserSource")
	*e = (UserSource)(value)
	return err
}

type NullUserSource struct {
	Enum  UserSource
	Valid bool
}

func WrapUserSource(enum UserSource) NullUserSource {
	return NullUserSource{Enum: enum, Valid: true}
}

func (n NullUserSource) Apply(s UserSource) UserSource {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullUserSource) Value() (interface{}, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullUserSource) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullUserSource) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullUserSource) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
