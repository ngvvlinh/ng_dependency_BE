package types

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/top/types/etc/connection_type"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/capi/dot"
)

type Config struct {
	Endpoints []ConfigEndpoint
}

type ConfigEndpoint struct {
	Provider connection_type.ConnectionProvider
	Endpoint string
}

type ConfigEndpoints []ConfigEndpoint

func (es ConfigEndpoints) GetByCarrier(carrier connection_type.ConnectionProvider) (string, bool) {
	for _, e := range es {
		if e.Provider == carrier {
			return e.Endpoint, true
		}
	}
	return "", false
}

type Driver interface {
	GetShipmentDriver(
		env string, locationQS location.QueryBus,
		connection *connectioning.Connection,
		shopConnection *connectioning.ShopConnection,
		endpoints ConfigEndpoints,
	) (ShipmentCarrier, error)

	GetAffiliateShipmentDriver(
		env string, locationQS location.QueryBus,
		connection *connectioning.Connection,
		endpoints ConfigEndpoints,
	) (ShipmentCarrier, error)
}

type ShipmentCarrier interface {
	Ping(context.Context) error

	GetAffiliateID() string

	CreateFulfillment(context.Context, *shipmodel.Fulfillment, *GetShippingServicesArgs, *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error)

	RefreshFulfillment(context.Context, *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error)

	UpdateFulfillmentInfo(context.Context, *shipmodel.Fulfillment) error

	UpdateFulfillmentCOD(context.Context, *shipmodel.Fulfillment) error

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

	ChargeableWeight int
	Length           int
	Width            int
	Height           int
	IncludeInsurance bool
	InsuranceValue   dot.NullInt
	BasketValue      int
	CODAmount        int
	Coupon           string
}

func (args *GetShippingServicesArgs) GetInsuranceAmount(maxValueFreeInsurance int) int {
	if args.IncludeInsurance {
		return args.InsuranceValue.Apply(args.BasketValue)
	}
	if args.BasketValue <= maxValueFreeInsurance {
		return args.BasketValue
	}
	return maxValueFreeInsurance
}

type SignInArgs struct {
	Identifier string // email or phone
	Password   string
	// Field is added to resolve a situation when signIn(new mechanism) with  GHN
	OTP string
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
	Token         string
	UserID        string
	ShopID        string
	IsRequiredOTP bool
}
