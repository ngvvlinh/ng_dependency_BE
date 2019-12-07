package model

import (
	"time"

	"etop.vn/api/top/types/etc/account_type"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

type CreateCreditCommand struct {
	Amount int
	ShopID dot.ID
	Type   account_type.AccountType
	PaidAt time.Time

	Result *CreditExtended
}

type GetCreditQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *CreditExtended
}

type GetCreditsQuery struct {
	ShopID dot.ID
	Paging *cm.Paging

	Result struct {
		Credits []*CreditExtended
		Total   int
	}
}

type ConfirmCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result struct {
		Updated int
	}
}

type UpdateCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID
	PaidAt time.Time
	Amount int

	Result *CreditExtended
}

type DeleteCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result struct {
		Deleted int
	}
}
