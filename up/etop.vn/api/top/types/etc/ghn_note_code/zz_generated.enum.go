// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package ghn_note_code

import (
	driver "database/sql/driver"
	fmt "fmt"

	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumGHNNoteCodeName = map[int]string{
	0: "unknown",
	1: "CHOTHUHANG",
	2: "CHOXEMHANGKHONGTHU",
	3: "KHONGCHOXEMHANG",
}

var enumGHNNoteCodeValue = map[string]int{
	"unknown":            0,
	"CHOTHUHANG":         1,
	"CHOXEMHANGKHONGTHU": 2,
	"KHONGCHOXEMHANG":    3,
}

func ParseGHNNoteCode(s string) (GHNNoteCode, bool) {
	val, ok := enumGHNNoteCodeValue[s]
	return GHNNoteCode(val), ok
}

func (e GHNNoteCode) Enum() int {
	return int(e)
}

func (e GHNNoteCode) Wrap() NullGHNNoteCode {
	return WrapGHNNoteCode(e)
}

func (e GHNNoteCode) Name() string {
	return enumGHNNoteCodeName[e.Enum()]
}

func (e GHNNoteCode) String() string {
	s, ok := enumGHNNoteCodeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("GHNNoteCode(%v)", e.Enum())
}

func (e GHNNoteCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumGHNNoteCodeName[e.Enum()] + "\""), nil
}

func (e *GHNNoteCode) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumGHNNoteCodeValue, data, "GHNNoteCode")
	if err != nil {
		return err
	}
	*e = GHNNoteCode(value)
	return nil
}

func (e GHNNoteCode) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *GHNNoteCode) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumGHNNoteCodeValue, src, "GHNNoteCode")
	*e = (GHNNoteCode)(value)
	return err
}

type NullGHNNoteCode struct {
	Enum  GHNNoteCode
	Valid bool
}

func WrapGHNNoteCode(enum GHNNoteCode) NullGHNNoteCode {
	return NullGHNNoteCode{Enum: enum, Valid: true}
}

func (n NullGHNNoteCode) Apply(s GHNNoteCode) GHNNoteCode {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullGHNNoteCode) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullGHNNoteCode) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullGHNNoteCode) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullGHNNoteCode) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
