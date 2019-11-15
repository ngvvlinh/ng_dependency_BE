package receipting

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
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
	GetReceiptByID(context.Context, *GetReceiptByIDArg) (*Receipt, error)
	GetReceiptByCode(ctx context.Context, code string, shopID int64) (*Receipt, error)
	ListReceipts(context.Context, *ListReceiptsArgs) (*ReceiptsResponse, error)
	ListReceiptsByIDs(context.Context, *GetReceiptbyIDsArgs) (*ReceiptsResponse, error)
	ListReceiptsByRefsAndStatus(context.Context, *ListReceiptsByRefsAndStatusArgs) (*ReceiptsResponse, error)
	ListReceiptsByTraderIDsAndStatuses(ctx context.Context, shopID int64, traderIDs []int64, statuses []etop.Status3) (*ReceiptsResponse, error)
	ListReceiptsByLedgerIDs(context.Context, *ListReceiptsByLedgerIDsArgs) (*ReceiptsResponse, error)
}

//-- queries --//
type GetReceiptByIDArg struct {
	ID     int64
	ShopID int64
}

type GetReceiptbyIDsArgs struct {
	IDs    []int64
	ShopID int64
}

type ListReceiptsArgs struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters
}

type ReceiptsResponse struct {
	Receipts                    []*Receipt
	Count                       int32
	TotalAmountConfirmedReceipt int64
	TotalAmountConfirmedPayment int64
	Paging                      meta.PageInfo
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
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	CreatedBy   int64
	CreatedType ReceiptCreatedType
	ConfirmedAt time.Time
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
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
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

type ListReceiptsByLedgerIDsArgs struct {
	ShopID    int64
	LedgerIDs []int64
	Paging    meta.Paging
	Filters   meta.Filters
}

type ListReceiptsByRefsAndStatusArgs struct {
	ShopID  int64
	RefIDs  []int64
	RefType ReceiptRefType
	Status  int32
}
