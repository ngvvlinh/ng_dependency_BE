// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package subject_referral

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumSubjectReferralName = map[int]string{
	3: "credit",
	7: "invoice",
	9: "subscription",
}

var enumSubjectReferralValue = map[string]int{
	"credit":       3,
	"invoice":      7,
	"subscription": 9,
}

func ParseSubjectReferral(s string) (SubjectReferral, bool) {
	val, ok := enumSubjectReferralValue[s]
	return SubjectReferral(val), ok
}

func ParseSubjectReferralWithDefault(s string, d SubjectReferral) SubjectReferral {
	val, ok := enumSubjectReferralValue[s]
	if !ok {
		return d
	}
	return SubjectReferral(val)
}

func (e SubjectReferral) Apply(d SubjectReferral) SubjectReferral {
	if e == 0 {
		return d
	}
	return e
}

func (e SubjectReferral) Enum() int {
	return int(e)
}

func (e SubjectReferral) Name() string {
	return enumSubjectReferralName[e.Enum()]
}

func (e SubjectReferral) String() string {
	s, ok := enumSubjectReferralName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("SubjectReferral(%v)", e.Enum())
}

func (e SubjectReferral) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumSubjectReferralName[e.Enum()] + "\""), nil
}

func (e *SubjectReferral) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumSubjectReferralValue, data, "SubjectReferral")
	if err != nil {
		return err
	}
	*e = SubjectReferral(value)
	return nil
}

func (e SubjectReferral) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *SubjectReferral) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumSubjectReferralValue, src, "SubjectReferral")
	*e = (SubjectReferral)(value)
	return err
}

func (e SubjectReferral) Wrap() NullSubjectReferral {
	return WrapSubjectReferral(e)
}

func ParseSubjectReferralWithNull(s dot.NullString, d SubjectReferral) NullSubjectReferral {
	if !s.Valid {
		return NullSubjectReferral{}
	}
	val, ok := enumSubjectReferralValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return SubjectReferral(val).Wrap()
}

func WrapSubjectReferral(enum SubjectReferral) NullSubjectReferral {
	return NullSubjectReferral{Enum: enum, Valid: true}
}

func (n NullSubjectReferral) Apply(s SubjectReferral) SubjectReferral {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullSubjectReferral) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullSubjectReferral) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullSubjectReferral) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullSubjectReferral) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
