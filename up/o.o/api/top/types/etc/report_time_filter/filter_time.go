package report_time_filter

// +enum
// +enum:zero=null
type TimeFilter int

type NullTimeFilter struct {
	Enum  TimeFilter
	Valid bool
}

const (
	// +enum=unknown
	Unknown TimeFilter = 0

	// +enum=month
	Month TimeFilter = 1

	// +enum=quater
	Quater TimeFilter = 2

	// +enum=year
	Year TimeFilter = 3
)

func (t TimeFilter) ShortName() string {
	switch t {
	case Month:
		return "T"
	case Quater:
		return "Q"
	default:
		return ""
	}
}
