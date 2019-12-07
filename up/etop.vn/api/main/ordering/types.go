package ordering

import (
	"context"
	"time"

	"etop.vn/api/top/types/etc/status5"

	"etop.vn/api/top/types/etc/status4"

	"etop.vn/api/top/types/etc/status3"

	"etop.vn/api/main/inventory"
	"etop.vn/api/main/ordering/types"
	ordertypes "etop.vn/api/main/ordering/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:api
// +gen:event:topic=event/ordering

type Aggregate interface {
	ValidateOrdersForShipping(context.Context, *ValidateOrdersForShippingArgs) (*ValidateOrdersForShippingResponse, error)
	ReserveOrdersForFfm(context.Context, *ReserveOrdersForFfmArgs) (*ReserveOrdersForFfmResponse, error)
	ReleaseOrdersForFfm(context.Context, *ReleaseOrdersForFfmArgs) (*ReleaseOrdersForFfmResponse, error)
	UpdateOrderShippingStatus(context.Context, *UpdateOrderShippingStatusArgs) (*UpdateOrderShippingStatusResponse, error)
	UpdateOrdersConfirmStatus(context.Context, *UpdateOrdersConfirmStatusArgs) (*UpdateOrdersConfirmStatusResponse, error)

	UpdateOrderPaymentInfo(context.Context, *UpdateOrderPaymentInfoArgs) error
}

type QueryService interface {
	GetOrderByID(context.Context, *GetOrderByIDArgs) (*Order, error)
	GetOrders(context.Context, *GetOrdersArgs) (*OrdersResponse, error)
	GetOrdersByIDsAndCustomerID(ctx context.Context, shopID dot.ID, IDs []dot.ID, customerID dot.ID) (*OrdersResponse, error)
	GetOrderByCode(context.Context, *GetOrderByCodeArgs) (*Order, error)
	ListOrdersByCustomerID(ctx context.Context, shopID, customerID dot.ID) (*OrdersResponse, error)
	ListOrdersByCustomerIDs(ctx context.Context, shopID dot.ID, customerIDs []dot.ID) (*OrdersResponse, error)
}

type GetOrderByIDArgs struct {
	ID dot.ID
}

type GetOrdersArgs struct {
	ShopID dot.ID
	IDs    []dot.ID
}

type GetOrderByCodeArgs struct {
	Code string
}

type OrdersResponse struct {
	Orders []*Order
}

type ValidateOrdersForShippingResponse struct {
}

type UpdateOrderShippingStatusResponse struct{}

type UpdateOrdersConfirmStatusResponse struct{}

type Order struct {
	ID              dot.ID
	ShopID          dot.ID
	Code            string
	CustomerAddress *types.Address
	ShippingAddress *types.Address
	CancelReason    string

	ConfirmStatus             status3.Status
	Status                    status5.Status
	FulfillmentShippingStatus status5.Status
	EtopPaymentStatus         status4.Status

	Lines         []*types.ItemLine
	TotalItems    int
	BasketValue   int
	TotalWeight   int
	OrderDiscount int
	TotalDiscount int
	TotalFee      int
	TotalAmount   int

	OrderNote string
	FeeLines  []OrderFeeLine
	Shipping  *shippingtypes.ShippingInfo

	FulfillmentType ordertypes.Fulfill
	FulfillmentIDs  []dot.ID

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProcessedAt time.Time
	ClosedAt    time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	PaymentStatus status4.Status
	PaymentID     dot.ID
	ReferralMeta  *ReferralMeta

	TradingShopID dot.ID
	CustomerID    dot.ID
}

type OrderFeeLine struct {
	Type   string
	Name   string
	Code   string
	Desc   string
	Amount int
}

// shipping means manually fulfill, Fulfillment or ShipnowFulfillment
type ValidateOrdersForShippingArgs struct {
	OrderIDs []dot.ID
}

type ReserveOrdersForFfmArgs struct {
	OrderIDs   []dot.ID
	Fulfill    ordertypes.Fulfill
	FulfillIDs []dot.ID
}

type ReserveOrdersForFfmResponse struct {
	Orders []*Order
}

type ReleaseOrdersForFfmArgs struct {
	OrderIDs []dot.ID
}

type ReleaseOrdersForFfmResponse struct {
	Updated int
}

type UpdateOrderShippingStatusArgs struct {
	ID                         dot.ID
	FulfillmentShippingStates  []string
	FulfillmentShippingStatus  status5.Status
	FulfillmentPaymentStatuses []int
	EtopPaymentStatus          status4.Status
	CODEtopPaidAt              time.Time
}

type UpdateOrdersConfirmStatusArgs struct {
	IDs           []dot.ID
	ShopConfirm   status3.Status
	ConfirmStatus status3.Status
}

type UpdateOrderPaymentInfoArgs struct {
	ID            dot.ID
	PaymentStatus status4.Status
	PaymentID     dot.ID
}

type ReferralMeta struct {
	ReferralCode   string `json:"referral_code"`
	ReferralAmount string `json:"referral_amount"`
}

type OrderPaymentSuccessEvent struct {
	meta.EventMeta
	OrderID dot.ID
}

type OrderConfirmingEvent struct {
	OrderID              dot.ID
	ShopID               dot.ID
	InventoryOverStock   bool
	Lines                []*types.ItemLine
	CustomerID           dot.ID
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type OrderConfirmedEvent struct {
	OrderID              dot.ID
	ShopID               dot.ID
	InventoryOverStock   bool
	AutoInventoryVoucher inventory.AutoInventoryVoucher
	UpdatedBy            dot.ID
}

type OrderCancelledEvent struct {
	OrderID              dot.ID
	ShopID               dot.ID
	AutoInventoryVoucher inventory.AutoInventoryVoucher
	UpdatedBy            dot.ID
}
