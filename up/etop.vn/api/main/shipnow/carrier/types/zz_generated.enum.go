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

var enumCarrierName = map[int]string{
	0: "default",
	1: "ahamove",
}

var enumCarrierValue = map[string]int{
	"default": 0,
	"ahamove": 1,
}

func ParseCarrier(s string) (Carrier, bool) {
	val, ok := enumCarrierValue[s]
	return Carrier(val), ok
}

func ParseCarrierWithDefault(s string, d Carrier) Carrier {
	val, ok := enumCarrierValue[s]
	if !ok {
		return d
	}
	return Carrier(val)
}

func ParseCarrierWithNull(s dot.NullString, d Carrier) NullCarrier {
	if !s.Valid {
		return NullCarrier{}
	}
	val, ok := enumCarrierValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return Carrier(val).Wrap()
}

func (e Carrier) Enum() int {
	return int(e)
}

func (e Carrier) Wrap() NullCarrier {
	return WrapCarrier(e)
}

func (e Carrier) Name() string {
	return enumCarrierName[e.Enum()]
}

func (e Carrier) String() string {
	s, ok := enumCarrierName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("Carrier(%v)", e.Enum())
}

func (e Carrier) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumCarrierName[e.Enum()] + "\""), nil
}

func (e *Carrier) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumCarrierValue, data, "Carrier")
	if err != nil {
		return err
	}
	*e = Carrier(value)
	return nil
}

func (e Carrier) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *Carrier) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumCarrierValue, src, "Carrier")
	*e = (Carrier)(value)
	return err
}

type NullCarrier struct {
	Enum  Carrier
	Valid bool
}

func WrapCarrier(enum Carrier) NullCarrier {
	return NullCarrier{Enum: enum, Valid: true}
}

func (n NullCarrier) Apply(s Carrier) Carrier {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullCarrier) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullCarrier) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullCarrier) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullCarrier) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
