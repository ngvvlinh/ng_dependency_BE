// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package price_modifier_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumPriceModifierTypeName = map[int]string{
	1: "percentage",
	2: "fixed_amount",
}

var enumPriceModifierTypeValue = map[string]int{
	"percentage":   1,
	"fixed_amount": 2,
}

func ParsePriceModifierType(s string) (PriceModifierType, bool) {
	val, ok := enumPriceModifierTypeValue[s]
	return PriceModifierType(val), ok
}

func ParsePriceModifierTypeWithDefault(s string, d PriceModifierType) PriceModifierType {
	val, ok := enumPriceModifierTypeValue[s]
	if !ok {
		return d
	}
	return PriceModifierType(val)
}

func (e PriceModifierType) Apply(d PriceModifierType) PriceModifierType {
	if e == 0 {
		return d
	}
	return e
}

func (e PriceModifierType) Enum() int {
	return int(e)
}

func (e PriceModifierType) Name() string {
	return enumPriceModifierTypeName[e.Enum()]
}

func (e PriceModifierType) String() string {
	s, ok := enumPriceModifierTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("PriceModifierType(%v)", e.Enum())
}

func (e PriceModifierType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumPriceModifierTypeName[e.Enum()] + "\""), nil
}

func (e *PriceModifierType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumPriceModifierTypeValue, data, "PriceModifierType")
	if err != nil {
		return err
	}
	*e = PriceModifierType(value)
	return nil
}

func (e PriceModifierType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *PriceModifierType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumPriceModifierTypeValue, src, "PriceModifierType")
	*e = (PriceModifierType)(value)
	return err
}

func (e PriceModifierType) Wrap() NullPriceModifierType {
	return WrapPriceModifierType(e)
}

func ParsePriceModifierTypeWithNull(s dot.NullString, d PriceModifierType) NullPriceModifierType {
	if !s.Valid {
		return NullPriceModifierType{}
	}
	val, ok := enumPriceModifierTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return PriceModifierType(val).Wrap()
}

func WrapPriceModifierType(enum PriceModifierType) NullPriceModifierType {
	return NullPriceModifierType{Enum: enum, Valid: true}
}

func (n NullPriceModifierType) Apply(s PriceModifierType) PriceModifierType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullPriceModifierType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullPriceModifierType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullPriceModifierType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullPriceModifierType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
