package cm

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"etop.vn/common/l"
)

const DateLayout = `2006-01-02`
const DateVNLayout = `02-01-2006`
const DateTimeVNLayout = `15:04 02-01-2006`

// TimeAsMillis ...
type TimeAsMillis time.Time

// MarshalJSON implements JSONMarshaler
func (t TimeAsMillis) MarshalJSON() ([]byte, error) {
	var b = make([]byte, 0, 16)
	b = append(b, '"')
	b = strconv.AppendInt(b, int64(ToTimestamp(time.Time(t))), 10)
	b = append(b, '"')
	return b, nil
}

// UnmarshalJSON implements JSONUnmarshaler
func (t *TimeAsMillis) UnmarshalJSON(b []byte) error {
	// Trim quotes
	if len(b) >= 2 && b[0] == '"' {
		b = b[1 : len(b)-1]
	}

	var v int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*t = TimeAsMillis(FromMillis(v))
	return nil
}

// IsZero checks whether time is zero
func (t *TimeAsMillis) IsZero() bool {
	return IsZeroTime(time.Time(*t))
}

// Timestamp represents time as number of milliseconds from 1970.
type Timestamp int64

const e6 = 1e6

// ToTimestamp converts from Go time to timestamp (in nanosecond).
func ToTimestamp(t time.Time) Timestamp {
	if IsZeroTime(t) {
		return 0
	}
	return Timestamp(t.UnixNano() / e6)
}

// Nanos converts from Go time to timestamp (in nanosecond).
func Nanos(t time.Time) int64 {
	if IsZeroTime(t) {
		return 0
	}
	return t.UnixNano() / e6
}

// Millis converts from Go time to timestamp (in millisecond).
func Millis(t time.Time) int64 {
	if IsZeroTime(t) {
		return 0
	}
	return t.UnixNano() / e6
}

// NanosP converts from Go time (with pointer) to timestamp (in nanosecond).
func NanosP(t *time.Time) int64 {
	if t == nil || IsZeroTime(*t) {
		return 0
	}
	return t.UnixNano() / e6
}

// MillisP converts from Go time (with pointer) to timestamp (in millisecond).
func MillisP(t *time.Time) int64 {
	if t == nil || IsZeroTime(*t) {
		return 0
	}
	return t.UnixNano() / e6
}

// FromNanos converts nanosecond to Go time
func FromNanos(t int64) time.Time {
	return Timestamp(t / 1e3).ToTime()
}

// FromNanosP converts nanosecond to Go time with pointer
func FromNanosP(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	tt := Timestamp(t / 1e3).ToTime()
	return &tt
}

// FromMillis converts millisecond to Go time
func FromMillis(t int64) time.Time {
	return Timestamp(t).ToTime()
}

// FromMillisP converts millisecond to Go time with pointer
func FromMillisP(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	tt := Timestamp(t).ToTime()
	return &tt
}

// ToTime converts from timestamp to Go time.
func (t Timestamp) ToTime() time.Time {
	if t == 0 {
		return time.Time{}
	}
	return time.Unix(int64(t)/1e3, int64(t%1e3)*e6).UTC()
}

// Unix converts Timestamp to seconds
func (t Timestamp) Unix() int64 {
	return int64(t) / 1e3
}

// UnixNano extracts nanoseconds from Timestamp
func (t Timestamp) UnixNano() int64 {
	return int64(t) * e6
}

// After reports whether the time instant t is after u
func (t Timestamp) After(u Timestamp) bool {
	return t > u
}

// Before reports whether the time instant t is before u
func (t Timestamp) Before(u Timestamp) bool {
	return t < u
}

// Add adds duration to timestamp
func (t Timestamp) Add(d time.Duration) Timestamp {
	return t + Timestamp(d/e6)
}

// Sub subs timestamp
func (t Timestamp) Sub(u Timestamp) time.Duration {
	return time.Duration((t - u) * e6)
}

// AddDays add a number of days to timestamp
func (t Timestamp) AddDays(days int) Timestamp {
	return t + Timestamp(days*24*60*60*1e3)
}

func (t Timestamp) String() string {
	return t.ToTime().Format(time.RFC3339)
}

// Millis returns the timestamp as number of milliseconds from 1970
func (t Timestamp) Millis() int64 {
	return int64(t)
}

// IsZero reports whether timestamp is zero
func (t Timestamp) IsZero() bool {
	return t == 0
}

// SQLFormat returns time string for using inside SQL query
func (t Timestamp) SQLFormat() string {
	return t.ToTime().Format("2016-01-02")
}

var zeroTime = time.Unix(0, 0)

// IsZeroTime checks whether given time is zero.
// When transport time using grpc, empty time is marshalled to time.Unix(0, 0).
func IsZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(zeroTime)
}

// Now returns current time.
func Now() Timestamp {
	return ToTimestamp(time.Now())
}

// Since returns the time elapsed since t. It is shorthand for cm.Now().Sub(t).
func Since(t Timestamp) time.Duration {
	return time.Since(t.ToTime())
}

func init() {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		ll.Fatal("Unable to load timezone", l.Error(err))
	}

	time.Local = location
	ll.Info("Set default timezone", l.String("location", location.String()))
}

// Handle format "Sáng 2017-07-01"
func FormatDateTimeEdgeCase(s string) (*time.Time, error) {
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	if !re.MatchString(s) {
		return nil, Error(InvalidArgument, "String has not format: `Sáng 2017-07-31`", nil)
	}
	layout := "2006-01-02"
	str := re.FindString(s)
	datetime, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		return nil, err
	}
	var hours int
	mins := 59
	sec := 00
	if strings.Contains(s, "Sáng") {
		hours = 11
	} else {
		hours = 17
	}
	datetime = datetime.Add(time.Hour*time.Duration(hours) +
		time.Minute*time.Duration(mins) +
		time.Second*time.Duration(sec))
	return &datetime, nil
}

func Clock(hour, min int) time.Duration {
	return time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute
}

func ParseDateFromTo(dateFrom, dateTo string) (from, to time.Time, err error) {
	if dateFrom == "" || dateTo == "" {
		err = Error(InvalidArgument, "Must provide both date_from and date_to", nil)
		return
	}
	from, err = time.ParseInLocation(DateLayout, dateFrom, time.Local)
	if err != nil {
		err = Error(InvalidArgument, "invalid date_from", err)
		return
	}
	to, err = time.ParseInLocation(DateLayout, dateTo, time.Local)
	if err != nil {
		err = Error(InvalidArgument, "invalid date_to", err)
		return
	}
	if !from.Before(to) {
		err = Error(InvalidArgument, "Ngày không hợp lệ", nil)
		return
	}
	return
}

// Format Date VN: 26/04/2019
func FormatDateVN(t time.Time) string {
	return t.In(time.Local).Format(DateVNLayout)
}

// Format Data Time VN: 26/04/2019 16:32
func FormatDateTimeVN(t time.Time) string {
	return t.In(time.Local).Format(DateTimeVNLayout)
}
