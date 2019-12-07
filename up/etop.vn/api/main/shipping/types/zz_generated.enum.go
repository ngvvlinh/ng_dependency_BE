// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package types

import (
	driver "database/sql/driver"
	fmt "fmt"

	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumFeeLineTypeName = map[int]string{
	0: "other",
	1: "main",
	2: "return",
	3: "adjustment",
	4: "cods",
	5: "insurance",
	6: "address_change",
	7: "discount",
}

var enumFeeLineTypeValue = map[string]int{
	"other":          0,
	"main":           1,
	"return":         2,
	"adjustment":     3,
	"cods":           4,
	"insurance":      5,
	"address_change": 6,
	"discount":       7,
}

func ParseFeeLineType(s string) (FeeLineType, bool) {
	val, ok := enumFeeLineTypeValue[s]
	return FeeLineType(val), ok
}

func (e FeeLineType) Enum() int {
	return int(e)
}

func (e FeeLineType) Wrap() NullFeeLineType {
	return WrapFeeLineType(e)
}

func (e FeeLineType) Name() string {
	return enumFeeLineTypeName[e.Enum()]
}

func (e FeeLineType) String() string {
	s, ok := enumFeeLineTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("FeeLineType(%v)", e.Enum())
}

func (e FeeLineType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumFeeLineTypeName[e.Enum()] + "\""), nil
}

func (e *FeeLineType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumFeeLineTypeValue, data, "FeeLineType")
	if err != nil {
		return err
	}
	*e = FeeLineType(value)
	return nil
}

func (e FeeLineType) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *FeeLineType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumFeeLineTypeValue, src, "FeeLineType")
	*e = (FeeLineType)(value)
	return err
}

type NullFeeLineType struct {
	Enum  FeeLineType
	Valid bool
}

func WrapFeeLineType(enum FeeLineType) NullFeeLineType {
	return NullFeeLineType{Enum: enum, Valid: true}
}

func (n NullFeeLineType) Apply(s FeeLineType) FeeLineType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullFeeLineType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullFeeLineType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullFeeLineType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullFeeLineType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}

var enumTryOnName = map[int]string{
	0: "unknown",
	1: "none",
	2: "open",
	3: "try",
}

var enumTryOnValue = map[string]int{
	"unknown": 0,
	"none":    1,
	"open":    2,
	"try":     3,
}

func ParseTryOn(s string) (TryOn, bool) {
	val, ok := enumTryOnValue[s]
	return TryOn(val), ok
}

func (e TryOn) Enum() int {
	return int(e)
}

func (e TryOn) Wrap() NullTryOn {
	return WrapTryOn(e)
}

func (e TryOn) Name() string {
	return enumTryOnName[e.Enum()]
}

func (e TryOn) String() string {
	s, ok := enumTryOnName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("TryOn(%v)", e.Enum())
}

func (e TryOn) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumTryOnName[e.Enum()] + "\""), nil
}

func (e *TryOn) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumTryOnValue, data, "TryOn")
	if err != nil {
		return err
	}
	*e = TryOn(value)
	return nil
}

func (e TryOn) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *TryOn) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumTryOnValue, src, "TryOn")
	*e = (TryOn)(value)
	return err
}

type NullTryOn struct {
	Enum  TryOn
	Valid bool
}

func WrapTryOn(enum TryOn) NullTryOn {
	return NullTryOn{Enum: enum, Valid: true}
}

func (n NullTryOn) Apply(s TryOn) TryOn {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullTryOn) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullTryOn) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullTryOn) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullTryOn) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
