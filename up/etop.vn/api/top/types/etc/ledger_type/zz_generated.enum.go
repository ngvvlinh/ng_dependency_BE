// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package ledger_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumLedgerTypeName = map[int]string{
	0: "unknown",
	1: "cash",
	2: "bank",
}

var enumLedgerTypeValue = map[string]int{
	"unknown": 0,
	"cash":    1,
	"bank":    2,
}

func ParseLedgerType(s string) (LedgerType, bool) {
	val, ok := enumLedgerTypeValue[s]
	return LedgerType(val), ok
}

func ParseLedgerTypeWithDefault(s string, d LedgerType) LedgerType {
	val, ok := enumLedgerTypeValue[s]
	if !ok {
		return d
	}
	return LedgerType(val)
}

func ParseLedgerTypeWithNull(s dot.NullString, d LedgerType) NullLedgerType {
	if !s.Valid {
		return NullLedgerType{}
	}
	val, ok := enumLedgerTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return LedgerType(val).Wrap()
}

func (e LedgerType) Enum() int {
	return int(e)
}

func (e LedgerType) Wrap() NullLedgerType {
	return WrapLedgerType(e)
}

func (e LedgerType) Name() string {
	return enumLedgerTypeName[e.Enum()]
}

func (e LedgerType) String() string {
	s, ok := enumLedgerTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("LedgerType(%v)", e.Enum())
}

func (e LedgerType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumLedgerTypeName[e.Enum()] + "\""), nil
}

func (e *LedgerType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumLedgerTypeValue, data, "LedgerType")
	if err != nil {
		return err
	}
	*e = LedgerType(value)
	return nil
}

func (e LedgerType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *LedgerType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumLedgerTypeValue, src, "LedgerType")
	*e = (LedgerType)(value)
	return err
}

type NullLedgerType struct {
	Enum  LedgerType
	Valid bool
}

func WrapLedgerType(enum LedgerType) NullLedgerType {
	return NullLedgerType{Enum: enum, Valid: true}
}

func (n NullLedgerType) Apply(s LedgerType) LedgerType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullLedgerType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullLedgerType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullLedgerType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullLedgerType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
