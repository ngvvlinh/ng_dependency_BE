package modelx

import (
	"time"

	"o.o/api/top/types/etc/credit_type"
	creditmodel "o.o/backend/com/main/credit/model"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

//TODO(Vu) Delete all

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

type DeleteCreditCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result struct {
		Deleted int
	}
}
