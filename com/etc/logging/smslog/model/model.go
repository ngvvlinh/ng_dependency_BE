package model

import (
	"time"

	"etop.vn/api/pb/etop/etc/status3"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenSMS(&SmsLog{})

// +convert:type=smslog.SmsLog
type SmsLog struct {
	ID         dot.ID
	ExternalID string
	Phone      string
	Provider   string
	Content    string
	CreatedAt  time.Time `sq:"create"`
	Status     status3.Status
	Error      string
}
