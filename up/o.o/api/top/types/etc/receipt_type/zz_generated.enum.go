// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package receipt_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumReceiptTypeName = map[int]string{
	0: "unknown",
	1: "receipt",
	2: "payment",
}

var enumReceiptTypeValue = map[string]int{
	"unknown": 0,
	"receipt": 1,
	"payment": 2,
}

var enumReceiptTypeMapLabel = map[string]map[string]string{
	"unknown": {
		"RefName": "",
	},
	"receipt": {
		"RefName": "thu",
	},
	"payment": {
		"RefName": "chi",
	},
}

func (e ReceiptType) GetLabelRefName() string {
	val := enumReceiptTypeName[int(e)]
	nameVal := enumReceiptTypeMapLabel[val]
	return nameVal["RefName"]
}
func ParseReceiptType(s string) (ReceiptType, bool) {
	val, ok := enumReceiptTypeValue[s]
	return ReceiptType(val), ok
}

func ParseReceiptTypeWithDefault(s string, d ReceiptType) ReceiptType {
	val, ok := enumReceiptTypeValue[s]
	if !ok {
		return d
	}
	return ReceiptType(val)
}

func (e ReceiptType) Apply(d ReceiptType) ReceiptType {
	if e == 0 {
		return d
	}
	return e
}

func (e ReceiptType) Enum() int {
	return int(e)
}

func (e ReceiptType) Name() string {
	return enumReceiptTypeName[e.Enum()]
}

func (e ReceiptType) String() string {
	s, ok := enumReceiptTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ReceiptType(%v)", e.Enum())
}

func (e ReceiptType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumReceiptTypeName[e.Enum()] + "\""), nil
}

func (e *ReceiptType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumReceiptTypeValue, data, "ReceiptType")
	if err != nil {
		return err
	}
	*e = ReceiptType(value)
	return nil
}

func (e ReceiptType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ReceiptType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumReceiptTypeValue, src, "ReceiptType")
	*e = (ReceiptType)(value)
	return err
}

func (e ReceiptType) Wrap() NullReceiptType {
	return WrapReceiptType(e)
}

func ParseReceiptTypeWithNull(s dot.NullString, d ReceiptType) NullReceiptType {
	if !s.Valid {
		return NullReceiptType{}
	}
	val, ok := enumReceiptTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ReceiptType(val).Wrap()
}

func WrapReceiptType(enum ReceiptType) NullReceiptType {
	return NullReceiptType{Enum: enum, Valid: true}
}

func (n NullReceiptType) Apply(s ReceiptType) ReceiptType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullReceiptType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullReceiptType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullReceiptType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullReceiptType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
