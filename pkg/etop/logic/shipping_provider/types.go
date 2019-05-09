package shipping_provider

import (
	"context"

	"etop.vn/backend/pkg/etop/model"
)

type ShippingProvider interface {
	CreateFulfillment(context.Context, *model.Order, *model.Fulfillment, GetShippingServicesArgs, *model.AvailableShippingService) (ffmToUpdate *model.Fulfillment, _ error)
	CancelFulfillment(context.Context, *model.Fulfillment, model.FfmAction) error
	GetShippingServices(ctx context.Context, args GetShippingServicesArgs) ([]*model.AvailableShippingService, error)
	// combine with ETOP services
	GetAllShippingServices(ctx context.Context, args GetShippingServicesArgs) ([]*model.AvailableShippingService, error)

	// Return "chuẩn" or "nhanh"
	ParseServiceCode(code string) (serviceName string, ok bool)
}

type GetShippingServicesArgs struct {
	AccountID        int64
	FromDistrictCode string
	ToDistrictCode   string

	ChargeableWeight int
	Length           int
	Width            int
	Height           int
	IncludeInsurance bool
	BasketValue      int
	CODAmount        int
}

func (args *GetShippingServicesArgs) GetInsuranceAmount() int {
	if args.IncludeInsurance {
		return args.BasketValue
	}
	return 0
}
