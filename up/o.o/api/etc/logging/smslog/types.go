package smslog

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
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
