// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package fb_feed_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumFbFeedTypeName = map[int]string{
	0:   "unknown",
	233: "web",
	278: "app",
}

var enumFbFeedTypeValue = map[string]int{
	"unknown": 0,
	"web":     233,
	"app":     278,
}

func ParseFbFeedType(s string) (FbFeedType, bool) {
	val, ok := enumFbFeedTypeValue[s]
	return FbFeedType(val), ok
}

func ParseFbFeedTypeWithDefault(s string, d FbFeedType) FbFeedType {
	val, ok := enumFbFeedTypeValue[s]
	if !ok {
		return d
	}
	return FbFeedType(val)
}

func (e FbFeedType) Apply(d FbFeedType) FbFeedType {
	if e == 0 {
		return d
	}
	return e
}

func (e FbFeedType) Enum() int {
	return int(e)
}

func (e FbFeedType) Name() string {
	return enumFbFeedTypeName[e.Enum()]
}

func (e FbFeedType) String() string {
	s, ok := enumFbFeedTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("FbFeedType(%v)", e.Enum())
}

func (e FbFeedType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumFbFeedTypeName[e.Enum()] + "\""), nil
}

func (e *FbFeedType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumFbFeedTypeValue, data, "FbFeedType")
	if err != nil {
		return err
	}
	*e = FbFeedType(value)
	return nil
}

func (e FbFeedType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return int64(e), nil
}

func (e *FbFeedType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumFbFeedTypeValue, src, "FbFeedType")
	*e = (FbFeedType)(value)
	return err
}

func (e FbFeedType) Wrap() NullFbFeedType {
	return WrapFbFeedType(e)
}

func ParseFbFeedTypeWithNull(s dot.NullString, d FbFeedType) NullFbFeedType {
	if !s.Valid {
		return NullFbFeedType{}
	}
	val, ok := enumFbFeedTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return FbFeedType(val).Wrap()
}

func WrapFbFeedType(enum FbFeedType) NullFbFeedType {
	return NullFbFeedType{Enum: enum, Valid: true}
}

func (n NullFbFeedType) Apply(s FbFeedType) FbFeedType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullFbFeedType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullFbFeedType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullFbFeedType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullFbFeedType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
