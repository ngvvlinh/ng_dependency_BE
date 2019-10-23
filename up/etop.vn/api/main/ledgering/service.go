package ledgering

import (
	"context"

	"etop.vn/api/main/identity"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateLedger(ctx context.Context, _ *CreateLedgerArgs) (*ShopLedger, error)
	UpdateLedger(ctx context.Context, _ *UpdateLedgerArgs) (*ShopLedger, error)
	DeleteLedger(ctx context.Context, ID, ShopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetLedgerByID(context.Context, *shopping.IDQueryShopArg) (*ShopLedger, error)
	ListLedgers(context.Context, *shopping.ListQueryShopArgs) (*ShopLedgersResponse, error)
}

//-- queries --//

type ShopLedgersResponse struct {
	Ledger []*ShopLedger
	Count  int32
	Paging meta.PageInfo
}

//-- commands --//

// +convert:create=ShopLedger
type CreateLedgerArgs struct {
	ShopID      int64
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        string
	CreatedBy   int64
}

// +convert:update=ShopLedger(ID,ShopID)
type UpdateLedgerArgs struct {
	ID          int64
	ShopID      int64
	Name        NullString
	BankAccount *identity.BankAccount
	Note        NullString
}
