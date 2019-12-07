// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package gender

import (
	driver "database/sql/driver"
	fmt "fmt"

	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumGenderName = map[int]string{
	0: "unknown",
	1: "male",
	2: "female",
	3: "other",
}

var enumGenderValue = map[string]int{
	"unknown": 0,
	"male":    1,
	"female":  2,
	"other":   3,
}

func ParseGender(s string) (Gender, bool) {
	val, ok := enumGenderValue[s]
	return Gender(val), ok
}

func (e Gender) Enum() int {
	return int(e)
}

func (e Gender) Wrap() NullGender {
	return WrapGender(e)
}

func (e Gender) Name() string {
	return enumGenderName[e.Enum()]
}

func (e Gender) String() string {
	s, ok := enumGenderName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("Gender(%v)", e.Enum())
}

func (e Gender) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumGenderName[e.Enum()] + "\""), nil
}

func (e *Gender) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumGenderValue, data, "Gender")
	if err != nil {
		return err
	}
	*e = Gender(value)
	return nil
}

func (e Gender) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *Gender) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumGenderValue, src, "Gender")
	*e = (Gender)(value)
	return err
}

type NullGender struct {
	Enum  Gender
	Valid bool
}

func WrapGender(enum Gender) NullGender {
	return NullGender{Enum: enum, Valid: true}
}

func (n NullGender) Apply(s Gender) Gender {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullGender) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullGender) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullGender) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullGender) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
