package types

import (
	"math"
	"strconv"

	"etop.vn/common/l"
)

var (
	ll = l.New()
)

// Int handles special case where we expect integer number but receive string or floating point number.
type Int int

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

func (v Int) Int() int {
	return int(v)
}
