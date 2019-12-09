// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package types

import (
	driver "database/sql/driver"
	fmt "fmt"

	"etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumFulfillName = map[int]string{
	0:  "none",
	1:  "manual",
	10: "shipment",
	11: "shipnow",
}

var enumFulfillValue = map[string]int{
	"none":     0,
	"manual":   1,
	"shipment": 10,
	"shipnow":  11,
}

func ParseFulfill(s string) (Fulfill, bool) {
	val, ok := enumFulfillValue[s]
	return Fulfill(val), ok
}

func ParseFulfillWithDefault(s string, d Fulfill) Fulfill {
	val, ok := enumFulfillValue[s]
	if !ok {
		return d
	}
	return Fulfill(val)
}

func ParseFulfillWithNull(s dot.NullString, d Fulfill) NullFulfill {
	if !s.Valid {
		return NullFulfill{}
	}
	val, ok := enumFulfillValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return Fulfill(val).Wrap()
}

func (e Fulfill) Enum() int {
	return int(e)
}

func (e Fulfill) Wrap() NullFulfill {
	return WrapFulfill(e)
}

func (e Fulfill) Name() string {
	return enumFulfillName[e.Enum()]
}

func (e Fulfill) String() string {
	s, ok := enumFulfillName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("Fulfill(%v)", e.Enum())
}

func (e Fulfill) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumFulfillName[e.Enum()] + "\""), nil
}

func (e *Fulfill) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumFulfillValue, data, "Fulfill")
	if err != nil {
		return err
	}
	*e = Fulfill(value)
	return nil
}

func (e Fulfill) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *Fulfill) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumFulfillValue, src, "Fulfill")
	*e = (Fulfill)(value)
	return err
}

type NullFulfill struct {
	Enum  Fulfill
	Valid bool
}

func WrapFulfill(enum Fulfill) NullFulfill {
	return NullFulfill{Enum: enum, Valid: true}
}

func (n NullFulfill) Apply(s Fulfill) Fulfill {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullFulfill) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullFulfill) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullFulfill) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullFulfill) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}