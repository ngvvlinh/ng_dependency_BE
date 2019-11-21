package inventory

import (
	"context"
	"time"

	"etop.vn/capi/dot"

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

	CreateInventoryVariant(context.Context, *CreateInventoryVariantArgs) (*InventoryVariant, error)

	CheckInventoryVariantsQuantity(context.Context, *CheckInventoryVariantQuantityRequest) error

	CreateInventoryVoucherByQuantityChange(context.Context, *CreateInventoryVoucherByQuantityChangeRequest) (*CreateInventoryVoucherByQuantityChangeResponse, error)

	CreateInventoryVoucherByReference(context.Context, *CreateInventoryVoucherByReferenceArgs) ([]*InventoryVoucher, error)
}

type QueryService interface {
	GetInventoryVariant(_ context.Context, ShopID int64, VariantID int64) (*InventoryVariant, error)

	GetInventoryVariants(context.Context, *GetInventoryRequest) (*GetInventoryVariantsResponse, error)

	GetInventoryVoucher(_ context.Context, ShopID int64, ID int64) (*InventoryVoucher, error)

	GetInventoryVouchers(ctx context.Context, _ *ListInventoryVouchersArgs) (*GetInventoryVouchersResponse, error)

	GetInventoryVariantsByVariantIDs(context.Context, *GetInventoryVariantsByVariantIDsArgs) (*GetInventoryVariantsResponse, error)

	GetInventoryVouchersByIDs(context.Context, *GetInventoryVouchersByIDArgs) (*GetInventoryVouchersResponse, error)

	GetInventoryVouchersByRefIDs(_ context.Context, RefIDs []int64, ShopID int64) (*GetInventoryVouchersResponse, error)

	GetInventoryVoucherByReference(ctx context.Context, ShopID int64, refID int64, refType InventoryRefType) (*GetInventoryVoucherByReferenceResponse, error)
}

// +convert:update=InventoryVoucher
type UpdateInventoryVoucherArgs struct {
	ID        int64
	ShopID    int64
	Title     dot.NullString
	UpdatedBy int64

	TraderID    dot.NullInt64
	TotalAmount int32

	Note  dot.NullString
	Lines []*InventoryVoucherItem
}

type ListInventoryVouchersArgs struct {
	ShopID  int64
	Paging  meta.Paging
	Filters meta.Filters
}

type CreateInventoryVariantArgs struct {
	ShopID    int64
	VariantID int64
}

type GetInventoryRequest struct {
	ShopID int64
	Paging *meta.Paging
}

type GetInventoryVariantsResponse struct {
	InventoryVariants []*InventoryVariant
}

type GetInventoryVariantsByVariantIDsArgs struct {
	ShopID     int64
	Paging     *meta.Paging
	VariantIDs []int64
}

type AdjustInventoryQuantityRespone struct {
	InventoryVariants []*InventoryVariant
	InventoryVouchers []*InventoryVoucher
}

type InventoryVariant struct {
	ShopID          int64
	VariantID       int64
	QuantityOnHand  int32
	QuantityPicked  int32
	PurchasePrice   int32
	QuantitySummary int32

	CreatedAt time.Time
	UpdatedAt time.Time
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
	ShopID    int64
	CreatedBy int64
	Title     string

	RefID   int64
	RefType InventoryRefType
	RefName InventoryVoucherRefName
	RefCode string

	TraderID    int64
	TotalAmount int32
	Type        InventoryVoucherType
	Note        string
	Lines       []*InventoryVoucherItem
}

type CreateInventoryVoucherByQuantityChangeRequest struct {
	ShopID int64

	RefID   int64
	RefType InventoryRefType
	RefName InventoryVoucherRefName
	RefCode string

	Note  string
	Title string

	Overstock bool

	CreatedBy int64
	Variants  []*InventoryVariantQuantityChange
}

type InventoryVariantQuantityChange struct {
	VariantID      int64
	QuantityChange int32
}

type CheckInventoryVariantQuantityRequest struct {
	Lines              []*InventoryVoucherItem
	InventoryOverStock bool
	ShopID             int64
	Type               InventoryVoucherType
}

type InventoryVoucher struct {
	ID        int64
	ShopID    int64
	Title     string
	Code      string
	CodeNorm  int32
	CreatedBy int64
	UpdatedBy int64

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	RefID   int64
	RefType InventoryRefType
	RefName InventoryVoucherRefName
	RefCode string

	TraderID    int64
	Trader      *Trader
	TotalAmount int32

	// enum "in" or "out"
	Type InventoryVoucherType

	CancelReason string
	Note         string
	Lines        []*InventoryVoucherItem
	Status       etop.Status3
}

type Trader struct {
	ID       int64
	Type     string
	FullName string
	Phone    string
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

type GetInventoryVoucherByReferenceResponse struct {
	InventoryVouchers []*InventoryVoucher
	Status            etop.Status4
}

type CreateInventoryVoucherByQuantityChangeResponse struct {
	TypeIn  *InventoryVoucher
	TypeOut *InventoryVoucher
}

type CreateInventoryVoucherByReferenceArgs struct {
	RefType   InventoryRefType
	RefID     int64
	Type      InventoryVoucherType
	ShopID    int64
	UserID    int64
	OverStock bool
}
