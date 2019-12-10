// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package payment_provider

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumPaymentProviderName = map[int]string{
	0: "unknown",
	1: "vtpay",
}

var enumPaymentProviderValue = map[string]int{
	"unknown": 0,
	"vtpay":   1,
}

func ParsePaymentProvider(s string) (PaymentProvider, bool) {
	val, ok := enumPaymentProviderValue[s]
	return PaymentProvider(val), ok
}

func ParsePaymentProviderWithDefault(s string, d PaymentProvider) PaymentProvider {
	val, ok := enumPaymentProviderValue[s]
	if !ok {
		return d
	}
	return PaymentProvider(val)
}

func ParsePaymentProviderWithNull(s dot.NullString, d PaymentProvider) NullPaymentProvider {
	if !s.Valid {
		return NullPaymentProvider{}
	}
	val, ok := enumPaymentProviderValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return PaymentProvider(val).Wrap()
}

func (e PaymentProvider) Enum() int {
	return int(e)
}

func (e PaymentProvider) Wrap() NullPaymentProvider {
	return WrapPaymentProvider(e)
}

func (e PaymentProvider) Name() string {
	return enumPaymentProviderName[e.Enum()]
}

func (e PaymentProvider) String() string {
	s, ok := enumPaymentProviderName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("PaymentProvider(%v)", e.Enum())
}

func (e PaymentProvider) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumPaymentProviderName[e.Enum()] + "\""), nil
}

func (e *PaymentProvider) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumPaymentProviderValue, data, "PaymentProvider")
	if err != nil {
		return err
	}
	*e = PaymentProvider(value)
	return nil
}

func (e PaymentProvider) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *PaymentProvider) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumPaymentProviderValue, src, "PaymentProvider")
	*e = (PaymentProvider)(value)
	return err
}

type NullPaymentProvider struct {
	Enum  PaymentProvider
	Valid bool
}

func WrapPaymentProvider(enum PaymentProvider) NullPaymentProvider {
	return NullPaymentProvider{Enum: enum, Valid: true}
}

func (n NullPaymentProvider) Apply(s PaymentProvider) PaymentProvider {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullPaymentProvider) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullPaymentProvider) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullPaymentProvider) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullPaymentProvider) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
