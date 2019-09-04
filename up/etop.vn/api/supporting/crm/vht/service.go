package vht

import (
	"context"
	"time"

	"etop.vn/api/meta"
)

type Aggregate interface {
	CreateOrUpdateCallHistoryBySDKCallID(context.Context, *VhtCallLog) (*VhtCallLog, error)
	CreateOrUpdateCallHistoryByCallID(context.Context, *VhtCallLog) (*VhtCallLog, error)
	PingServerVht(context.Context, *meta.Empty) error
	SyncVhtCallHistories(context.Context, *SyncVhtCallHistoriesArgs) error
}

type QueryService interface {
	GetCallHistories(context.Context, *GetCallHistoriesArgs) (*GetCallHistoriesResponse, error)
	GetLastCallHistory(context.Context, meta.Paging) (*VhtCallLog, error)
}

type VhtCallLog struct {
	Direction       int32
	CdrID           string
	CallID          string
	SipCallID       string
	SdkCallID       string
	Cause           string
	Q850Cause       string
	FromExtension   string
	ToExtension     string
	FromNumber      string
	ToNumber        string
	Duration        int32
	TimeStarted     time.Time
	TimeConnected   time.Time
	TimeEnded       time.Time
	RecordingPath   string
	RecordingUrl    string
	RecordFileSize  int32
	EtopAccountID   int64
	VtigerAccountID string
}

type GetCallHistoriesArgs struct {
	Paging     *meta.Paging
	TextSearch string
}

type GetCallHistoriesResponse struct {
	VhtCallLog []*VhtCallLog
}

type SyncVhtCallHistoriesArgs struct {
	SyncTime time.Time
}
