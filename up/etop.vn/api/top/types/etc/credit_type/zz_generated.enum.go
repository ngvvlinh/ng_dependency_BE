// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package credit_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumCreditTypeName = map[int]string{
	1: "shop",
}

var enumCreditTypeValue = map[string]int{
	"shop": 1,
}

func ParseCreditType(s string) (CreditType, bool) {
	val, ok := enumCreditTypeValue[s]
	return CreditType(val), ok
}

func ParseCreditTypeWithDefault(s string, d CreditType) CreditType {
	val, ok := enumCreditTypeValue[s]
	if !ok {
		return d
	}
	return CreditType(val)
}

func ParseCreditTypeWithNull(s dot.NullString, d CreditType) NullCreditType {
	if !s.Valid {
		return NullCreditType{}
	}
	val, ok := enumCreditTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return CreditType(val).Wrap()
}

func (e CreditType) Enum() int {
	return int(e)
}

func (e CreditType) Wrap() NullCreditType {
	return WrapCreditType(e)
}

func (e CreditType) Name() string {
	return enumCreditTypeName[e.Enum()]
}

func (e CreditType) String() string {
	s, ok := enumCreditTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("CreditType(%v)", e.Enum())
}

func (e CreditType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumCreditTypeName[e.Enum()] + "\""), nil
}

func (e *CreditType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumCreditTypeValue, data, "CreditType")
	if err != nil {
		return err
	}
	*e = CreditType(value)
	return nil
}

func (e CreditType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *CreditType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumCreditTypeValue, src, "CreditType")
	*e = (CreditType)(value)
	return err
}

type NullCreditType struct {
	Enum  CreditType
	Valid bool
}

func WrapCreditType(enum CreditType) NullCreditType {
	return NullCreditType{Enum: enum, Valid: true}
}

func (n NullCreditType) Apply(s CreditType) CreditType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullCreditType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullCreditType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullCreditType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullCreditType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
