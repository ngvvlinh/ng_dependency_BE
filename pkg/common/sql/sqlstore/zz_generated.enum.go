// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package sqlstore

import (
	driver "database/sql/driver"
	fmt "fmt"

	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumPagingFieldName = map[int]string{
	0: "unknown",
	1: "id",
	2: "updated_at",
	3: "external_created_time",
	4: "last_message_at",
	5: "created_at",
}

var enumPagingFieldValue = map[string]int{
	"unknown":               0,
	"id":                    1,
	"updated_at":            2,
	"external_created_time": 3,
	"last_message_at":       4,
	"created_at":            5,
}

func ParsePagingField(s string) (PagingField, bool) {
	val, ok := enumPagingFieldValue[s]
	return PagingField(val), ok
}

func ParsePagingFieldWithDefault(s string, d PagingField) PagingField {
	val, ok := enumPagingFieldValue[s]
	if !ok {
		return d
	}
	return PagingField(val)
}

func (e PagingField) Apply(d PagingField) PagingField {
	if e == 0 {
		return d
	}
	return e
}

func (e PagingField) Enum() int {
	return int(e)
}

func (e PagingField) Name() string {
	return enumPagingFieldName[e.Enum()]
}

func (e PagingField) String() string {
	s, ok := enumPagingFieldName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("PagingField(%v)", e.Enum())
}

func (e PagingField) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumPagingFieldName[e.Enum()] + "\""), nil
}

func (e *PagingField) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumPagingFieldValue, data, "PagingField")
	if err != nil {
		return err
	}
	*e = PagingField(value)
	return nil
}

func (e PagingField) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *PagingField) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumPagingFieldValue, src, "PagingField")
	*e = (PagingField)(value)
	return err
}
