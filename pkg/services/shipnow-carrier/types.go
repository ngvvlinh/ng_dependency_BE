package shipnow_carrier

import (
	"context"

	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
)

type ShipnowCarrier interface {
	CreateExternalShipnow(context.Context, *carrier.CreateExternalShipnowCommand) (_ error)
	CancelExternalShipnow(context.Context, *carrier.CancelExternalShipnowCommand) error
	GetShippingServices(context.Context, GetShippingServiceArgs) ([]*carrier.AvailableShippingService, error)
	ParseServiceCode(code string) (serviceName string, ok bool)
}

type GetShippingServiceArgs struct {
	AccountID        int64
	FromProvinceCode string
	FromDistrictCode string

	DeliveryPoints []*shipnowtypes.DeliveryPoint
}
