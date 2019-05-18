package carrier

import (
	"context"
	"time"
)

type Manager interface {
	CreateExternalShipping(ctx context.Context, ffm *CreateExternalShipnowCommand) error
	CancelExternalShipping(ctx context.Context, ffm *CancelExternalShipnowCommand) error
}

type CreateExternalShipnowCommand struct {
	// TODO
}

type CancelExternalShipnowCommand struct {
	ShipnowFulfillmentID int64
	ExternalShipnowID    int64
}

type ExternalShipnow struct {
	// TODO
}

type AvailableShippingService struct {
	Name string
	// ServiceFee: Tổng phí giao hàng (đã bao gồm phí chính + các phụ phí khác)
	ServiceFee int
	// ShippingFeeMain: Phí chính giao hàng
	ShippingFeeMain  int
	Carrier          string
	CarrierServiceID string

	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
}
