package types

import (
	"context"

	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/capi/dot"
)

type ShipmentCarrier interface {
	Ping(context.Context) error

	GetAffiliateID() string

	CreateFulfillment(context.Context, *shipmodel.Fulfillment, *GetShippingServicesArgs, *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error)

	UpdateFulfillment(context.Context, *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error)

	CancelFulfillment(context.Context, *shipmodel.Fulfillment) error

	GetShippingServices(ctx context.Context, args *GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error)

	// Return "chuẩn" or "nhanh"
	GetServiceName(code string) (serviceName string, ok bool)

	ParseServiceID(code string) (serviceID string, err error)

	GetMaxValueFreeInsuranceFee() int

	SignIn(context.Context, *SignInArgs) (*AccountResponse, error)

	SignUp(context.Context, *SignUpArgs) (*AccountResponse, error)
}

type GetShippingServicesArgs struct {
	// ArbitraryID: this is privided as a seed, for stable randomization
	ArbitraryID      dot.ID
	AccountID        dot.ID
	FromDistrictCode string
	FromWardCode     string
	ToDistrictCode   string
	ToWardCode       string

	ChargeableWeight       int
	Length                 int
	Width                  int
	Height                 int
	IncludeInsurance       bool
	BasketValue            int
	CODAmount              int
	IncludeTopshipServices bool
}

func (args *GetShippingServicesArgs) GetInsuranceAmount(maxValueFreeInsurance int) int {
	if args.IncludeInsurance {
		return args.BasketValue
	}
	if args.BasketValue <= maxValueFreeInsurance {
		return args.BasketValue
	}
	return maxValueFreeInsurance
}

type SignInArgs struct {
	Email    string
	Password string
}

type SignUpArgs struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Province string
	District string
	Address  string
}

type AccountResponse struct {
	Token  string
	UserID string
}
