package purchaseorder

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreatePurchaseOrder(ctx context.Context, _ *CreatePurchaseOrderArgs) (*PurchaseOrder, error)
	UpdatePurchaseOrder(ctx context.Context, _ *UpdatePurchaseOrderArgs) (*PurchaseOrder, error)
	CancelPurchaseOrder(ctx context.Context, _ *CancelPurchaseOrderArgs) (updated int, _ error)
	ConfirmPurchaseOrder(ctx context.Context, _ *ConfirmPurchaseOrderArgs) (updated int, _ error)
	DeletePurchaseOrder(ctx context.Context, ID, shopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetPurchaseOrderByID(context.Context, *shopping.IDQueryShopArg) (*PurchaseOrder, error)
	GetPurchaseOrdersByIDs(ctx context.Context, IDs []int64, ShopID int64) (*PurchaseOrdersResponse, error)
	ListPurchaseOrders(context.Context, *shopping.ListQueryShopArgs) (*PurchaseOrdersResponse, error)
	ListPurchaseOrdersBySupplierIDsAndStatuses(ctx context.Context, shopID int64, supplierIDs []int64, statuses []etop.Status3) (*PurchaseOrdersResponse, error)
	ListPurchaseOrdersByReceiptID(ctx context.Context, receiptID, shopID int64) (*PurchaseOrdersResponse, error)
}

//-- queries --//

type PurchaseOrdersResponse struct {
	PurchaseOrders []*PurchaseOrder
	Count          int32
	Paging         meta.PageInfo
}

//-- commands --//

// +convert:create=PurchaseOrder
type CreatePurchaseOrderArgs struct {
	ShopID        int64
	SupplierID    int64
	BasketValue   int64
	TotalDiscount int64
	TotalAmount   int64
	Note          string
	Lines         []*PurchaseOrderLine
	CreatedBy     int64
}

// +convert:update=PurchaseOrder(ID,ShopID)
type UpdatePurchaseOrderArgs struct {
	ID            int64
	ShopID        int64
	SupplierID    NullInt64
	BasketValue   NullInt64
	TotalDiscount NullInt64
	TotalAmount   NullInt64
	Note          NullString
	Lines         []*PurchaseOrderLine
}

type CancelPurchaseOrderArgs struct {
	ID     int64
	ShopID int64
	Reason string
}

type ConfirmPurchaseOrderArgs struct {
	ID                   int64
	AutoInventoryVoucher inventory.AutoInventoryVoucher
	ShopID               int64
}
