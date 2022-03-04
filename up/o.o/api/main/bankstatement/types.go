package bankstatement

import (
	"time"

	"o.o/capi/dot"
)

// +gen:event:topic=event/bank_statement

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
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type BankStatementCreatedEvent struct {
	ID dot.ID
}
