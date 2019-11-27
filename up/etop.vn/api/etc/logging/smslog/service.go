package smslog

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/pb/etop/etc/status3"
)

// +gen:api

type Aggregate interface {
	CreateSmsLog(ctx context.Context, _ *CreateSmsArgs) error
}

// -- querys -- //

type GetSmsLogs struct {
	Filters meta.Filters
}

//-- commands --//

// +convert:create=SmsLog
type CreateSmsArgs struct {
	ExternalID string
	Phone      string
	Provider   string
	Content    string
	Status     status3.Status
	Error      string
}
