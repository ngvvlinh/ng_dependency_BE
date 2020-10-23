package receipting

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/top/types/etc/receipt_mode"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	. "o.o/capi/dot"
	dot "o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateReceipt(ctx context.Context, _ *CreateReceiptArgs) (*Receipt, error)
	UpdateReceipt(ctx context.Context, _ *UpdateReceiptArgs) (*Receipt, error)
	CancelReceipt(ctx context.Context, _ *CancelReceiptArgs) (updated int, _ error)
	ConfirmReceipt(ctx context.Context, _ *ConfirmReceiptArgs) (updated int, _ error)
	DeleteReceipt(ctx context.Context, ID, shopID dot.ID) (deleted int, _ error)
	CancelReceiptByRefID(ctx context.Context, _ *CancelReceiptByRefIDRequest) error
}

type QueryService interface {
	GetReceiptByID(context.Context, *GetReceiptByIDArg) (*Receipt, error)
	GetReceiptByCode(ctx context.Context, code string, shopID dot.ID) (*Receipt, error)
	ListReceipts(context.Context, *ListReceiptsArgs) (*ReceiptsResponse, error)
	ListReceiptsByIDs(context.Context, *GetReceiptbyIDsArgs) (*ReceiptsResponse, error)
	ListReceiptsByRefsAndStatus(context.Context, *ListReceiptsByRefsAndStatusArgs) (*ReceiptsResponse, error)
	ListReceiptsByRefsAndStatusAndType(context.Context, *ListReceiptsByRefsAndStatusAndTypeArgs) ([]*Receipt, error)
	ListReceiptsByTraderIDsAndStatuses(ctx context.Context, shopID dot.ID, traderIDs []dot.ID, statuses []status3.Status) (*ReceiptsResponse, error)
	ListReceiptsByLedgerIDs(context.Context, *ListReceiptsByLedgerIDsArgs) (*ReceiptsResponse, error)
}

type CancelReceiptByRefIDRequest struct {
	UpdatedBy dot.ID
	ShopID    dot.ID
	RefID     dot.ID
	RefType   receipt_ref.ReceiptRef
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
	Type        receipt_type.ReceiptType
	Status      int
	Description string
	Amount      int
	LedgerID    dot.ID
	RefIDs      []dot.ID
	RefType     receipt_ref.ReceiptRef
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	CreatedBy   dot.ID
	Mode        receipt_mode.ReceiptMode
	ConfirmedAt time.Time
	Note        string
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
	RefType     receipt_ref.NullReceiptRef
	Lines       []*ReceiptLine
	Trader      *Trader
	PaidAt      time.Time
	Note        NullString
}

type CancelReceiptArgs struct {
	ID           dot.ID
	ShopID       dot.ID
	CancelReason string
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
	ShopID     dot.ID
	RefIDs     []dot.ID
	RefType    receipt_ref.ReceiptRef
	Status     int
	IsContains bool
}

type ListReceiptsByRefsAndStatusAndTypeArgs struct {
	ShopID      dot.ID
	RefIDs      []dot.ID
	RefType     receipt_ref.ReceiptRef
	ReceiptType receipt_type.ReceiptType
	Status      status3.Status
	IsContains  bool
}
