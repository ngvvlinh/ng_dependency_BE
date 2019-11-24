package dot

import (
	"database/sql"
	"database/sql/driver"
	"strconv"
	"time"
)

type ID int64

func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id ID) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 32)
	b = append(b, '"')
	b = strconv.AppendInt(b, int64(id), 10)
	b = append(b, '"')
	return b, nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*id = 0
		return nil
	}
	if data[0] == '"' {
		data = data[1 : len(data)-1]
	}
	i, err := strconv.ParseInt(string(data), 10, 64)
	*id = ID(i)
	return err
}

func (id ID) Value() (driver.Value, error) {
	return int64(id), nil
}

func (i *ID) Scan(src interface{}) error {
	var ni sql.NullInt64
	err := ni.Scan(src)
	if err == nil && ni.Valid {
		*i = ID(ni.Int64)
	}
	return err
}

type NullBool struct {
	Bool  bool
	Valid bool
}

type NullString struct {
	String string
	Valid  bool
}

type NullInt64 struct {
	Int64 int64
	Valid bool
}

type NullInt32 struct {
	Int32 int32
	Valid bool
}

type NullInt struct {
	Int   int
	Valid bool
}

type NullID struct {
	ID    ID
	Valid bool
}

func Bool(b bool) NullBool       { return NullBool{Bool: b, Valid: true} }
func String(s string) NullString { return NullString{String: s, Valid: true} }
func Int64(i int64) NullInt64    { return NullInt64{Int64: i, Valid: true} }
func Int32(i int32) NullInt32    { return NullInt32{Int32: i, Valid: true} }
func Int(i int) NullInt          { return NullInt{Int: i, Valid: true} }

func PBool(b *bool) NullBool {
	if b == nil {
		return NullBool{}
	}
	return Bool(*b)
}

func PString(s *string) NullString {
	if s == nil {
		return NullString{}
	}
	return String(*s)
}

func PInt64(i *int64) NullInt64 {
	if i == nil {
		return NullInt64{}
	}
	return Int64(*i)
}

func PInt32(i *int32) NullInt32 {
	if i == nil {
		return NullInt32{}
	}
	return Int32(*i)
}

func PInt(i *int) NullInt {
	if i == nil {
		return NullInt{}
	}
	return Int(*i)
}

func PID(i *ID) NullID {
	if i == nil {
		return NullID{}
	}
	return NullID{ID: *i, Valid: true}
}

func (ns NullString) Apply(s string) string {
	if ns.Valid {
		return ns.String
	}
	return s
}

func (nb NullBool) Apply(b bool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return b
}

func (ni NullInt64) Apply(i int64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return i
}

func (ni NullInt32) Apply(i int32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return i
}

func (ni NullInt) Apply(i int) int {
	if ni.Valid {
		return ni.Int
	}
	return i
}

func (ni NullID) Apply(i ID) ID {
	if ni.Valid {
		return ni.ID
	}
	return i
}

var jsonNull = []byte("null")
var zeroTime = time.Unix(0, 0)

// IsZeroTime checks whether given time is zero.
// When transport time using grpc, empty time is marshalled to time.Unix(0, 0).
func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(zeroTime)
}

type Time time.Time

func (t Time) ToTime() time.Time {
	return time.Time(t)
}

func (t Time) IsZero() bool {
	return time.Time(t).IsZero() || time.Time(t).Equal(zeroTime)
}

func (t Time) MarshalJSON() ([]byte, error) {
	if IsZeroTime(time.Time(t)) {
		return jsonNull, nil
	}
	tt := time.Time(t)
	tt = tt.Add(-time.Duration(tt.Nanosecond()))
	return tt.MarshalJSON()
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var tt time.Time
	err := tt.UnmarshalJSON(data)
	if err != nil || IsZeroTime(tt) {
		*t = Time{}
	} else {
		*t = Time(tt)
	}
	return err
}
