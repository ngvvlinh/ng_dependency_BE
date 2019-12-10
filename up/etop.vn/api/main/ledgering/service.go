package ledgering

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/api/top/types/etc/ledger_type"
	. "etop.vn/capi/dot"
	dot "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateLedger(ctx context.Context, _ *CreateLedgerArgs) (*ShopLedger, error)
	UpdateLedger(ctx context.Context, _ *UpdateLedgerArgs) (*ShopLedger, error)
	DeleteLedger(ctx context.Context, ID, ShopID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	GetLedgerByID(context.Context, *shopping.IDQueryShopArg) (*ShopLedger, error)
	// AccountNumber of BankAccount
	GetLedgerByAccountNumber(ctx context.Context, accountNumber string, shopID dot.ID) (*ShopLedger, error)
	ListLedgers(context.Context, *shopping.ListQueryShopArgs) (*ShopLedgersResponse, error)
	ListLedgersByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) (*ShopLedgersResponse, error)
	ListLedgersByType(ctx context.Context, ledgerType ledger_type.LedgerType, shopID dot.ID) (*ShopLedgersResponse, error)
}

//-- queries --//

type ShopLedgersResponse struct {
	Ledgers []*ShopLedger
	Count   int
	Paging  meta.PageInfo
}

//-- commands --//

// +convert:create=ShopLedger
type CreateLedgerArgs struct {
	ShopID      dot.ID
	Name        string
	BankAccount *identity.BankAccount
	Note        string
	Type        ledger_type.LedgerType
	CreatedBy   dot.ID
}

// +convert:update=ShopLedger(ID,ShopID)
type UpdateLedgerArgs struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        NullString
	BankAccount *identity.BankAccount
	Note        NullString
}
