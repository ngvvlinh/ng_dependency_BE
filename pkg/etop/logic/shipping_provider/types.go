package shipping_provider

import (
	"context"

	"o.o/api/top/types/etc/shipping_provider"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type CarrierDriver interface {
	Code() shipping_provider.ShippingProvider

	CreateFulfillment(context.Context, *ordermodel.Order, *shipmodel.Fulfillment, GetShippingServicesArgs, *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error)
	CancelFulfillment(context.Context, *shipmodel.Fulfillment, model.FfmAction) error
	GetShippingServices(ctx context.Context, args GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error)
	// combine with ETOP services
	GetAllShippingServices(ctx context.Context, args GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error)

	// Return "chuẩn" or "nhanh"
	ParseServiceCode(code string) (serviceName string, ok bool)

	GetMaxValueFreeInsuranceFee() int
}

type GetShippingServicesArgs struct {
	AccountID        dot.ID
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

func (args *GetShippingServicesArgs) GetInsuranceAmount(maxValueFreeInsuranceFee int) int {
	if args.IncludeInsurance {
		return args.BasketValue
	}
	if args.BasketValue <= maxValueFreeInsuranceFee {
		return args.BasketValue
	}
	return maxValueFreeInsuranceFee
}
