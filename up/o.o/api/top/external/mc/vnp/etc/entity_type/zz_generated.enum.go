// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package entity_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumEntityTypeName = map[int]string{
	53: "shipnow_fulfillment",
}

var enumEntityTypeValue = map[string]int{
	"shipnow_fulfillment": 53,
}

func ParseEntityType(s string) (EntityType, bool) {
	val, ok := enumEntityTypeValue[s]
	return EntityType(val), ok
}

func ParseEntityTypeWithDefault(s string, d EntityType) EntityType {
	val, ok := enumEntityTypeValue[s]
	if !ok {
		return d
	}
	return EntityType(val)
}

func (e EntityType) Apply(d EntityType) EntityType {
	if e == 0 {
		return d
	}
	return e
}

func (e EntityType) Enum() int {
	return int(e)
}

func (e EntityType) Name() string {
	return enumEntityTypeName[e.Enum()]
}

func (e EntityType) String() string {
	s, ok := enumEntityTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("EntityType(%v)", e.Enum())
}

func (e EntityType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumEntityTypeName[e.Enum()] + "\""), nil
}

func (e *EntityType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumEntityTypeValue, data, "EntityType")
	if err != nil {
		return err
	}
	*e = EntityType(value)
	return nil
}

func (e EntityType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *EntityType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumEntityTypeValue, src, "EntityType")
	*e = (EntityType)(value)
	return err
}

func (e EntityType) Wrap() NullEntityType {
	return WrapEntityType(e)
}

func ParseEntityTypeWithNull(s dot.NullString, d EntityType) NullEntityType {
	if !s.Valid {
		return NullEntityType{}
	}
	val, ok := enumEntityTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return EntityType(val).Wrap()
}

func WrapEntityType(enum EntityType) NullEntityType {
	return NullEntityType{Enum: enum, Valid: true}
}

func (n NullEntityType) Apply(s EntityType) EntityType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullEntityType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullEntityType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullEntityType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullEntityType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
