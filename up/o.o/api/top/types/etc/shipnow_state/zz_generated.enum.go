// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package shipnow_state

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumStateName = map[int]string{
	0:   "default",
	1:   "created",
	2:   "assigning",
	3:   "picking",
	4:   "delivering",
	5:   "delivered",
	6:   "returning",
	7:   "returned",
	101: "unknown",
	126: "undeliverable",
	127: "cancelled",
}

var enumStateValue = map[string]int{
	"default":       0,
	"created":       1,
	"assigning":     2,
	"picking":       3,
	"delivering":    4,
	"delivered":     5,
	"returning":     6,
	"returned":      7,
	"unknown":       101,
	"undeliverable": 126,
	"cancelled":     127,
}

func ParseState(s string) (State, bool) {
	val, ok := enumStateValue[s]
	return State(val), ok
}

func ParseStateWithDefault(s string, d State) State {
	val, ok := enumStateValue[s]
	if !ok {
		return d
	}
	return State(val)
}

func (e State) Enum() int {
	return int(e)
}

func (e State) Name() string {
	return enumStateName[e.Enum()]
}

func (e State) String() string {
	s, ok := enumStateName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("State(%v)", e.Enum())
}

func (e State) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumStateName[e.Enum()] + "\""), nil
}

func (e *State) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumStateValue, data, "State")
	if err != nil {
		return err
	}
	*e = State(value)
	return nil
}

func (e State) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *State) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumStateValue, src, "State")
	*e = (State)(value)
	return err
}

func (e State) Wrap() NullState {
	return WrapState(e)
}

func ParseStateWithNull(s dot.NullString, d State) NullState {
	if !s.Valid {
		return NullState{}
	}
	val, ok := enumStateValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return State(val).Wrap()
}

func WrapState(enum State) NullState {
	return NullState{Enum: enum, Valid: true}
}

func (n NullState) Apply(s State) State {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullState) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullState) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullState) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullState) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
