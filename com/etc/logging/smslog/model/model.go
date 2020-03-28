package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +convert:type=smslog.SmsLog
// +sqlgen
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
