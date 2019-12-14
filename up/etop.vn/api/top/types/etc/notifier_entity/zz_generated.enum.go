// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package notifier_entity

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumNotifierEntityName = map[int]string{
	0: "unknown",
	1: "fulfillment",
	2: "money_transaction_shipping",
}

var enumNotifierEntityValue = map[string]int{
	"unknown":                    0,
	"fulfillment":                1,
	"money_transaction_shipping": 2,
}

func ParseNotifierEntity(s string) (NotifierEntity, bool) {
	val, ok := enumNotifierEntityValue[s]
	return NotifierEntity(val), ok
}

func ParseNotifierEntityWithDefault(s string, d NotifierEntity) NotifierEntity {
	val, ok := enumNotifierEntityValue[s]
	if !ok {
		return d
	}
	return NotifierEntity(val)
}

func (e NotifierEntity) Enum() int {
	return int(e)
}

func (e NotifierEntity) Name() string {
	return enumNotifierEntityName[e.Enum()]
}

func (e NotifierEntity) String() string {
	s, ok := enumNotifierEntityName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("NotifierEntity(%v)", e.Enum())
}

func (e NotifierEntity) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumNotifierEntityName[e.Enum()] + "\""), nil
}

func (e *NotifierEntity) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumNotifierEntityValue, data, "NotifierEntity")
	if err != nil {
		return err
	}
	*e = NotifierEntity(value)
	return nil
}

func (e NotifierEntity) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *NotifierEntity) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumNotifierEntityValue, src, "NotifierEntity")
	*e = (NotifierEntity)(value)
	return err
}

func (e NotifierEntity) Wrap() NullNotifierEntity {
	return WrapNotifierEntity(e)
}

func ParseNotifierEntityWithNull(s dot.NullString, d NotifierEntity) NullNotifierEntity {
	if !s.Valid {
		return NullNotifierEntity{}
	}
	val, ok := enumNotifierEntityValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return NotifierEntity(val).Wrap()
}

func WrapNotifierEntity(enum NotifierEntity) NullNotifierEntity {
	return NullNotifierEntity{Enum: enum, Valid: true}
}

func (n NullNotifierEntity) Apply(s NotifierEntity) NotifierEntity {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullNotifierEntity) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullNotifierEntity) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullNotifierEntity) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullNotifierEntity) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
