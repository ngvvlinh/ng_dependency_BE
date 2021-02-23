// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package charge_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumChargeTypeName = map[int]string{
	1: "prepaid",
	5: "postpaid",
	9: "free",
}

var enumChargeTypeValue = map[string]int{
	"prepaid":  1,
	"postpaid": 5,
	"free":     9,
}

func ParseChargeType(s string) (ChargeType, bool) {
	val, ok := enumChargeTypeValue[s]
	return ChargeType(val), ok
}

func ParseChargeTypeWithDefault(s string, d ChargeType) ChargeType {
	val, ok := enumChargeTypeValue[s]
	if !ok {
		return d
	}
	return ChargeType(val)
}

func (e ChargeType) Apply(d ChargeType) ChargeType {
	if e == 0 {
		return d
	}
	return e
}

func (e ChargeType) Enum() int {
	return int(e)
}

func (e ChargeType) Name() string {
	return enumChargeTypeName[e.Enum()]
}

func (e ChargeType) String() string {
	s, ok := enumChargeTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ChargeType(%v)", e.Enum())
}

func (e ChargeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumChargeTypeName[e.Enum()] + "\""), nil
}

func (e *ChargeType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumChargeTypeValue, data, "ChargeType")
	if err != nil {
		return err
	}
	*e = ChargeType(value)
	return nil
}

func (e ChargeType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ChargeType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumChargeTypeValue, src, "ChargeType")
	*e = (ChargeType)(value)
	return err
}

func (e ChargeType) Wrap() NullChargeType {
	return WrapChargeType(e)
}

func ParseChargeTypeWithNull(s dot.NullString, d ChargeType) NullChargeType {
	if !s.Valid {
		return NullChargeType{}
	}
	val, ok := enumChargeTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ChargeType(val).Wrap()
}

func WrapChargeType(enum ChargeType) NullChargeType {
	return NullChargeType{Enum: enum, Valid: true}
}

func (n NullChargeType) Apply(s ChargeType) ChargeType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullChargeType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullChargeType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullChargeType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullChargeType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
