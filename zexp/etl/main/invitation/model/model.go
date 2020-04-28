package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
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
	Status     status3.Status `sql_type:"int2"`
	InvitedBy  dot.ID
	AcceptedAt time.Time
	RejectedAt time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Rid        dot.ID
}
