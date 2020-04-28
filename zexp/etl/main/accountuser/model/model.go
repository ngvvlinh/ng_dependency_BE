package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type AccountUser struct {
	AccountID dot.ID
	UserID    dot.ID

	Status         status3.Status `sql_type:"int2"` // 1: activated, -1: rejected/disabled, 0: pending
	ResponseStatus status3.Status `sql_type:"int2"` // 1: accepted,  -1: rejected, 0: pending

	CreatedAt time.Time
	UpdatedAt time.Time

	Permission `sq:"inline"`

	FullName  string
	ShortName string
	Position  string

	InvitationSentAt     time.Time
	InvitationSentBy     dot.ID
	InvitationAcceptedAt time.Time
	InvitationRejectedAt time.Time

	DisabledAt    time.Time
	DisabledBy    int8
	DisableReason string

	Rid dot.ID
}

type Permission struct {
	Roles       []string
	Permissions []string
}
