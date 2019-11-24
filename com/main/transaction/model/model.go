package model

import (
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenTransaction(&Transaction{})

type Transaction struct {
	ID        dot.ID
	Amount    int
	AccountID dot.ID
	Status    int
	Type      string
	Note      string
	Metadata  *TransactionMetadata
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

type TransactionMetadata struct {
	ReferralType string
	ReferralIDs  []dot.ID
}
