package purchaseorder

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreatePurchaseOrder(ctx context.Context, _ *CreatePurchaseOrderArgs) (*PurchaseOrder, error)
	UpdatePurchaseOrder(ctx context.Context, _ *UpdatePurchaseOrderArgs) (*PurchaseOrder, error)
	CancelPurchaseOrder(ctx context.Context, _ *CancelPurchaseOrderArgs) (updated int, _ error)
	ConfirmPurchaseOrder(ctx context.Context, _ *ConfirmPurchaseOrderArgs) (updated int, _ error)
	DeletePurchaseOrder(ctx context.Context, ID, shopID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	GetPurchaseOrderByID(context.Context, *shopping.IDQueryShopArg) (*PurchaseOrder, error)
	GetPurchaseOrdersByIDs(ctx context.Context, IDs []dot.ID, ShopID dot.ID) (*PurchaseOrdersResponse, error)
	ListPurchaseOrders(context.Context, *shopping.ListQueryShopArgs) (*PurchaseOrdersResponse, error)
	ListPurchaseOrdersBySupplierIDsAndStatuses(ctx context.Context, shopID dot.ID, supplierIDs []dot.ID, statuses []status3.Status) (*PurchaseOrdersResponse, error)
	ListPurchaseOrdersByReceiptID(ctx context.Context, receiptID, shopID dot.ID) (*PurchaseOrdersResponse, error)
}

//-- queries --//

type PurchaseOrdersResponse struct {
	PurchaseOrders []*PurchaseOrder
	Paging         meta.PageInfo
}

//-- commands --//

// +convert:create=PurchaseOrder
type CreatePurchaseOrderArgs struct {
	ShopID        dot.ID
	SupplierID    dot.ID
	BasketValue   int
	TotalDiscount int
	TotalAmount   int
	Note          string
	Lines         []*PurchaseOrderLine
	CreatedBy     dot.ID
}

// +convert:update=PurchaseOrder(ID,ShopID)
type UpdatePurchaseOrderArgs struct {
	ID            dot.ID
	ShopID        dot.ID
	BasketValue   NullInt
	TotalDiscount NullInt
	TotalAmount   NullInt
	Note          NullString
	Lines         []*PurchaseOrderLine
}

type CancelPurchaseOrderArgs struct {
	ID                   dot.ID
	ShopID               dot.ID
	Reason               string
	UpdatedBy            dot.ID
	InventoryOverStock   bool
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type ConfirmPurchaseOrderArgs struct {
	ID                   dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	ShopID               dot.ID
}
