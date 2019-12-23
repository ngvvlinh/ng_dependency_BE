package refund

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateRefund(ctx context.Context, _ *CreateRefundArgs) (*Refund, error)
	UpdateRefund(ctx context.Context, _ *UpdateRefundArgs) (*Refund, error)
	CancelRefund(ctx context.Context, _ *CancelRefundArgs) (*Refund, error)
	ConfirmRefund(ctx context.Context, _ *ConfirmRefundArgs) (*Refund, error)
}

type QueryService interface {
	GetRefunds(context.Context, *GetRefundsArgs) (*GetRefundsResponse, error)
	GetRefundByID(context.Context, *GetRefundByIDArgs) (*Refund, error)
	GetRefundsByIDs(context.Context, *GetRefundsByIDsArgs) ([]*Refund, error)
}

type GetRefundsResponse struct {
	PageInfor meta.PageInfo
	Refunds   []*Refund
}

// +convert:create=Refund
type CreateRefundArgs struct {
	Lines     []*RefundLine
	OrderID   dot.ID
	Discount  int
	ShopID    dot.ID
	CreatedBy dot.ID
	Note      string
}

// +convert:update=Refund
type UpdateRefundArgs struct {
	Lines    []*RefundLine
	ID       dot.ID
	ShopID   dot.ID
	Discount dot.NullInt
	UpdateBy dot.ID
	Note     dot.NullString
}

type CancelRefundArgs struct {
	ShopID       dot.ID
	ID           dot.ID
	UpdatedBy    dot.ID
	CancelReason string
}

type ConfirmRefundArgs struct {
	ShopID               dot.ID
	ID                   dot.ID
	UpdatedBy            dot.ID
	AutoInventoryVoucher inventory.AutoInventoryVoucher
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
