// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package shipping_fee_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	"etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumShippingFeeTypeName = map[int]string{
	0: "unknown",
	1: "main",
	2: "return",
	3: "adjustment",
	4: "insurance",
	5: "tax",
	6: "other",
	7: "cods",
	8: "address_change",
	9: "discount",
}

var enumShippingFeeTypeValue = map[string]int{
	"unknown":        0,
	"main":           1,
	"return":         2,
	"adjustment":     3,
	"insurance":      4,
	"tax":            5,
	"other":          6,
	"cods":           7,
	"address_change": 8,
	"discount":       9,
}

func ParseShippingFeeType(s string) (ShippingFeeType, bool) {
	val, ok := enumShippingFeeTypeValue[s]
	return ShippingFeeType(val), ok
}

func ParseShippingFeeTypeWithDefault(s string, d ShippingFeeType) ShippingFeeType {
	val, ok := enumShippingFeeTypeValue[s]
	if !ok {
		return d
	}
	return ShippingFeeType(val)
}

func ParseShippingFeeTypeWithNull(s dot.NullString, d ShippingFeeType) NullShippingFeeType {
	if !s.Valid {
		return NullShippingFeeType{}
	}
	val, ok := enumShippingFeeTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ShippingFeeType(val).Wrap()
}

func (e ShippingFeeType) Enum() int {
	return int(e)
}

func (e ShippingFeeType) Wrap() NullShippingFeeType {
	return WrapShippingFeeType(e)
}

func (e ShippingFeeType) Name() string {
	return enumShippingFeeTypeName[e.Enum()]
}

func (e ShippingFeeType) String() string {
	s, ok := enumShippingFeeTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ShippingFeeType(%v)", e.Enum())
}

func (e ShippingFeeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumShippingFeeTypeName[e.Enum()] + "\""), nil
}

func (e *ShippingFeeType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumShippingFeeTypeValue, data, "ShippingFeeType")
	if err != nil {
		return err
	}
	*e = ShippingFeeType(value)
	return nil
}

func (e ShippingFeeType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ShippingFeeType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumShippingFeeTypeValue, src, "ShippingFeeType")
	*e = (ShippingFeeType)(value)
	return err
}

type NullShippingFeeType struct {
	Enum  ShippingFeeType
	Valid bool
}

func WrapShippingFeeType(enum ShippingFeeType) NullShippingFeeType {
	return NullShippingFeeType{Enum: enum, Valid: true}
}

func (n NullShippingFeeType) Apply(s ShippingFeeType) ShippingFeeType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullShippingFeeType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullShippingFeeType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullShippingFeeType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullShippingFeeType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
