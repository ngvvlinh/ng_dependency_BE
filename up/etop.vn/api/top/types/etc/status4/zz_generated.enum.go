// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package status4

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumStatusName = map[int]string{
	0:  "Z",
	1:  "P",
	2:  "S",
	-1: "N",
}

var enumStatusValue = map[string]int{
	"Z": 0,
	"P": 1,
	"S": 2,
	"N": -1,
}

func ParseStatus(s string) (Status, bool) {
	val, ok := enumStatusValue[s]
	return Status(val), ok
}

func ParseStatusWithDefault(s string, d Status) Status {
	val, ok := enumStatusValue[s]
	if !ok {
		return d
	}
	return Status(val)
}

func ParseStatusWithNull(s dot.NullString, d Status) NullStatus {
	if !s.Valid {
		return NullStatus{}
	}
	val, ok := enumStatusValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return Status(val).Wrap()
}

func (e Status) Enum() int {
	return int(e)
}

func (e Status) Wrap() NullStatus {
	return WrapStatus(e)
}

func (e Status) Name() string {
	return enumStatusName[e.Enum()]
}

func (e Status) String() string {
	s, ok := enumStatusName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("Status(%v)", e.Enum())
}

func (e Status) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumStatusName[e.Enum()] + "\""), nil
}

func (e *Status) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumStatusValue, data, "Status")
	if err != nil {
		return err
	}
	*e = Status(value)
	return nil
}

func (e Status) Value() (driver.Value, error) {
	return int64(e), nil
}

func (e *Status) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumStatusValue, src, "Status")
	*e = (Status)(value)
	return err
}

type NullStatus struct {
	Enum  Status
	Valid bool
}

func WrapStatus(enum Status) NullStatus {
	return NullStatus{Enum: enum, Valid: true}
}

func (n NullStatus) Apply(s Status) Status {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullStatus) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullStatus) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullStatus) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullStatus) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
