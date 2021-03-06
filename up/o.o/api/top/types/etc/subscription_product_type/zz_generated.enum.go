// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package subscription_product_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumProductSubscriptionTypeName = map[int]string{
	0: "unknown",
	1: "ecomify",
	4: "telecom-extension",
}

var enumProductSubscriptionTypeValue = map[string]int{
	"unknown":           0,
	"ecomify":           1,
	"telecom-extension": 4,
}

func ParseProductSubscriptionType(s string) (ProductSubscriptionType, bool) {
	val, ok := enumProductSubscriptionTypeValue[s]
	return ProductSubscriptionType(val), ok
}

func ParseProductSubscriptionTypeWithDefault(s string, d ProductSubscriptionType) ProductSubscriptionType {
	val, ok := enumProductSubscriptionTypeValue[s]
	if !ok {
		return d
	}
	return ProductSubscriptionType(val)
}

func (e ProductSubscriptionType) Apply(d ProductSubscriptionType) ProductSubscriptionType {
	if e == 0 {
		return d
	}
	return e
}

func (e ProductSubscriptionType) Enum() int {
	return int(e)
}

func (e ProductSubscriptionType) Name() string {
	return enumProductSubscriptionTypeName[e.Enum()]
}

func (e ProductSubscriptionType) String() string {
	s, ok := enumProductSubscriptionTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ProductSubscriptionType(%v)", e.Enum())
}

func (e ProductSubscriptionType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumProductSubscriptionTypeName[e.Enum()] + "\""), nil
}

func (e *ProductSubscriptionType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumProductSubscriptionTypeValue, data, "ProductSubscriptionType")
	if err != nil {
		return err
	}
	*e = ProductSubscriptionType(value)
	return nil
}

func (e ProductSubscriptionType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ProductSubscriptionType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumProductSubscriptionTypeValue, src, "ProductSubscriptionType")
	*e = (ProductSubscriptionType)(value)
	return err
}

func (e ProductSubscriptionType) Wrap() NullProductSubscriptionType {
	return WrapProductSubscriptionType(e)
}

func ParseProductSubscriptionTypeWithNull(s dot.NullString, d ProductSubscriptionType) NullProductSubscriptionType {
	if !s.Valid {
		return NullProductSubscriptionType{}
	}
	val, ok := enumProductSubscriptionTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ProductSubscriptionType(val).Wrap()
}

func WrapProductSubscriptionType(enum ProductSubscriptionType) NullProductSubscriptionType {
	return NullProductSubscriptionType{Enum: enum, Valid: true}
}

func (n NullProductSubscriptionType) Apply(s ProductSubscriptionType) ProductSubscriptionType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullProductSubscriptionType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullProductSubscriptionType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullProductSubscriptionType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullProductSubscriptionType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
