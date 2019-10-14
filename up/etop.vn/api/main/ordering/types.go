package ordering

import (
	"context"
	"time"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering/types"
	ordertypes "etop.vn/api/main/ordering/types"
	shippingtypes "etop.vn/api/main/shipping/types"
)

// +gen:api

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
	GetOrdersByIDsAndCustomerID(ctx context.Context, shopID int64, IDs []int64, customerID int64) (*OrdersResponse, error)
	GetOrderByCode(context.Context, *GetOrderByCodeArgs) (*Order, error)
}

type GetOrderByIDArgs struct {
	ID int64
}

type GetOrdersArgs struct {
	ShopID int64
	IDs    []int64
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
	ID              int64
	ShopID          int64
	Code            string
	CustomerAddress *types.Address
	ShippingAddress *types.Address
	CancelReason    string

	ConfirmStatus             etoptypes.Status3
	Status                    etoptypes.Status5
	FulfillmentShippingStatus etoptypes.Status5
	EtopPaymentStatus         etoptypes.Status4

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
	FulfillmentIDs  []int64

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProcessedAt time.Time
	ClosedAt    time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time

	PaymentStatus etoptypes.Status4
	PaymentID     int64
	ReferralMeta  *ReferralMeta
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
	OrderIDs []int64
}

type ReserveOrdersForFfmArgs struct {
	OrderIDs   []int64
	Fulfill    ordertypes.Fulfill
	FulfillIDs []int64
}

type ReserveOrdersForFfmResponse struct {
	Orders []*Order
}

type ReleaseOrdersForFfmArgs struct {
	OrderIDs []int64
}

type ReleaseOrdersForFfmResponse struct {
	Updated int
}

type UpdateOrderShippingStatusArgs struct {
	ID                         int64
	FulfillmentShippingStates  []string
	FulfillmentShippingStatus  etoptypes.Status5
	FulfillmentPaymentStatuses []int
	EtopPaymentStatus          etoptypes.Status4
	CODEtopPaidAt              time.Time
}

type UpdateOrdersConfirmStatusArgs struct {
	IDs           []int64
	ShopConfirm   etoptypes.Status3
	ConfirmStatus etoptypes.Status3
}

type UpdateOrderPaymentInfoArgs struct {
	ID            int64
	PaymentStatus etoptypes.Status4
	PaymentID     int64
}

type ReferralMeta struct {
	ReferralCode   string `json:"referral_code"`
	ReferralAmount string `json:"referral_amount"`
}
