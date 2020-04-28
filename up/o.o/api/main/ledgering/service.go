package ledgering

import (
	"context"

	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/shopping"
	"o.o/api/top/types/etc/ledger_type"
	. "o.o/capi/dot"
	dot "o.o/capi/dot"
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
	Paging  meta.PageInfo
}

//-- commands --//

// +convert:create=ShopLedger
type CreateLedgerArgs struct {
	ShopID      dot.ID
	Name        string
	BankAccount *identitytypes.BankAccount
	Note        string
	Type        ledger_type.LedgerType
	CreatedBy   dot.ID
}

// +convert:update=ShopLedger(ID,ShopID)
type UpdateLedgerArgs struct {
	ID          dot.ID
	ShopID      dot.ID
	Name        NullString
	BankAccount *identitytypes.BankAccount
	Note        NullString
}
