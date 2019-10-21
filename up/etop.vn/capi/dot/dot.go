package dot

import "strconv"

type ID int64

func (id ID) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, 32)
	b = append(b, '"')
	b = strconv.AppendInt(b, int64(id), 10)
	b = append(b, '"')
	return b, nil
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

func Bool(b bool) NullBool       { return NullBool{Bool: b, Valid: true} }
func String(s string) NullString { return NullString{String: s, Valid: true} }
func Int64(i int64) NullInt64    { return NullInt64{Int64: i, Valid: true} }
func Int32(i int32) NullInt32    { return NullInt32{Int32: i, Valid: true} }

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
