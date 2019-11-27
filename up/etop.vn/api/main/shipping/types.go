package shipping

import (
	"context"
	"time"

	"etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

// +gen:api

type AggregateBus struct{ capi.Bus }

type Aggregate interface {
	GetFulfillmentByID(context.Context, *GetFulfillmentByIDQueryArgs) (*Fulfillment, error)

	CreateFulfillment(context.Context, *CreateFulfillmentArgs) (*meta.Empty, error)
	ConfirmFulfillment(context.Context, *ConfirmFulfillmentArgs) (*meta.Empty, error)
	CancelFulfillment(context.Context, *CancelFulfillmentArgs) (*meta.Empty, error)
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

type Address struct {
	// @required
	FullName string

	Phone string

	Email string

	// optional
	Company string

	// @required
	Address1 string

	// optional
	Address2 string

	Location
}

type Location struct {
	DistrictCode string
	WardCode     string
}

type Fulfillment struct {
	ID        dot.ID
	OrderID   dot.ID
	ShopID    dot.ID
	PartnerID dot.ID
	SelfURL   string

	Lines []*ItemLine

	TotalItems int

	WeightInfo

	ValueInfo
}

type ItemLine struct {
	OrderID dot.ID

	ProductName string
	ProductID   dot.ID
	VariantID   dot.ID
	IsOutside   bool
	ImageURL    string
	Attribute   []Attribute

	Quantity     int
	ListPrice    int
	RetailPrice  int
	PaymentPrice int
}

type Attribute struct {
	Name  string
	Value string
}

type WeightInfo struct {
	GrossWeight      int
	ChargeableWeight int
	Length           int
	Width            int
	Height           int
}

type ValueInfo struct {
	BasketValue      int
	CODAmount        int
	IncludeInsurance bool
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

type CreateFulfillmentArgs struct {
	OrderID dot.ID

	PickupAddress *Address

	ShippingAddress *Address

	ReturnAddress *Address

	Carrier string

	ShippingServiceCode string

	ShippingServiceFee string

	WeightInfo

	ValueInfo

	TryOn types.TryOn

	ShippingNote string
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
