package dot

import (
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

type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

type NullID struct {
	ID    ID
	Valid bool
}

func NID(id ID) NullID              { return NullID{ID: id, Valid: true} }
func Bool(b bool) NullBool          { return NullBool{Bool: b, Valid: true} }
func String(s string) NullString    { return NullString{String: s, Valid: true} }
func Int64(i int64) NullInt64       { return NullInt64{Int64: i, Valid: true} }
func Int32(i int32) NullInt32       { return NullInt32{Int32: i, Valid: true} }
func Int(i int) NullInt             { return NullInt{Int: i, Valid: true} }
func Float64(f float64) NullFloat64 { return NullFloat64{Float64: f, Valid: true} }

func PBool(b *bool) NullBool {
	if b == nil {
		return NullBool{}
	}
	return Bool(*b)
}

func PString(s NullString) NullString {
	return s
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

func (n NullInt) Apply(i int) int {
	if n.Valid {
		return n.Int
	}
	return i
}

func (ni NullID) Apply(i ID) ID {
	if ni.Valid {
		return ni.ID
	}
	return i
}

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
