package receipting

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	. "etop.vn/capi/dot"
	dot "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateReceipt(ctx context.Context, _ *CreateReceiptArgs) (*Receipt, error)
	UpdateReceipt(ctx context.Context, _ *UpdateReceiptArgs) (*Receipt, error)
	CancelReceipt(ctx context.Context, _ *CancelReceiptArgs) (updated int, _ error)
	ConfirmReceipt(ctx context.Context, _ *ConfirmReceiptArgs) (updated int, _ error)
	DeleteReceipt(ctx context.Context, ID, shopID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	GetReceiptByID(context.Context, *GetReceiptByIDArg) (*Receipt, error)
	GetReceiptByCode(ctx context.Context, code string, shopID dot.ID) (*Receipt, error)
	ListReceipts(context.Context, *ListReceiptsArgs) (*ReceiptsResponse, error)
	ListReceiptsByIDs(context.Context, *GetReceiptbyIDsArgs) (*ReceiptsResponse, error)
	ListReceiptsByRefsAndStatus(context.Context, *ListReceiptsByRefsAndStatusArgs) (*ReceiptsResponse, error)
	ListReceiptsByTraderIDsAndStatuses(ctx context.Context, shopID dot.ID, traderIDs []dot.ID, statuses []status3.Status) (*ReceiptsResponse, error)
	ListReceiptsByLedgerIDs(context.Context, *ListReceiptsByLedgerIDsArgs) (*ReceiptsResponse, error)
}

//-- queries --//
type GetReceiptByIDArg struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetReceiptbyIDsArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
}

type ListReceiptsArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ReceiptsResponse struct {
	Receipts                    []*Receipt
	Count                       int
	TotalAmountConfirmedReceipt int
	TotalAmountConfirmedPayment int
	Paging                      meta.PageInfo
}

//-- commands --//

// +convert:create=Receipt
type CreateReceiptArgs struct {
	ShopID      dot.ID
	TraderID    dot.ID
	Title       string
	Type        ReceiptType
	Status      int
	Description string
	Amount      int
	LedgerID    dot.ID
	RefIDs      []dot.ID
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	CreatedBy   dot.ID
	CreatedType ReceiptCreatedType
	ConfirmedAt time.Time
}

// +convert:update=Receipt(ID,ShopID)
type UpdateReceiptArgs struct {
	ID          dot.ID
	ShopID      dot.ID
	TraderID    NullID
	Title       NullString
	Description NullString
	Amount      NullInt
	LedgerID    NullID
	RefIDs      []dot.ID
	RefType     ReceiptRefType
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
}

type CancelReceiptArgs struct {
	ID     dot.ID
	ShopID dot.ID
	Reason string
}

type ConfirmReceiptArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type ListReceiptsByLedgerIDsArgs struct {
	ShopID    dot.ID
	LedgerIDs []dot.ID
	Paging    meta.Paging
	Filters   meta.Filters
}

type ListReceiptsByRefsAndStatusArgs struct {
	ShopID  dot.ID
	RefIDs  []dot.ID
	RefType ReceiptRefType
	Status  int
}
