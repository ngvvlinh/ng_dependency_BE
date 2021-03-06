// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package report_time_filter

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumTimeFilterName = map[int]string{
	0: "unknown",
	1: "month",
	2: "quater",
	3: "year",
}

var enumTimeFilterValue = map[string]int{
	"unknown": 0,
	"month":   1,
	"quater":  2,
	"year":    3,
}

func ParseTimeFilter(s string) (TimeFilter, bool) {
	val, ok := enumTimeFilterValue[s]
	return TimeFilter(val), ok
}

func ParseTimeFilterWithDefault(s string, d TimeFilter) TimeFilter {
	val, ok := enumTimeFilterValue[s]
	if !ok {
		return d
	}
	return TimeFilter(val)
}

func (e TimeFilter) Apply(d TimeFilter) TimeFilter {
	if e == 0 {
		return d
	}
	return e
}

func (e TimeFilter) Enum() int {
	return int(e)
}

func (e TimeFilter) Name() string {
	return enumTimeFilterName[e.Enum()]
}

func (e TimeFilter) String() string {
	s, ok := enumTimeFilterName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("TimeFilter(%v)", e.Enum())
}

func (e TimeFilter) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumTimeFilterName[e.Enum()] + "\""), nil
}

func (e *TimeFilter) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumTimeFilterValue, data, "TimeFilter")
	if err != nil {
		return err
	}
	*e = TimeFilter(value)
	return nil
}

func (e TimeFilter) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *TimeFilter) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumTimeFilterValue, src, "TimeFilter")
	*e = (TimeFilter)(value)
	return err
}

func (e TimeFilter) Wrap() NullTimeFilter {
	return WrapTimeFilter(e)
}

func ParseTimeFilterWithNull(s dot.NullString, d TimeFilter) NullTimeFilter {
	if !s.Valid {
		return NullTimeFilter{}
	}
	val, ok := enumTimeFilterValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return TimeFilter(val).Wrap()
}

func WrapTimeFilter(enum TimeFilter) NullTimeFilter {
	return NullTimeFilter{Enum: enum, Valid: true}
}

func (n NullTimeFilter) Apply(s TimeFilter) TimeFilter {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullTimeFilter) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullTimeFilter) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullTimeFilter) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullTimeFilter) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
