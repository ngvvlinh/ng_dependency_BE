package call_state

import "o.o/api/top/types/etc/status5"

// +enum
// +enum:zero=null
type CallState int

type NullCallState struct {
	Enum  CallState
	Valid bool
}

const (
	// +enum=unknown
	Unknown CallState = 0

	// +enum=answered
	Answered CallState = 1

	// +enum=not_answered
	NotAnswered CallState = 2
)

func (s CallState) ToStatus5() status5.Status {
	switch s {
	case Answered:
		return status5.P
	case NotAnswered:
		return status5.NS
	default:
		return status5.NS
	}
}
