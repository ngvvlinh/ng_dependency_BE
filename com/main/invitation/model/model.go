package model

import (
	"time"

	"etop.vn/api/main/etop"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenInvitation(&Invitation{})

type Invitation struct {
	ID         int64
	AccountID  int64
	Email      string
	Roles      []string
	Token      string
	Status     etop.Status3
	InvitedBy  int64
	AcceptedAt time.Time
	RejectedAt time.Time
	ExpiresAt  time.Time
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
	DeletedAt  time.Time
}
