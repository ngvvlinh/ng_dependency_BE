package model

import (
	"time"

	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type Transaction struct {
	ID        dot.ID
	Amount    int
	AccountID dot.ID
	Status    status3.Status
	Type      transaction.TransactionType
	Note      string
	Metadata  *TransactionMetadata
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

type TransactionMetadata struct {
	ReferralType string
	ReferralIDs  []dot.ID
}
