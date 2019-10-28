package receipting

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateReceipt(ctx context.Context, _ *CreateReceiptArgs) (*Receipt, error)
	UpdateReceipt(ctx context.Context, _ *UpdateReceiptArgs) (*Receipt, error)
	CancelReceipt(ctx context.Context, _ *CancelReceiptArgs) (updated int, _ error)
	ConfirmReceipt(ctx context.Context, _ *ConfirmReceiptArgs) (updated int, _ error)
	DeleteReceipt(ctx context.Context, ID, shopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetReceiptByID(context.Context, *shopping.IDQueryShopArg) (*Receipt, error)
	GetReceiptByCode(ctx context.Context, code string, shopID int64) (*Receipt, error)
	ListReceipts(context.Context, *shopping.ListQueryShopArgs) (*ReceiptsResponse, error)
	ListReceiptsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*ReceiptsResponse, error)
	ListReceiptsByRefIDs(context.Context, *shopping.IDsQueryShopArgs) (*ReceiptsResponse, error)
	ListReceiptsByTraderIDs(ctx context.Context, shopID int64, traderIDs []int64) (*ReceiptsResponse, error)
	ListReceiptsByLedgerID(ctx context.Context, shopID, ledgerID int64) (*ReceiptsResponse, error)
}

//-- queries --//

type ReceiptsResponse struct {
	Receipts []*Receipt
	Count    int32
	Paging   meta.PageInfo
}

//-- commands --//

// +convert:create=Receipt
type CreateReceiptArgs struct {
	ShopID      int64
	TraderID    int64
	Title       string
	Type        ReceiptType
	Status      int32
	Description string
	Amount      int32
	LedgerID    int64
	RefIDs      []int64
	Lines       []*ReceiptLine
	PaidAt      time.Time
	CreatedBy   int64
	CreatedType string
}

// +convert:update=Receipt(ID,ShopID)
type UpdateReceiptArgs struct {
	ID          int64
	ShopID      int64
	TraderID    NullInt64
	Title       NullString
	Description NullString
	Amount      NullInt32
	LedgerID    NullInt64
	RefIDs      []int64
	Lines       []*ReceiptLine
	PaidAt      time.Time
	CreatedType NullString
}

type CancelReceiptArgs struct {
	ID     int64
	ShopID int64
	Reason string
}

type ConfirmReceiptArgs struct {
	ID     int64
	ShopID int64
}
