package transaction

import (
	"time"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/capi/dot"
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
	ID        dot.ID
	Amount    int
	AccountID dot.ID
	Status    etoptypes.Status3
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
