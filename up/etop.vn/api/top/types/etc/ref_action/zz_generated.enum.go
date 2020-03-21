// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package ref_action

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumRefActionName = map[int]string{
	0: "create",
	1: "cancel",
}

var enumRefActionValue = map[string]int{
	"create": 0,
	"cancel": 1,
}

func ParseRefAction(s string) (RefAction, bool) {
	val, ok := enumRefActionValue[s]
	return RefAction(val), ok
}

func ParseRefActionWithDefault(s string, d RefAction) RefAction {
	val, ok := enumRefActionValue[s]
	if !ok {
		return d
	}
	return RefAction(val)
}

func (e RefAction) Apply(d RefAction) RefAction {
	if e == 0 {
		return d
	}
	return e
}

func (e RefAction) Enum() int {
	return int(e)
}

func (e RefAction) Name() string {
	return enumRefActionName[e.Enum()]
}

func (e RefAction) String() string {
	s, ok := enumRefActionName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("RefAction(%v)", e.Enum())
}

func (e RefAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumRefActionName[e.Enum()] + "\""), nil
}

func (e *RefAction) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumRefActionValue, data, "RefAction")
	if err != nil {
		return err
	}
	*e = RefAction(value)
	return nil
}

func (e RefAction) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *RefAction) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumRefActionValue, src, "RefAction")
	*e = (RefAction)(value)
	return err
}

func (e RefAction) Wrap() NullRefAction {
	return WrapRefAction(e)
}

func ParseRefActionWithNull(s dot.NullString, d RefAction) NullRefAction {
	if !s.Valid {
		return NullRefAction{}
	}
	val, ok := enumRefActionValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return RefAction(val).Wrap()
}

func WrapRefAction(enum RefAction) NullRefAction {
	return NullRefAction{Enum: enum, Valid: true}
}

func (n NullRefAction) Apply(s RefAction) RefAction {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullRefAction) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullRefAction) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullRefAction) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullRefAction) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
