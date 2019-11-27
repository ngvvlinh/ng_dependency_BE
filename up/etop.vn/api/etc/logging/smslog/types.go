package smslog

import (
	"time"

	"etop.vn/api/pb/etop/etc/status3"
	"etop.vn/capi/dot"
)

type SmsLog struct {
	ID         dot.ID
	ExternalID string
	Phone      string
	Provider   string
	Content    string
	CreatedAt  time.Time
	Status     status3.Status
	Error      string
}
