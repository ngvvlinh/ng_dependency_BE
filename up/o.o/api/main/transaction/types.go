package transaction

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

type (
	TransactionType string
	ReferralType    string
)

func (t TransactionType) String() string { return string(t) }
func (t ReferralType) String() string    { return string(t) }

var (
	TransactionTypeAffiliate TransactionType = "affiliate"

	ReferralTypeOrder ReferralType = "order"
)

type Transaction struct {
	ID        dot.ID
	Amount    int
	AccountID dot.ID
	Status    status3.Status
	Type      TransactionType
	Note      string
	Metadata  *TransactionMetadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionMetadata struct {
	ReferralType ReferralType
	ReferralIDs  []dot.ID
}
