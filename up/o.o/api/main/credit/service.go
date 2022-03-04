package credit

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/top/types/etc/credit_type"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCredit(context.Context, *CreateCreditArgs) (*CreditExtended, error)
	ConfirmCredit(context.Context, *ConfirmCreditArgs) (*CreditExtended, error)
	DeleteCredit(context.Context, *DeleteCreditArgs) (int, error)
}

type QueryService interface {
	GetCredit(context.Context, *GetCreditArgs) (*CreditExtended, error)
	ListCredits(context.Context, *ListCreditsArgs) (*ListCreditsResponse, error)
}

// +convert:create=Credit
type CreateCreditArgs struct {
	Amount          int
	ShopID          dot.ID
	Type            credit_type.CreditType
	PaidAt          time.Time
	Classify        credit_type.CreditClassify
	BankStatementID dot.ID
}

type ConfirmCreditArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type DeleteCreditArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetCreditArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type ListCreditsArgs struct {
	ShopID   dot.ID
	Classify credit_type.NullCreditClassify
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
	Paging   *meta.Paging
}

type ListCreditsResponse struct {
	Credits []*CreditExtended
	Paging  *meta.PageInfo
}

type GetTotalCreditArgs struct {
	UserID   dot.ID
	ShopIDs  []dot.ID
	Classify credit_type.CreditClassify
}
