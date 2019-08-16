package shipping

import (
	"context"
	"time"

	"etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
)

type AggregateBus struct{ meta.Bus }

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

	Fee int32

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
	ID        int64
	OrderID   int64
	ShopID    int64
	PartnerID int64
	SelfURL   string

	Lines []*ItemLine

	TotalItems int32

	WeightInfo

	ValueInfo
}

type ItemLine struct {
	OrderID int64

	ProductName string
	ProductID   int64
	VariantID   int64
	IsOutside   bool
	ImageURL    string
	Attribute   []Attribute

	Quantity     int32
	ListPrice    int32
	RetailPrice  int32
	PaymentPrice int32
}

type Attribute struct {
	Name  string
	Value string
}

type WeightInfo struct {
	GrossWeight      int32
	ChargeableWeight int32
	Length           int32
	Width            int32
	Height           int32
}

type ValueInfo struct {
	BasketValue      int32
	CODAmount        int32
	IncludeInsurance bool
}

type ExternalShipmentData struct {
	State string

	ShippingFee int32

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
	OrderID int64

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
	FulfillmentID int64
}

type CancelFulfillmentArgs struct {
	FulfillmentID int64

	CancelReason string
}

//-- Queries --//

type GetFulfillmentByIDQueryArgs struct {
	FulfillmentID int64
}
