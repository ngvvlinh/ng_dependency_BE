package inventory

import (
	"context"
	"time"

	catalogtype "etop.vn/api/main/catalog/types"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/inventory_type"
	"etop.vn/api/top/types/etc/inventory_voucher_ref"
	"etop.vn/api/top/types/etc/ref_action"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateInventoryVoucher(_ context.Context, Overstock bool, _ *CreateInventoryVoucherArgs) (*InventoryVoucher, error)

	ConfirmInventoryVoucher(context.Context, *ConfirmInventoryVoucherArgs) (*InventoryVoucher, error)

	CancelInventoryVoucher(context.Context, *CancelInventoryVoucherArgs) (*InventoryVoucher, error)

	AdjustInventoryQuantity(_ context.Context, Overstock bool, _ *AdjustInventoryQuantityArgs) (*AdjustInventoryQuantityRespone, error)

	UpdateInventoryVoucher(context.Context, *UpdateInventoryVoucherArgs) (*InventoryVoucher, error)

	CreateInventoryVariant(context.Context, *CreateInventoryVariantArgs) (*InventoryVariant, error)

	CheckInventoryVariantsQuantity(context.Context, *CheckInventoryVariantQuantityRequest) error

	CreateInventoryVoucherByQuantityChange(context.Context, *CreateInventoryVoucherByQuantityChangeRequest) (*CreateInventoryVoucherByQuantityChangeResponse, error)

	CreateInventoryVoucherByReference(context.Context, *CreateInventoryVoucherByReferenceArgs) ([]*InventoryVoucher, error)

	UpdateInventoryVariantCostPrice(context.Context, *UpdateInventoryVariantCostPriceRequest) (*InventoryVariant, error)

	CancelInventoryByRefID(context.Context, *CancelInventoryByRefIDRequest) (*CancelInventoryByRefIDResponse, error)
}

type QueryService interface {
	GetInventoryVariant(_ context.Context, ShopID dot.ID, VariantID dot.ID) (*InventoryVariant, error)

	GetInventoryVariants(context.Context, *GetInventoryRequest) (*GetInventoryVariantsResponse, error)

	GetInventoryVoucher(_ context.Context, ShopID dot.ID, ID dot.ID) (*InventoryVoucher, error)

	GetInventoryVouchers(ctx context.Context, _ *ListInventoryVouchersArgs) (*GetInventoryVouchersResponse, error)

	GetInventoryVariantsByVariantIDs(context.Context, *GetInventoryVariantsByVariantIDsArgs) (*GetInventoryVariantsResponse, error)

	GetInventoryVouchersByIDs(context.Context, *GetInventoryVouchersByIDArgs) (*GetInventoryVouchersResponse, error)

	GetInventoryVouchersByRefIDs(_ context.Context, RefIDs []dot.ID, ShopID dot.ID) (*GetInventoryVouchersResponse, error)

	GetInventoryVoucherByReference(ctx context.Context, ShopID dot.ID, refID dot.ID, refType inventory_voucher_ref.InventoryVoucherRef) (*GetInventoryVoucherByReferenceResponse, error)

	ListInventoryVariantsByVariantIDs(context.Context, *ListInventoryVariantsByVariantIDsArgs) (*GetInventoryVariantsResponse, error)
}

// +convert:update=InventoryVoucher
type UpdateInventoryVoucherArgs struct {
	ID        dot.ID
	ShopID    dot.ID
	Title     dot.NullString
	UpdatedBy dot.ID

	TraderID    dot.NullID
	TotalAmount int

	Note  dot.NullString
	Lines []*InventoryVoucherItem
}

type ListInventoryVouchersArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type CreateInventoryVariantArgs struct {
	ShopID    dot.ID
	VariantID dot.ID
}

type GetInventoryRequest struct {
	ShopID dot.ID
	Paging *meta.Paging
}

type GetInventoryVariantsResponse struct {
	InventoryVariants []*InventoryVariant
}

type GetInventoryVariantsByVariantIDsArgs struct {
	ShopID     dot.ID
	Paging     *meta.Paging
	VariantIDs []dot.ID
}

type AdjustInventoryQuantityRespone struct {
	InventoryVariants []*InventoryVariant
	InventoryVouchers []*InventoryVoucher
}

type InventoryVariant struct {
	ShopID          dot.ID
	VariantID       dot.ID
	QuantityOnHand  int
	QuantityPicked  int
	CostPrice       int
	QuantitySummary int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetInventoryVouchersResponse struct {
	InventoryVoucher []*InventoryVoucher
}

type ConfirmInventoryVoucherArgs struct {
	ShopID    dot.ID
	ID        dot.ID
	UpdatedBy dot.ID
}

type GetInventoryVouchersArgs struct {
	ShopID dot.ID
	Paging *meta.Paging
}

type GetInventoryVouchersByIDArgs struct {
	ShopID dot.ID
	Paging *meta.Paging
	IDs    []dot.ID
}

// +convert:create=InventoryVoucher
type CreateInventoryVoucherArgs struct {
	ShopID    dot.ID
	CreatedBy dot.ID
	Title     string

	RefID   dot.ID
	RefType inventory_voucher_ref.InventoryVoucherRef
	RefCode string

	TraderID    dot.ID
	TotalAmount int
	Type        inventory_type.InventoryVoucherType
	Lines       []*InventoryVoucherItem
}

type CreateInventoryVoucherByQuantityChangeRequest struct {
	ShopID dot.ID

	RefID   dot.ID
	RefType inventory_voucher_ref.InventoryVoucherRef
	RefCode string

	Title string

	Overstock bool

	CreatedBy dot.ID
	Lines     []*InventoryVariantQuantityChange
}

type InventoryVariantQuantityChange struct {
	ItemInfo       *InventoryVoucherItem
	QuantityChange int
}

type CheckInventoryVariantQuantityRequest struct {
	Lines              []*InventoryVoucherItem
	InventoryOverStock bool
	ShopID             dot.ID
	Type               inventory_type.InventoryVoucherType
}

type InventoryVoucher struct {
	ID        dot.ID
	ShopID    dot.ID
	Title     string
	Code      string
	CodeNorm  int
	CreatedBy dot.ID
	UpdatedBy dot.ID

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	RefID     dot.ID
	RefType   inventory_voucher_ref.InventoryVoucherRef
	RefName   string
	RefCode   string
	RefAction ref_action.RefAction

	TraderID    dot.ID
	Trader      *Trader
	TotalAmount int

	// enum "in" or "out"
	Type inventory_type.InventoryVoucherType

	CancelReason string
	Note         string
	Lines        []*InventoryVoucherItem
	Status       status3.Status
}

type Trader struct {
	ID       dot.ID
	Type     string
	FullName string
	Phone    string
}

type InventoryVoucherItem struct {
	ProductID   dot.ID
	ProductName string
	VariantID   dot.ID
	VariantName string

	Quantity int
	Price    int

	Code       string
	ImageURL   string
	Attributes []*catalogtype.Attribute
}

type CancelInventoryVoucherArgs struct {
	ShopID       dot.ID
	ID           dot.ID
	UpdatedBy    dot.ID
	CancelReason string
}

type AdjustInventoryQuantityArgs struct {
	ShopID dot.ID
	Lines  []*InventoryVariant
	Title  string
	UserID dot.ID
	Note   string
}

type InventoryVoucherConfirmEvent struct {
	ShopID dot.ID
	Line   []*InventoryVoucherItem
}

type GetInventoryVoucherByReferenceResponse struct {
	InventoryVouchers []*InventoryVoucher
	Status            status4.Status
}

type CreateInventoryVoucherByQuantityChangeResponse struct {
	TypeIn  *InventoryVoucher
	TypeOut *InventoryVoucher
}

type CreateInventoryVoucherByReferenceArgs struct {
	RefType   inventory_voucher_ref.InventoryVoucherRef
	RefID     dot.ID
	Type      inventory_type.InventoryVoucherType
	ShopID    dot.ID
	UserID    dot.ID
	OverStock bool
}

type UpdateInventoryVariantCostPriceRequest struct {
	ShopID    dot.ID
	VariantID dot.ID
	CostPrice int
}

type CancelInventoryByRefIDResponse struct {
	InventoryVouchers []*InventoryVoucher
}

type CancelInventoryByRefIDRequest struct {
	RefID                dot.ID
	ShopID               dot.ID
	RefType              inventory_voucher_ref.InventoryVoucherRef
	InventoryOverStock   bool
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	UpdateBy             dot.ID
}

type ListInventoryVariantsByVariantIDsArgs struct {
	ShopID     dot.ID
	VariantIDs []dot.ID
}
