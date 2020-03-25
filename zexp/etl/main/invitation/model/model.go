package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +sqlgen
type Invitation struct {
	ID         dot.ID
	AccountID  dot.ID
	Email      string
	Phone      string
	FullName   string
	ShortName  string
	Roles      []string
	Token      string
	Status     status3.Status
	InvitedBy  dot.ID
	AcceptedAt time.Time
	RejectedAt time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Rid        dot.ID
}
