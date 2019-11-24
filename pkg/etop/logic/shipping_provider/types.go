package shipping_provider

import (
	"context"

	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

type ShippingProvider interface {
	CreateFulfillment(context.Context, *ordermodel.Order, *shipmodel.Fulfillment, GetShippingServicesArgs, *model.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error)
	CancelFulfillment(context.Context, *shipmodel.Fulfillment, model.FfmAction) error
	GetShippingServices(ctx context.Context, args GetShippingServicesArgs) ([]*model.AvailableShippingService, error)
	// combine with ETOP services
	GetAllShippingServices(ctx context.Context, args GetShippingServicesArgs) ([]*model.AvailableShippingService, error)

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
