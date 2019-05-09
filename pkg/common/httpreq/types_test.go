package httpreq

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var t0, t1 time.Time

func init() {
	local, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		panic(err)
	}

	layout := "2006-01-02T15:04:05.999"
	t0, err = time.ParseInLocation(layout, "2010-10-20T01:02:03.000", local)
	if err != nil {
		panic(err)
	}

	// All times are presented in local timezone
	t1 = t0.Add(450 * time.Millisecond)
}

func TestJSONInt(t *testing.T) {
	tests := []struct {
		s  string
		i  Int
		ok bool
	}{
		{`null`, 0, true},
		{`""`, 0, true},
		{`10`, 10, true},
		{`10.0`, 10, true},
		{`0.0`, 0, true},
		{`10.5`, 10, true}, // Ignore float number
		{`"10"`, 10, true},
		{`"10.0"`, 10, true},
		{`"0.0"`, 0, true},
		{`"10.5"`, 10, true}, // Ignore float number
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			var output Int
			err := json.Unmarshal([]byte(tt.s), &output)
			if tt.ok {
				assert.NoError(t, err)
				assert.Equal(t, tt.i, output)
			} else {
				assert.Error(t, err, "Expected integer but got float number")
			}
		})
	}
}

func TestJSONFloat(t *testing.T) {
	tests := []struct {
		s  string
		f  Float
		ok bool
	}{
		{`null`, 0, true},
		{`10.5`, 10.5, true},
		{`"10.5"`, 10.5, true},
		{`"oops"`, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			var output Float
			err := json.Unmarshal([]byte(tt.s), &output)
			if tt.ok {
				assert.NoError(t, err)
				assert.Equal(t, tt.f, output)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestJSONTime(t *testing.T) {
	tests := []struct {
		s string
		t time.Time
	}{
		{`null`, time.Time{}},
		{`""`, time.Time{}},
		{`"/Date(1287511323000)/"`, t0},
		{`"/Date(1287511323450)/"`, t1},
		{`"2010-10-20T01:02:03"`, t0},
		{`"2010-10-20T01:02:03+07:00"`, t0},
		{`"2010-10-20T01:02:03.45"`, t1},
		{`"2010-10-20T01:02:03.4500000"`, t1},
		{`"2010-10-20T01:02:03.45+07:00"`, t1},
		{`"2010-10-20T01:02:03.4500000+07:00"`, t1},
		{`"2010-10-20T01:02:03Z"`, t0.In(time.UTC).Add(7 * time.Hour)},
		{`"2010-10-20T01:02:03.45Z"`, t1.In(time.UTC).Add(7 * time.Hour)},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			var output Time
			err := json.Unmarshal([]byte(tt.s), &output)
			assert.NoError(t, err)
			assert.Equal(t, Time(tt.t), output)
		})
	}

	t.Run(`"invalid"`, func(t *testing.T) {
		var output Time
		err := json.Unmarshal([]byte(`"invalid"`), &output)
		assert.EqualError(t, err, `Unable to parse time "invalid"`)
	})
}

func Test_parseAsMilliseconds(t *testing.T) {
	tests := []struct {
		s  string
		t  time.Time
		ok bool
	}{
		{`null`, time.Time{}, false},
		{`"invalid"`, time.Time{}, false},
		{`"/Date(1287511323000)/"`, t0, true},
		{`"/Date(1287511323450)/"`, t1, true},
		{`"1287511323000"`, t0, true},
		{`1287511323000`, t0, true},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			output, ok := parseAsMiliseconds([]byte(tt.s))
			assert.Equal(t, tt.t, output)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func Test_parseAsISO8601(t *testing.T) {
	tests := []struct {
		s  string
		t  time.Time
		ok bool
	}{
		{``, time.Time{}, false},
		{`invalid`, time.Time{}, false},
		{`2010-10-20T01:02:03`, t0, true},
		{`2010-10-20T01:02:03.000`, t0, true},
		{`2010-10-20T01:02:03.0000000`, t0, true},
		{`2010-10-20T01:02:03+07:00`, t0, true},
		{`2010-10-20T01:02:03.45`, t1, true},
		{`2010-10-20T01:02:03.45+07:00`, t1, true},
		{`2010-10-20T01:02:03.450000+07:00`, t1, true},
		{`2010-10-20T01:02:03Z`, t0.In(time.UTC).Add(7 * time.Hour), true},
		{`2010-10-20T01:02:03.45Z`, t1.In(time.UTC).Add(7 * time.Hour), true},
		{`2010-10-20T01:02:03.450123Z`, t1.In(time.UTC).Add(7 * time.Hour).Add(123 * time.Microsecond), true},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			output, ok := parseAsISO8601([]byte(tt.s))
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.t, output)
		})
	}
}
