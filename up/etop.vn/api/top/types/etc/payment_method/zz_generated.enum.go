// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package payment_method

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumPaymentMethodName = map[int]string{
	0: "unknown",
	1: "cod",
	2: "bank",
	3: "other",
	4: "vtpay",
}

var enumPaymentMethodValue = map[string]int{
	"unknown": 0,
	"cod":     1,
	"bank":    2,
	"other":   3,
	"vtpay":   4,
}

func ParsePaymentMethod(s string) (PaymentMethod, bool) {
	val, ok := enumPaymentMethodValue[s]
	return PaymentMethod(val), ok
}

func ParsePaymentMethodWithDefault(s string, d PaymentMethod) PaymentMethod {
	val, ok := enumPaymentMethodValue[s]
	if !ok {
		return d
	}
	return PaymentMethod(val)
}

func (e PaymentMethod) Enum() int {
	return int(e)
}

func (e PaymentMethod) Name() string {
	return enumPaymentMethodName[e.Enum()]
}

func (e PaymentMethod) String() string {
	s, ok := enumPaymentMethodName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("PaymentMethod(%v)", e.Enum())
}

func (e PaymentMethod) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumPaymentMethodName[e.Enum()] + "\""), nil
}

func (e *PaymentMethod) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumPaymentMethodValue, data, "PaymentMethod")
	if err != nil {
		return err
	}
	*e = PaymentMethod(value)
	return nil
}

func (e PaymentMethod) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *PaymentMethod) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumPaymentMethodValue, src, "PaymentMethod")
	*e = (PaymentMethod)(value)
	return err
}

func (e PaymentMethod) Wrap() NullPaymentMethod {
	return WrapPaymentMethod(e)
}

func ParsePaymentMethodWithNull(s dot.NullString, d PaymentMethod) NullPaymentMethod {
	if !s.Valid {
		return NullPaymentMethod{}
	}
	val, ok := enumPaymentMethodValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return PaymentMethod(val).Wrap()
}

func WrapPaymentMethod(enum PaymentMethod) NullPaymentMethod {
	return NullPaymentMethod{Enum: enum, Valid: true}
}

func (n NullPaymentMethod) Apply(s PaymentMethod) PaymentMethod {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullPaymentMethod) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullPaymentMethod) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullPaymentMethod) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullPaymentMethod) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
