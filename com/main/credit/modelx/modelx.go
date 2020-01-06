package modelx

import (
	"time"

	"etop.vn/api/top/types/etc/credit_type"
	creditmodel "etop.vn/backend/com/main/credit/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

type CreateCreditCommand struct {
	Amount int
	ShopID dot.ID
	Type   credit_type.CreditType
	PaidAt time.Time

	Result *creditmodel.CreditExtended
}

type GetCreditQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *creditmodel.CreditExtended
}

type GetCreditsQuery struct {
	ShopID dot.ID
	Paging *cm.Paging

	Result struct {
		Credits []*creditmodel.CreditExtended
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

	Result *creditmodel.CreditExtended
}

type DeleteCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result struct {
		Deleted int
	}
}
