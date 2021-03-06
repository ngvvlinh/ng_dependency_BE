// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package types

import (
	driver "database/sql/driver"
	fmt "fmt"

	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumShippingTypeName = map[int]string{
	0:  "none",
	1:  "manual",
	10: "shipment",
	11: "shipnow",
}

var enumShippingTypeValue = map[string]int{
	"none":     0,
	"manual":   1,
	"shipment": 10,
	"shipnow":  11,
}

func ParseShippingType(s string) (ShippingType, bool) {
	val, ok := enumShippingTypeValue[s]
	return ShippingType(val), ok
}

func ParseShippingTypeWithDefault(s string, d ShippingType) ShippingType {
	val, ok := enumShippingTypeValue[s]
	if !ok {
		return d
	}
	return ShippingType(val)
}

func (e ShippingType) Apply(d ShippingType) ShippingType {
	if e == 0 {
		return d
	}
	return e
}

func (e ShippingType) Enum() int {
	return int(e)
}

func (e ShippingType) Name() string {
	return enumShippingTypeName[e.Enum()]
}

func (e ShippingType) String() string {
	s, ok := enumShippingTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ShippingType(%v)", e.Enum())
}

func (e ShippingType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumShippingTypeName[e.Enum()] + "\""), nil
}

func (e *ShippingType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumShippingTypeValue, data, "ShippingType")
	if err != nil {
		return err
	}
	*e = ShippingType(value)
	return nil
}

func (e ShippingType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return int64(e), nil
}

func (e *ShippingType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumShippingTypeValue, src, "ShippingType")
	*e = (ShippingType)(value)
	return err
}
