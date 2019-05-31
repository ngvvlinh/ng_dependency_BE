package ordering

import (
	"context"
	"time"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/ordering/types"
	ordertypes "etop.vn/api/main/ordering/types"
	shippingtypes "etop.vn/api/main/shipping/types"
)

type Aggregate interface {
	ValidateOrdersForShipping(ctx context.Context, args *ValidateOrdersForShippingArgs) (*ValidateOrdersForShippingResponse, error)
	ReserveOrdersForFfm(ctx context.Context, args *ReserveOrdersForFfmArgs) (*ReserveOrdersForFfmResponse, error)
	ReleaseOrdersForFfm(ctx context.Context, args *ReleaseOrdersForFfmArgs) (*ReleaseOrdersForFfmResponse, error)
	UpdateOrderShippingStatus(ctx context.Context, args *UpdateOrderShippingStatusArgs) (*UpdateOrderShippingStatusResponse, error)
	UpdateOrdersConfirmStatus(ctx context.Context, args *UpdateOrdersConfirmStatusArgs) (*UpdateOrdersConfirmStatusResponse, error)
}

type QueryService interface {
	GetOrderByID(ctx context.Context, args *GetOrderByIDArgs) (*Order, error)
	GetOrders(ctx context.Context, args *GetOrdersArgs) (*OrdersResponse, error)
}

type GetOrderByIDArgs struct {
	ID int64
}

type GetOrdersArgs struct {
	ShopID int64
	IDs    []int64
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

	Fulfill    ordertypes.Fulfill
	FulfillIDs []int64

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProcessedAt time.Time
	ClosedAt    time.Time
	ConfirmedAt time.Time
	CancelledAt time.Time
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
