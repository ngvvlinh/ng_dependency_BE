// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package shipping_provider

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "etop.vn/capi/dot"
	mix "etop.vn/capi/mix"
)

var __jsonNull = []byte("null")

var enumShippingProviderName = map[int]string{
	0:  "unknown",
	22: "all",
	20: "manual",
	19: "ghn",
	21: "ghtk",
	23: "vtpost",
	24: "etop",
	25: "partner",
}

var enumShippingProviderValue = map[string]int{
	"unknown": 0,
	"all":     22,
	"manual":  20,
	"ghn":     19,
	"ghtk":    21,
	"vtpost":  23,
	"etop":    24,
	"partner": 25,
}

func ParseShippingProvider(s string) (ShippingProvider, bool) {
	val, ok := enumShippingProviderValue[s]
	return ShippingProvider(val), ok
}

func ParseShippingProviderWithDefault(s string, d ShippingProvider) ShippingProvider {
	val, ok := enumShippingProviderValue[s]
	if !ok {
		return d
	}
	return ShippingProvider(val)
}

func (e ShippingProvider) Apply(d ShippingProvider) ShippingProvider {
	if e == 0 {
		return d
	}
	return e
}

func (e ShippingProvider) Enum() int {
	return int(e)
}

func (e ShippingProvider) Name() string {
	return enumShippingProviderName[e.Enum()]
}

func (e ShippingProvider) String() string {
	s, ok := enumShippingProviderName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("ShippingProvider(%v)", e.Enum())
}

func (e ShippingProvider) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumShippingProviderName[e.Enum()] + "\""), nil
}

func (e *ShippingProvider) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumShippingProviderValue, data, "ShippingProvider")
	if err != nil {
		return err
	}
	*e = ShippingProvider(value)
	return nil
}

func (e ShippingProvider) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *ShippingProvider) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumShippingProviderValue, src, "ShippingProvider")
	*e = (ShippingProvider)(value)
	return err
}

func (e ShippingProvider) Wrap() NullShippingProvider {
	return WrapShippingProvider(e)
}

func ParseShippingProviderWithNull(s dot.NullString, d ShippingProvider) NullShippingProvider {
	if !s.Valid {
		return NullShippingProvider{}
	}
	val, ok := enumShippingProviderValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return ShippingProvider(val).Wrap()
}

func WrapShippingProvider(enum ShippingProvider) NullShippingProvider {
	return NullShippingProvider{Enum: enum, Valid: true}
}

func (n NullShippingProvider) Apply(s ShippingProvider) ShippingProvider {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullShippingProvider) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullShippingProvider) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullShippingProvider) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullShippingProvider) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
