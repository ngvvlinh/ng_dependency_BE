package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenInvitation(&Invitation{})

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
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
	DeletedAt  time.Time
}
