package types

import (
	"context"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
	"o.o/api/top/types/etc/connection_type"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/capi/dot"
)

const CarrierNote = "KHÔNG TỰ Ý HOÀN HÀNG. Gọi shop nếu giao 1 phần/thất bại."

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
		identityQS identity.QueryBus,
		connection *connectioning.Connection,
		shopConnection *connectioning.ShopConnection,
		shippingcodeQS shippingcode.QueryBus,
	) (ShipmentCarrier, error)

	GetAffiliateShipmentDriver(
		env string, locationQS location.QueryBus,
		identityQS identity.QueryBus,
		connection *connectioning.Connection,
		shippingcodeQS shippingcode.QueryBus,
	) (ShipmentCarrier, error)
}

type ShipmentCarrier interface {
	Ping(context.Context) error

	GetAffiliateID() string

	GenerateToken(context.Context) (*GenerateTokenResponse, error)

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

	ChargeableWeight int // gram
	Length           int // cm
	Width            int // cm
	Height           int // cm
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

type GenerateTokenResponse struct {
	AccessToken string
	ExpiresAt   time.Time
	TokenType   string
	ExpiresIn   int
}

type ShipmentSync interface {
	Init(context.Context) error

	Start(context.Context) error

	Stop(context.Context) error
}
