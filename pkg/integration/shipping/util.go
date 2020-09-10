package shipping

import (
	"regexp"
	"strings"
	"time"

	cm "o.o/backend/pkg/common"
)

var reDateTimeShipping = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// ParseDateTimeShipping handles strings returned from external api
//
//     "Sáng 2020-02-02"  -> "11:59 2020-02-02"
//     "Chiều 2020-02-02" -> "17:59 2020-02-02"
func ParseDateTimeShipping(s string) (*time.Time, error) {
	if !reDateTimeShipping.MatchString(s) {
		return nil, cm.Error(cm.InvalidArgument, "string must have format: `Sáng 2017-07-31`", nil)
	}
	s = strings.ToLower(s)

	layout := "2006-01-02"
	str := reDateTimeShipping.FindString(s)
	datetime, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		return nil, err
	}
	var hours time.Duration
	if strings.Contains(s, "sáng") {
		hours = 11
	} else {
		hours = 17
	}
	datetime = datetime.Add(hours*time.Hour + 59*time.Minute)
	return &datetime, nil
}

func AppendString(s, appendStr string) string {
	if !strings.HasSuffix(s, ".") {
		s += ". \n"
	} else {
		s += "\n"
	}
	s += appendStr
	return s
}
