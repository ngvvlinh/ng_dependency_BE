package shipnow

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/ordering/types"
	carriertypes "etop.vn/api/main/shipnow/carrier/types"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
)

type Aggregate interface {
	CreateShipnowFulfillment(context.Context, *CreateShipnowFulfillmentArgs) (*ShipnowFulfillment, error)

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
	OrderIds            []int64
	Carrier             carriertypes.Carrier
	ShopId              int64
	ShippingServiceCode string
	ShippingServiceFee  int32
	ShippingNote        string
	RequestPickupAt     time.Time
	PickupAddress       *types.Address
}

type ConfirmShipnowFulfillmentArgs struct {
	Id     int64
	ShopId int64
}

type CancelShipnowFulfillmentArgs struct {
	Id           int64
	ShopId       int64
	CancelReason string
}

type UpdateShipnowFulfillmentArgs struct {
	Id                  int64
	OrderIds            []int64
	Carrier             carriertypes.Carrier
	ShopId              int64
	ShippingServiceCode string
	ShippingServiceFee  int32
	ShippingNote        string
	RequestPickupAt     time.Time
	PickupAddress       *types.Address
}

type UpdateShipnowFulfillmentCarrierInfoArgs struct {
	Id                         int64
	ShippingCode               string
	ShippingState              shipnowtypes.State
	TotalFee                   int32
	FeeLines                   []*shippingtypes.FeeLine
	CarrierFeeLines            []*shippingtypes.FeeLine
	ShippingCreatedAt          time.Time
	EtopPaymentStatus          etop.Status4
	ShippingStatus             etop.Status5
	Status                     etop.Status5
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
	Id             int64
	SyncStatus     etop.Status4
	Status         etop.Status5
	ConfirmStatus  etop.Status3
	ShippingStatus etop.Status5
	SyncStates     *SyncStates
	ShippingState  shipnowtypes.State
}

type GetShipnowServicesCommandResult struct {
	Services []*shipnowtypes.ShipnowService
}

type GetShipnowFulfillmentQueryArgs struct {
	Id     int64
	ShopId int64
}

type GetShipnowFulfillmentsQueryArgs struct {
	ShopIds []int64
	Paging  *meta.Paging
	Filters []*meta.Filter
}

type GetShipnowFulfillmentQueryResult struct {
	ShipnowFulfillment *ShipnowFulfillment
}

type GetShipnowFulfillmentsQueryResult struct {
	ShipnowFulfillments []*ShipnowFulfillment
	Count               int32
}

type GetShipnowFulfillmentByShippingCodeQueryArgs struct {
	ShippingCode string
}

type GetShipnowServicesArgs struct {
	ShopId         int64
	OrderIds       []int64
	PickupAddress  *types.Address
	DeliveryPoints []*DeliveryPoint
}

type GetShipnowServicesResult struct {
	Services []*shipnowtypes.ShipnowService
}
