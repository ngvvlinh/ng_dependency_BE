package model

import (
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenInvitation(&Invitation{})

type Invitation struct {
	ID         dot.ID
	AccountID  dot.ID
	Email      string
	Roles      []string
	Token      string
	Status     etop.Status3
	InvitedBy  dot.ID
	AcceptedAt time.Time
	RejectedAt time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
	DeletedAt  time.Time
}
