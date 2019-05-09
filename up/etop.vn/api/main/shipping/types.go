package shipping

import (
	"context"
	"time"

	"etop.vn/api/meta"
)

type ProcessManagerBus struct {
	meta.Bus
}

type Aggregate interface {
	GetFulfillmentByID(ctx context.Context, query *GetFulfillmentByIDQuery) (*Fulfillment, error)

	CreateFulfillment(ctx context.Context, cmd *CreateFulfillmentCommand) error
	ConfirmFulfillment(ctx context.Context, cmd *ConfirmFulfillmentCommand) error
	CancelFulfillment(ctx context.Context, cmd *CancelFulfillmentCommand) error
}

type ProcessManager interface {
}

//-- Types --//

type TryOn int

const (
	TryOnNone TryOn = 1
	TryOnOpen TryOn = 2
	TryOnTry  TryOn = 3
)

func (code TryOn) String() string {
	switch code {
	case TryOnNone:
		return "none"
	case TryOnOpen:
		return "open"
	case TryOnTry:
		return "try"
	default:
		return "unknown"
	}
}

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
	ID         int64
	OrderID    int64
	ShopID     int64
	SupplierID int64
	PartnerID  int64
	SelfURL    string

	Lines []*ItemLine

	TotalItems int32

	WeightInfo

	ValueInfo
}

type ItemLine struct {
	OrderID    int64
	SupplierID int64

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

type CreateFulfillmentCommand struct {
	OrderID int64

	PickupAddress *Address

	ShippingAddress *Address

	ReturnAddress *Address

	Carrier string

	ShippingServiceCode string

	ShippingServiceFee string

	WeightInfo

	ValueInfo

	TryOn TryOn

	ShippingNote string

	Result *Fulfillment
}

type ConfirmFulfillmentCommand struct {
	FulfillmentID int64
}

type CancelFulfillmentCommand struct {
	FulfillmentID int64

	CancelReason string
}

//-- Queries --//

type GetFulfillmentByIDQuery struct {
	FulfillmentID int64

	Result *Fulfillment
}
