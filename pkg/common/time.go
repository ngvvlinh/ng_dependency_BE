package cm

import (
	"time"

	"etop.vn/common/l"
)

const DateLayout = `2006-01-02`
const DateVNLayout = `02-01-2006`
const DateTimeVNLayout = `15:04 02-01-2006`

func init() {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		ll.Fatal("Unable to load timezone", l.Error(err))
	}

	time.Local = location
	ll.Info("Set default timezone", l.String("location", location.String()))
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
