package ordering

import (
	"context"
	"time"

	"o.o/api/main/ordering/types"
	ordertypes "o.o/api/main/ordering/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/fee"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/order_source"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +gen:api
// +gen:event:topic=event/ordering

type Aggregate interface {
	ValidateOrdersForShipping(context.Context, *ValidateOrdersForShippingArgs) (*ValidateOrdersForShippingResponse, error)
	ReserveOrdersForFfm(context.Context, *ReserveOrdersForFfmArgs) (*ReserveOrdersForFfmResponse, error)
	ReleaseOrdersForFfm(context.Context, *ReleaseOrdersForFfmArgs) (*ReleaseOrdersForFfmResponse, error)
	UpdateOrderShippingStatus(context.Context, *UpdateOrderShippingStatusArgs) (*UpdateOrderShippingStatusResponse, error)
	UpdateOrdersConfirmStatus(context.Context, *UpdateOrdersConfirmStatusArgs) (*UpdateOrdersConfirmStatusResponse, error)
	UpdateOrderCustomerInfo(context.Context, *UpdateOrderCustomerInfoArgs) error

	UpdateOrderPaymentInfo(context.Context, *UpdateOrderPaymentInfoArgs) error
	CompleteOrder(_ context.Context, OrderID dot.ID, ShopID dot.ID) error
	UpdateOrderStatus(context.Context, *UpdateOrderStatusArgs) error
	UpdateOrderPaymentStatus(context.Context, *UpdateOrderPaymentStatusArgs) error
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
	PartnerID       dot.ID
	Code            string
	EdCode          string
	Customer        *OrderCustomer
	CustomerAddress *types.Address
	BillingAddress  *types.Address
	ShippingAddress *types.Address
	CancelReason    string

	ConfirmStatus             status3.Status
	ShopConfirm               status3.Status
	Status                    status5.Status
	FulfillmentShippingStatus status5.Status
	EtopPaymentStatus         status4.Status

	Lines           []*types.ItemLine
	Discounts       []*OrderDiscount
	TotalItems      int
	BasketValue     int
	TotalWeight     int
	OrderDiscount   int
	TotalDiscount   int
	TotalFee        int
	TotalAmount     int
	ShopCOD         int
	ShopShippingFee int

	OrderNote    string
	FeeLines     []OrderFeeLine
	Shipping     *shippingtypes.ShippingInfo
	ShippingNote string

	FulfillmentType ordertypes.ShippingType
	FulfillmentIDs  []dot.ID

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProcessedAt time.Time
	ClosedAt    time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	CreatedBy dot.ID

	PaymentStatus status4.Status
	PaymentID     dot.ID
	ReferralMeta  *ReferralMeta

	TradingShopID dot.ID
	CustomerID    dot.ID

	OrderSourceType order_source.Source
	ExternalOrderID string
	PaymentMethod   payment_method.PaymentMethod
	ReferenceURL    string
	GhnNoteCode     ghn_note_code.GHNNoteCode
	TryOn           try_on.TryOnCode

	PreOrder bool
}

type OrderFeeLine struct {
	Type   fee.FeeType
	Name   string
	Code   string
	Desc   string
	Amount int
}

type OrderCustomer struct {
	FirstName     string
	LastName      string
	FullName      string
	Email         string
	Phone         string
	Gender        string
	Birthday      string
	VerifiedEmail bool
	ExternalID    string
}

func (m *OrderCustomer) GetFullName() string {
	if m.FullName != "" {
		return m.FullName
	}
	return m.FirstName + " " + m.LastName
}

type OrderDiscount struct {
	Code   string
	Type   string
	Amount int
}

func (m *Order) GetTotalFee() int {
	if m.TotalFee == 0 && m.ShopShippingFee != 0 {
		return m.ShopShippingFee
	}
	return m.TotalFee
}

// shipping means manually fulfill, Fulfillment or ShipnowFulfillment
type ValidateOrdersForShippingArgs struct {
	OrderIDs []dot.ID
}

type ReserveOrdersForFfmArgs struct {
	OrderIDs   []dot.ID
	Fulfill    ordertypes.ShippingType
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
	FulfillmentStatuses        []int
	EtopPaymentStatus          status4.Status
	CODEtopPaidAt              time.Time
}

type UpdateOrdersConfirmStatusArgs struct {
	IDs           []dot.ID
	ShopConfirm   status3.Status
	ConfirmStatus status3.Status
}

type UpdateOrderCustomerInfoArgs struct {
	ID       dot.ID
	FullName dot.NullString
	Phone    dot.NullString
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
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
}

type OrderConfirmedEvent struct {
	OrderID              dot.ID
	ShopID               dot.ID
	InventoryOverStock   bool
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	UpdatedBy            dot.ID
}

type OrderCancelledEvent struct {
	OrderID              dot.ID
	OrderCode            string
	ShopID               dot.ID
	AutoInventoryVoucher inventory_auto.AutoInventoryVoucher
	UpdatedBy            dot.ID
}

type UpdateOrderStatusArgs struct {
	OrderID dot.ID
	ShopID  dot.ID
	Status  status5.Status
}

type UpdateOrderPaymentStatusArgs struct {
	OrderID       dot.ID
	ShopID        dot.ID
	PaymentStatus status4.NullStatus
}
