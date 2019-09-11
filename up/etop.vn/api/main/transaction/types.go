package transaction

import (
	"time"

	etoptypes "etop.vn/api/main/etop"
)

type (
	TransactionType string
	ReferralType    string
)

var (
	TransactionTypeAffiliate TransactionType = "affiliate"

	ReferralTypeOrder ReferralType = "order"
)

type Transaction struct {
	ID        int64
	Amount    int
	AccountID int64
	Status    etoptypes.Status3
	Type      TransactionType
	Note      string
	Metadata  *TransactionMetadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionMetadata struct {
	ReferralType ReferralType
	ReferralIDs  []int64
}
