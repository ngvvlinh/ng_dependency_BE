package model

import "time"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenTransaction(&Transaction{})

type Transaction struct {
	ID        int64
	Amount    int
	AccountID int64
	Status    int
	Type      string
	Note      string
	Metadata  *TransactionMetadata
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

type TransactionMetadata struct {
	ReferralType string
	ReferralIDs  []int64
}
