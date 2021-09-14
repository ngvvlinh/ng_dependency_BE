package driver

import (
	"time"

	"o.o/api/etelecom/call_state"
)

type GetCallLogsRequest struct {
	StartedAt time.Time
	EndedAt   time.Time
	ScrollID  string // VHT pagination
}

type GetCallLogsResponse struct {
	CallLogs []*CallLog
	ScrollID string // VHT pagination
}

type CallLog struct {
	CallID        string
	CallStatus    string
	Caller        string
	Callee        string
	Direction     string
	StartedAt     time.Time
	EndedAt       time.Time
	Duration      int
	AudioURLs     []string
	CallTargets   []*CallTarget
	CallState     call_state.CallState
	HotlineNumber string
	SessionID     string
}

type CallTarget struct {
	TargetNumber string
	TalkDuration int
	CallState    call_state.CallState
	AnsweredTime time.Time
	EndedTime    time.Time
	AddTime      time.Time
	EndReason    string
	FailCode     int
	RingDuration int
	RingTime     time.Time
	Status       string
	TrunkName    string
}
