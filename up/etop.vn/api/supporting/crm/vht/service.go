package vht

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:api

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
	Direction       int
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
	Duration        int
	TimeStarted     time.Time
	TimeConnected   time.Time
	TimeEnded       time.Time
	RecordingPath   string
	RecordingUrl    string
	RecordFileSize  int
	EtopAccountID   dot.ID
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
