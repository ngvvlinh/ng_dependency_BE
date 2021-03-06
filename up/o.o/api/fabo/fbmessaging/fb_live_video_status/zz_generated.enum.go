// +build !generator

// Code generated by generator enum. DO NOT EDIT.

package fb_live_video_status

import (
	driver "database/sql/driver"
	fmt "fmt"

	dot "o.o/capi/dot"
	mix "o.o/capi/mix"
)

var __jsonNull = []byte("null")

var enumFbLiveVideoStatusName = map[int]string{
	0:   "unknown",
	54:  "created",
	63:  "live",
	97:  "live_stopped",
	100: "cancelled",
}

var enumFbLiveVideoStatusValue = map[string]int{
	"unknown":      0,
	"created":      54,
	"live":         63,
	"live_stopped": 97,
	"cancelled":    100,
}

func ParseFbLiveVideoStatus(s string) (FbLiveVideoStatus, bool) {
	val, ok := enumFbLiveVideoStatusValue[s]
	return FbLiveVideoStatus(val), ok
}

func ParseFbLiveVideoStatusWithDefault(s string, d FbLiveVideoStatus) FbLiveVideoStatus {
	val, ok := enumFbLiveVideoStatusValue[s]
	if !ok {
		return d
	}
	return FbLiveVideoStatus(val)
}

func (e FbLiveVideoStatus) Apply(d FbLiveVideoStatus) FbLiveVideoStatus {
	if e == 0 {
		return d
	}
	return e
}

func (e FbLiveVideoStatus) Enum() int {
	return int(e)
}

func (e FbLiveVideoStatus) Name() string {
	return enumFbLiveVideoStatusName[e.Enum()]
}

func (e FbLiveVideoStatus) String() string {
	s, ok := enumFbLiveVideoStatusName[e.Enum()]
	if ok {
		return s
	}
	return fmt.Sprintf("FbLiveVideoStatus(%v)", e.Enum())
}

func (e FbLiveVideoStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + enumFbLiveVideoStatusName[e.Enum()] + "\""), nil
}

func (e *FbLiveVideoStatus) UnmarshalJSON(data []byte) error {
	value, err := mix.UnmarshalJSONEnumInt(enumFbLiveVideoStatusValue, data, "FbLiveVideoStatus")
	if err != nil {
		return err
	}
	*e = FbLiveVideoStatus(value)
	return nil
}

func (e FbLiveVideoStatus) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil
	}
	return e.String(), nil
}

func (e *FbLiveVideoStatus) Scan(src interface{}) error {
	value, err := mix.ScanEnumInt(enumFbLiveVideoStatusValue, src, "FbLiveVideoStatus")
	*e = (FbLiveVideoStatus)(value)
	return err
}

func (e FbLiveVideoStatus) Wrap() NullFbLiveVideoStatus {
	return WrapFbLiveVideoStatus(e)
}

func ParseFbLiveVideoStatusWithNull(s dot.NullString, d FbLiveVideoStatus) NullFbLiveVideoStatus {
	if !s.Valid {
		return NullFbLiveVideoStatus{}
	}
	val, ok := enumFbLiveVideoStatusValue[s.String]
	if !ok {
		return d.Wrap()
	}
	return FbLiveVideoStatus(val).Wrap()
}

func WrapFbLiveVideoStatus(enum FbLiveVideoStatus) NullFbLiveVideoStatus {
	return NullFbLiveVideoStatus{Enum: enum, Valid: true}
}

func (n NullFbLiveVideoStatus) Apply(s FbLiveVideoStatus) FbLiveVideoStatus {
	if n.Valid {
		return n.Enum
	}
	return s
}

func (n NullFbLiveVideoStatus) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Enum.Value()
}

func (n *NullFbLiveVideoStatus) Scan(src interface{}) error {
	if src == nil {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.Scan(src)
}

func (n NullFbLiveVideoStatus) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return n.Enum.MarshalJSON()
	}
	return __jsonNull, nil
}

func (n *NullFbLiveVideoStatus) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Enum, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return n.Enum.UnmarshalJSON(data)
}
