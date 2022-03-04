package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type BankStatement struct {
	ID                    dot.ID
	Amount                int
	Description           string // format: {mã shop} {số điện thoại}
	AccountID             dot.ID
	TransferedAt          time.Time
	ExternalTransactionID string
	SenderName            string
	SenderBankAccount     string
	OtherInfo             map[string]string
	CreatedAt             time.Time `sq:"create" json:"-"`
	UpdatedAt             time.Time `sq:"update" json:"-"`
}
