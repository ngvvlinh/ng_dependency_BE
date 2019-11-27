package model

import (
	"time"

	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenVthCallHistory(&VhtCallHistory{})

type VhtCallHistory struct {
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
	Direction       int
	TimeStarted     time.Time
	TimeConnected   time.Time
	TimeEnded       time.Time
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	RecordingPath   string
	RecordingURL    string
	RecordFileSize  int
	EtopAccountID   dot.ID
	VtigerAccountID string
	SyncStatus      string
	OData           string
	SearchNorm      string
}

func (p *VhtCallHistory) BeforeInsertOrUpdate() error {
	s := p.ToNumber + " " + p.FromNumber
	p.SearchNorm = validate.NormalizeSearch(s)
	return nil
}
