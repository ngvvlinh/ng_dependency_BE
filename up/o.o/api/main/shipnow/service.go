package shipnow

import (
	"context"
	"time"

	"o.o/api/main/ordering/types"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateShipnowFulfillment(context.Context, *CreateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	CreateShipnowFulfillmentV2(context.Context, *CreateShipnowFulfillmentV2Args) (*ShipnowFulfillment, error)

	ConfirmShipnowFulfillment(context.Context, *ConfirmShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	CancelShipnowFulfillment(context.Context, *CancelShipnowFulfillmentArgs) (*meta.Empty, error)

	UpdateShipnowFulfillment(context.Context, *UpdateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

	UpdateShipnowFulfillmentCarrierInfo(context.Context, *UpdateShipnowFulfillmentCarrierInfoArgs) (*ShipnowFulfillment, error)

	UpdateShipnowFulfillmentState(context.Context, *UpdateShipnowFulfillmentStateArgs) (*ShipnowFulfillment, error)

	GetShipnowServices(context.Context, *GetShipnowServicesArgs) (*GetShipnowServicesResult, error)
}

type QueryService interface {
	GetShipnowFulfillment(context.Context, *GetShipnowFulfillmentQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
	GetShipnowFulfillments(context.Context, *GetShipnowFulfillmentsQueryArgs) (*GetShipnowFulfillmentsQueryResult, error)
	GetShipnowFulfillmentByShippingCode(context.Context, *GetShipnowFulfillmentByShippingCodeQueryArgs) (*GetShipnowFulfillmentQueryResult, error)
}

type CreateShipnowFulfillmentArgs struct {
	OrderIds            []dot.ID
	Carrier             carriertypes.Carrier
	ShopId              dot.ID
	ShippingServiceCode string
	ShippingServiceFee  int
	ShippingNote        string
	RequestPickupAt     time.Time
	PickupAddress       *types.Address
}

type CreateShipnowFulfillmentV2Args struct {
	DeliveryPoints      []*OrderShippingInfo
	Carrier             carriertypes.Carrier
	ShopID              dot.ID
	ShippingServiceCode string
	ShippingServiceFee  int
	ShippingNote        string
	RequestPickupAt     time.Time
	PickupAddress       *types.Address
}

type OrderShippingInfo struct {
	OrderID         dot.ID
	ShippingAddress *types.Address
	ShippingNote    string
	shippingtypes.WeightInfo
	shippingtypes.ValueInfo
	TryOn try_on.TryOnCode
}

type ConfirmShipnowFulfillmentArgs struct {
	Id     dot.ID
	ShopId dot.ID
}

type CancelShipnowFulfillmentArgs struct {
	Id           dot.ID
	ShopId       dot.ID
	CancelReason string
}

type UpdateShipnowFulfillmentArgs struct {
	Id                  dot.ID
	OrderIds            []dot.ID
	Carrier             carriertypes.Carrier
	ShopId              dot.ID
	ShippingServiceCode string
	ShippingServiceFee  int
	ShippingNote        string
	RequestPickupAt     time.Time
	PickupAddress       *types.Address
}

type UpdateShipnowFulfillmentCarrierInfoArgs struct {
	Id                         dot.ID
	ShippingCode               string
	ShippingState              shipnow_state.State
	TotalFee                   int
	FeeLines                   []*shippingtypes.FeeLine
	CarrierFeeLines            []*shippingtypes.FeeLine
	ShippingCreatedAt          time.Time
	EtopPaymentStatus          status4.Status
	ShippingStatus             status5.Status
	Status                     status5.Status
	CodEtopTransferedAt        time.Time
	ShippingPickingAt          time.Time
	ShippingDeliveringAt       time.Time
	ShippingDeliveredAt        time.Time
	ShippingCancelledAt        time.Time
	ShippingServiceName        string
	CancelReason               string
	ShippingSharedLink         string
	ShippingServiceDescription string
}

type UpdateShipnowFulfillmentStateArgs struct {
	Id             dot.ID
	SyncStatus     status4.Status
	Status         status5.Status
	ConfirmStatus  status3.Status
	ShippingStatus status5.Status
	SyncStates     *SyncStates
	ShippingState  shipnow_state.State
}

type GetShipnowServicesCommandResult struct {
	Services []*shipnowtypes.ShipnowService
}

type GetShipnowFulfillmentQueryArgs struct {
	Id     dot.ID
	ShopId dot.ID
}

type GetShipnowFulfillmentsQueryArgs struct {
	ShopIds []dot.ID
	Paging  *meta.Paging
	Filters []*meta.Filter
}

type GetShipnowFulfillmentQueryResult struct {
	ShipnowFulfillment *ShipnowFulfillment
}

type GetShipnowFulfillmentsQueryResult struct {
	ShipnowFulfillments []*ShipnowFulfillment
	Count               int
}

type GetShipnowFulfillmentByShippingCodeQueryArgs struct {
	ShippingCode string
}

type GetShipnowServicesArgs struct {
	ShopId         dot.ID
	OrderIds       []dot.ID
	PickupAddress  *types.Address
	DeliveryPoints []*DeliveryPoint
}

type GetShipnowServicesResult struct {
	Services []*shipnowtypes.ShipnowService
}
