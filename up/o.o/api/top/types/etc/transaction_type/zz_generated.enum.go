// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package transaction_type

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumTransactionTypeName = map[int]string{
	9:  "credit",
	13: "invoice",
}

var enumTransactionTypeValue = map[string]int{
	"credit":  9,
	"invoice": 13,
}

func ParseTransactionType(s string) (TransactionType, bool) {
	val, ok := enumTransactionTypeValue[s]
	return TransactionType(val), ok
}

func ParseTransactionTypeWithDefault(s string, d TransactionType) TransactionType {
	val, ok := enumTransactionTypeValue[s]
	if !ok {
		return d
	}
	return TransactionType(val)
}

func (e TransactionType) Apply(d TransactionType) TransactionType {
	if e == 0 {
		return d
	}
	return e
}

func (e TransactionType) Enum() int {
	return int(e)
}

func (e TransactionType) Name() string {
	return enumTransactionTypeName[e.Enum()]
}

func (e TransactionType) String() string {
	s, ok := enumTransactionTypeName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("TransactionType(%v)", e.Enum())
}

func (e TransactionType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumTransactionTypeName[e.Enum()] + "\""), nil
}

func (e *TransactionType) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumTransactionTypeValue, data, "TransactionType")
	if err != nil {
		return err
	}
	*e = TransactionType(value)
	return nil
}

func (e TransactionType) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *TransactionType) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumTransactionTypeValue, src, "TransactionType")
	*e = (TransactionType)(value)
	return err
}

func (e TransactionType) Wrap() NullTransactionType {
	return WrapTransactionType(e)
}

func ParseTransactionTypeWithNull(s dot.NullString, d TransactionType) NullTransactionType {
	if !s.Valid {
		return NullTransactionType{}
	}
	val, ok := enumTransactionTypeValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return TransactionType(val).Wrap()
}

func WrapTransactionType(enum TransactionType) NullTransactionType {
	return NullTransactionType{Enum: enum, Valid: true}
}

func (n NullTransactionType) Apply(s TransactionType) TransactionType {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullTransactionType) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullTransactionType) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullTransactionType) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullTransactionType) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}