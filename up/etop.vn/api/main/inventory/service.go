package inventory

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	CreateInventoryVoucher(_ context.Context, Overstock bool, _ *CreateInventoryVoucherArgs) (*InventoryVoucher, error)

	ConfirmInventoryVoucher(context.Context, *ConfirmInventoryVoucherArgs) (*InventoryVoucher, error)

	CancelInventoryVoucher(context.Context, *CancelInventoryVoucherArgs) (*InventoryVoucher, error)

	AdjustInventoryQuantity(_ context.Context, Overstock bool, _ *AdjustInventoryQuantityArgs) (*AdjustInventoryQuantityRespone, error)

	UpdateInventoryVoucher(context.Context, *UpdateInventoryVoucherArgs) (*InventoryVoucher, error)

	CreateInventoryVariant(context.Context, *CreateInventoryVariantArgs) error
}

type QueryService interface {
	GetInventoryVoucher(_ context.Context, ShopID int64, ID int64) (*InventoryVoucher, error)

	GetInventory(_ context.Context, ShopID int64, VariantID int64) (*InventoryVariant, error)

	GetInventories(context.Context, *GetInventoryRequest) (*GetInventoriesResponse, error)

	GetInventoryVouchers(_ context.Context, ShopID int64, Paging *meta.Paging) (*GetInventoryVouchersResponse, error)

	GetInventoriesByVariantIDs(context.Context, *GetInventoriesByVariantIDsArgs) (*GetInventoriesResponse, error)

	GetInventoryVouchersByIDs(context.Context, *GetInventoryVouchersByIDArgs) (*GetInventoryVouchersResponse, error)
}

// +convert:update=InventoryVoucher
type UpdateInventoryVoucherArgs struct {
	ID        int64
	ShopID    int64
	Title     string
	UpdatedBy int64

	TraderID    int64
	TotalAmount int32

	CancelledReason string
	Note            string
	Lines           []*InventoryVoucherItem
}

type CreateInventoryVariantArgs struct {
	ShopID    int64
	VariantID int64
}

type GetInventoryRequest struct {
	ShopID int64
	Paging *meta.Paging
}

type GetInventoriesResponse struct {
	Inventories []*InventoryVariant
}

type GetInventoriesByVariantIDsArgs struct {
	ShopID     int64
	Paging     *meta.Paging
	VariantIDs []int64
}

type AdjustInventoryQuantityRespone struct {
	Inventory         []*InventoryVariant
	InventoryVouchers []*InventoryVoucher
}

type InventoryVariant struct {
	ShopID          int64
	VariantID       int64
	QuantityOnHand  int32
	QuantityPicked  int32
	PurchasePrice   int32
	QuantitySummary int32
}

type GetInventoryVouchersResponse struct {
	InventoryVoucher []*InventoryVoucher
}

type ConfirmInventoryVoucherArgs struct {
	ShopID    int64
	ID        int64
	UpdatedBy int64
}

type GetInventoryVouchersArgs struct {
	ShopID int64
	Paging *meta.Paging
}

type GetInventoryVouchersByIDArgs struct {
	ShopID int64
	Paging *meta.Paging
	IDs    []int64
}

// +convert:create=InventoryVoucher
type CreateInventoryVoucherArgs struct {
	ShopID      int64
	CreatedBy   int64
	Title       string
	TraderID    int64
	TotalAmount int32
	Type        string
	Note        string
	Lines       []*InventoryVoucherItem
}

type InventoryVoucher struct {
	ID     int64
	ShopID int64
	Title  string

	CreatedBy int64
	UpdatedBy int64

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	TraderID    int64
	TotalAmount int32

	// enum "in" or "out"
	Type string

	CancelledReason string
	Note            string
	Lines           []*InventoryVoucherItem
	Status          etop.Status3
}

type InventoryVoucherItem struct {
	VariantID int64
	Price     int32
	Quantity  int32
}

type CancelInventoryVoucherArgs struct {
	ShopID    int64
	ID        int64
	UpdatedBy int64
	Reason    string
}

type AdjustInventoryQuantityArgs struct {
	ShopID int64
	Lines  []*InventoryVariant
	Title  string
	UserID int64
	Note   string
}

type InventoryVoucherConfirmEvent struct {
	ShopID int64
	Line   []*InventoryVoucherItem
}
