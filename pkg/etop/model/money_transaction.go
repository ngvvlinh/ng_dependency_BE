package model

import (
	"time"

	cm "etop.vn/backend/pkg/common"
)

type CreateCreditCommand struct {
	Amount     int
	ShopID     int64
	SupplierID int64
	Type       AccountType
	PaidAt     time.Time

	Result *CreditExtended
}

type GetCreditQuery struct {
	ID         int64
	ShopID     int64
	SupplierID int64

	Result *CreditExtended
}

type GetCreditsQuery struct {
	ShopID int64
	Paging *cm.Paging

	Result struct {
		Credits []*CreditExtended
		Total   int
	}
}

type ConfirmCreditCommand struct {
	ID     int64
	ShopID int64

	Result struct {
		Updated int
	}
}

type UpdateCreditCommand struct {
	ID     int64
	ShopID int64
	PaidAt time.Time
	Amount int

	Result *CreditExtended
}

type DeleteCreditCommand struct {
	ID     int64
	ShopID int64

	Result struct {
		Deleted int
	}
}
