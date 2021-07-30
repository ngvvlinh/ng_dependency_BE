// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package sms_provider

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumSmsProviderName = map[int]string{
	0: "unknown",
	1: "mock",
	2: "telegram",
	3: "vietguys",
	4: "incom",
}

var enumSmsProviderValue = map[string]int{
	"unknown":  0,
	"mock":     1,
	"telegram": 2,
	"vietguys": 3,
	"incom":    4,
}

func ParseSmsProvider(s string) (SmsProvider, bool) {
	val, ok := enumSmsProviderValue[s]
	return SmsProvider(val), ok
}

func ParseSmsProviderWithDefault(s string, d SmsProvider) SmsProvider {
	val, ok := enumSmsProviderValue[s]
	if !ok {
		return d
	}
	return SmsProvider(val)
}

func (e SmsProvider) Apply(d SmsProvider) SmsProvider {
	if e == 0 {
		return d
	}
	return e
}

func (e SmsProvider) Enum() int {
	return int(e)
}

func (e SmsProvider) Name() string {
	return enumSmsProviderName[e.Enum()]
}

func (e SmsProvider) String() string {
	s, ok := enumSmsProviderName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("SmsProvider(%v)", e.Enum())
}

func (e SmsProvider) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumSmsProviderName[e.Enum()] + "\""), nil
}

func (e *SmsProvider) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumSmsProviderValue, data, "SmsProvider")
	if err != nil {
		return err
	}
	*e = SmsProvider(value)
	return nil
}

func (e SmsProvider) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *SmsProvider) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumSmsProviderValue, src, "SmsProvider")
	*e = (SmsProvider)(value)
	return err
}

func (e SmsProvider) Wrap() NullSmsProvider {
	return WrapSmsProvider(e)
}

func ParseSmsProviderWithNull(s dot.NullString, d SmsProvider) NullSmsProvider {
	if !s.Valid {
		return NullSmsProvider{}
	}
	val, ok := enumSmsProviderValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return SmsProvider(val).Wrap()
}

func WrapSmsProvider(enum SmsProvider) NullSmsProvider {
	return NullSmsProvider{Enum: enum, Valid: true}
}

func (n NullSmsProvider) Apply(s SmsProvider) SmsProvider {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullSmsProvider) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullSmsProvider) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullSmsProvider) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullSmsProvider) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}