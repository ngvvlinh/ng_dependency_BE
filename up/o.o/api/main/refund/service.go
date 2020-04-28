package refund

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateRefund(context.Context, *CreateRefundArgs) (*Refund, error)
	UpdateRefund(context.Context, *UpdateRefundArgs) (*Refund, error)
	CancelRefund(context.Context, *CancelRefundArgs) (*Refund, error)
	ConfirmRefund(context.Context, *ConfirmRefundArgs) (*Refund, error)
}

type QueryService interface {
	GetRefunds(context.Context, *GetRefundsArgs) (*GetRefundsResponse, error)
	GetRefundByID(context.Context, *GetRefundByIDArgs) (*Refund, error)
	GetRefundsByIDs(context.Context, *GetRefundsByIDsArgs) ([]*Refund, error)
	GetRefundsByOrderID(context.Context, *GetRefundsByOrderID) ([]*Refund, error)
}

type GetRefundsByOrderID struct {
	OrderID dot.ID
	ShopID  dot.ID
}

type GetRefundsResponse struct {
	PageInfo meta.PageInfo
	Refunds  []*Refund
}

// +convert:create=Refund
type CreateRefundArgs struct {
	Lines           []*RefundLine
	OrderID         dot.ID
	AdjustmentLines []*types.AdjustmentLine
	TotalAdjustment int
	TotalAmount     int
	BasketValue     int
	ShopID          dot.ID
	CreatedBy       dot.ID
	Note            string
}

// +convert:update=Refund
type UpdateRefundArgs struct {
	Lines           []*RefundLine
	ID              dot.ID
	ShopID          dot.ID
	AdjustmentLines []*types.AdjustmentLine
	TotalAdjustment dot.NullInt
	TotalAmount     dot.NullInt
	BasketValue     dot.NullInt
	UpdateBy        dot.ID
	Note            dot.NullString
}

type CancelRefundArgs struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	CancelReason         string
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type ConfirmRefundArgs struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type GetRefundByIDArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetRefundsByIDsArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
}

type GetRefundsArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}
