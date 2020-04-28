package purchaserefund

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreatePurchaseRefund(ctx context.Context, _ *CreatePurchaseRefundArgs) (*PurchaseRefund, error)
	UpdatePurchaseRefund(ctx context.Context, _ *UpdatePurchaseRefundArgs) (*PurchaseRefund, error)
	CancelPurchaseRefund(ctx context.Context, _ *CancelPurchaseRefundArgs) (*PurchaseRefund, error)
	ConfirmPurchaseRefund(ctx context.Context, _ *ConfirmPurchaseRefundArgs) (*PurchaseRefund, error)
}

type QueryService interface {
	ListPurchaseRefunds(context.Context, *GetPurchaseRefundsArgs) (*GetPurchaseRefundsResponse, error)
	GetPurchaseRefundByID(context.Context, *GetPurchaseRefundByIDArgs) (*PurchaseRefund, error)
	GetPurchaseRefundsByIDs(context.Context, *GetPurchaseRefundsByIDsArgs) ([]*PurchaseRefund, error)
	GetPurchaseRefundsByPurchaseOrderID(context.Context, *GetPurchaseRefundsByPurchaseOrderIDRequest) ([]*PurchaseRefund, error)
}

type GetPurchaseRefundsByPurchaseOrderIDRequest struct {
	PurchaseOrderID dot.ID
	ShopID          dot.ID
}

type GetPurchaseRefundsResponse struct {
	PageInfo        meta.PageInfo
	PurchaseRefunds []*PurchaseRefund
	Count           int
}

// +convert:create=PurchaseRefund
type CreatePurchaseRefundArgs struct {
	Lines           []*PurchaseRefundLine
	PurchaseOrderID dot.ID
	AdjustmentLines []*types.AdjustmentLine
	TotalAdjustment int
	TotalAmount     int
	BasketValue     int
	ShopID          dot.ID
	CreatedBy       dot.ID
	Note            string
}

// +convert:update=PurchaseRefund
type UpdatePurchaseRefundArgs struct {
	Lines           []*PurchaseRefundLine
	ID              dot.ID
	ShopID          dot.ID
	AdjustmentLines []*types.AdjustmentLine
	TotalAdjustment dot.NullInt
	TotalAmount     dot.NullInt
	UpdateBy        dot.ID
	BasketValue     dot.NullInt
	Note            dot.NullString
}

type CancelPurchaseRefundArgs struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	CancelReason         string
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	InventoryOverStock   bool
}

type ConfirmPurchaseRefundArgs struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	InventoryOverStock   bool
}

type GetPurchaseRefundByIDArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetPurchaseRefundsByIDsArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
}

type GetPurchaseRefundsArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}
