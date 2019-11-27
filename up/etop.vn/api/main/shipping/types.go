package shipping

import (
	"context"
	"time"

	catalogtype "etop.vn/api/main/catalog/types"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/try_on"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

// +gen:api
// +gen:event:topic=event/shipping

type AggregateBus struct{ capi.Bus }

type Aggregate interface {
	CreateFulfillments(context.Context, *CreateFulfillmentsArgs) (fulfillmentID []dot.ID, _ error)
}

//-- Types --//

type ShippingService struct {
	Code string

	Name string

	Fee int

	Carrier string

	EstimatedPickupAt time.Time

	EstimatedDeliveryAt time.Time
}

type Fulfillment struct {
	ID        dot.ID
	OrderID   dot.ID
	ShopID    dot.ID
	PartnerID dot.ID
	SelfURL   string

	Lines []*ItemLine

	ShopConfirm   status3.Status
	ConfirmStatus status3.Status

	TotalItems    int
	TotalDiscount int
	TotalAmount   int

	types.WeightInfo

	types.ValueInfo
}

type ItemLine struct {
	OrderID dot.ID

	ProductName string
	ProductID   dot.ID
	VariantID   dot.ID
	IsOutside   bool
	ImageURL    string
	Attribute   []catalogtype.Attribute

	Quantity     int
	ListPrice    int
	RetailPrice  int
	PaymentPrice int
}

type Attribute struct {
	Name  string
	Value string
}

type ExternalShipmentData struct {
	State string

	ShippingFee int

	// ShippingData

	// ShippingLogs

	ShippingFeeLine []*ShippingFeeLine

	CreatedAt time.Time

	UpdatedAt time.Time

	PickingAt time.Time

	PickedAt time.Time

	HoldingAt time.Time

	DeliveringAt time.Time

	DeliveredAt time.Time

	ReturningAt time.Time

	ReturnedAt time.Time
}

type ShippingFeeLine struct {
	ShippingFeeType string

	Cost int

	ExternalServiceID string

	ExternalServiceName string

	ExternalServiceType string

	//
	ExternalShipmentID string
}

//-- Commands --//

type CreateFulfillmentsArgs struct {
	ShopID dot.ID

	OrderID dot.ID

	PickupAddress *ordertypes.Address

	ShippingAddress *ordertypes.Address

	ReturnAddress *ordertypes.Address

	ShippingType ordertypes.ShippingType

	ShippingServiceCode string

	ShippingServiceFee int

	ShippingServiceName string

	types.WeightInfo

	types.ValueInfo

	TryOn try_on.TryOnCode

	ShippingNote string

	ConnectionID dot.ID

	ShopCarrierID dot.ID
}

type ConfirmFulfillmentArgs struct {
	FulfillmentID dot.ID
}

type CancelFulfillmentArgs struct {
	FulfillmentID dot.ID

	CancelReason string
}

//-- Queries --//

type GetFulfillmentByIDQueryArgs struct {
	FulfillmentID dot.ID
}

type FulfillmentCreatingEvent struct {
	meta.EventMeta
	ShopID      dot.ID
	ShippingFee int
}
