package httpreq

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"etop.vn/common/l"
)

// Bool handles null, string and bool from json as bool
type Bool bool

// UnmarshalJSON implements json.Unmarshaler
func (b *Bool) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	s := string(data)
	if s == "null" {
		return nil
	}

	bb, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	*b = Bool(bb)
	return nil
}

// String handles null, int, float, etc. from json as string
type String string

// UnmarshalJSON implements json.Unmarshaler
func (s *String) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		*s = String(data[1 : len(data)-1])
		return nil
	}
	ss := String(data)
	if ss != "null" {
		*s = ss
	}
	return nil
}

func (s String) String() string {
	return string(s)
}

// Int handles special case where we expect integer number but receive string or floating point number.
type Int int

// ErrFloatJSON is returned by unmarshalling a floating point number into Int
var ErrFloatJSON = errors.New("expected integer but got float number")

// UnmarshalJSON parses float number or string as int.
func (v *Int) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Handle string as number
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}

	s := string(data)
	if s == "" || s == "null" {
		return nil
	}

	i, err := strconv.Atoi(s)
	if err == nil {
		*v = Int(i)
		return nil
	}

	// Handle float number as int
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*v = Int(f)

	// The float number must be equal to integer part
	_, frac := math.Modf(f)
	if frac != 0 {
		ll.Warn("expect int but got float", l.Any("f", f))
	}
	return nil
}

func (v Int) String() string {
	return strconv.Itoa(int(v))
}

// Float handles null, string and float from json as float
type Float float64

// UnmarshalJSON implements json.Unmarshaler
func (f *Float) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	s := string(data)
	if s == "null" {
		return nil
	}

	ff, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = Float(ff)
	return nil
}

// Time is used to decode time returned by PDS.
// There are a few variants we have to deal with:
//
// 		/Date(1485157445877)/
//		2016-05-11T09:37:39.23
//		2016-05-11T09:37.39+07:00
type Time time.Time

// MarshalJSON encodes time in JSON
func (t Time) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

// ToTime converts to Go time.Time
func (t Time) ToTime() time.Time {
	tt := time.Time(t)
	if tt.Year() < 1990 {
		return time.Time{}
	}
	return tt
}

func (t Time) IsZero() bool {
	tt := time.Time(t)
	return tt.IsZero() || tt.Year() < 1990
}

// UnmarshalJSON parses JSON time
func (t *Time) UnmarshalJSON(data []byte) error {
	if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
		if tt, err := time.Parse(time.RFC3339, string(data)); err == nil {
			*t = Time(tt)
			return nil
		}
		if tt, ok := parseAsMiliseconds(data); ok {
			*t = Time(tt)
			return nil
		}
		if tt, ok := parseAsISO8601(data[1 : len(data)-1]); ok {
			*t = Time(tt)
			return nil
		}
		ll.Error("Unable to parse time", l.String("data", string(data)))
		return fmt.Errorf(`unable to parse time %s`, data)
	}

	// Zero time
	return nil
}

var reDateString = regexp.MustCompile(`[0-9]{12,14}`)

// Parse time string which has following format:
//
//  	/Date(1485157445877)/
func parseAsMiliseconds(data []byte) (time.Time, bool) {
	match := reDateString.Find(data)
	if match == nil {
		return time.Time{}, false
	}
	s := string(match)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, false
	}

	// This function returns time in local timezone.
	return time.Unix(n/1e3, n%1e3*1e6), true
}

var (
	timeLayout  = `2006-01-02T15:04:05`
	timeISO8601 = `2006-01-02T15:04:05Z07:00`
	reISO8601   = regexp.MustCompile(`^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]{1,9})?(.*))$`)
)

func ParseAsISO8601(data []byte) (time.Time, bool) {
	return parseAsISO8601(data)
}

// Parse time string which has following format (ISO 8601).
// If 'Z' is present, parse as UTC time; otherwise as local time.
//
//    2016-05-11T09:37:39.23
func parseAsISO8601(data []byte) (time.Time, bool) {
	match := reISO8601.FindSubmatch(data)
	if match == nil {
		return time.Time{}, false
	}

	var t time.Time
	var err error

	// Detect time with or without timezone
	if len(match[3]) > 0 {
		t, err = time.Parse(timeISO8601, string(data))
	} else {
		s := string(match[1])
		t, err = time.ParseInLocation(timeLayout, s, time.Local)
	}
	return t, err == nil
}
