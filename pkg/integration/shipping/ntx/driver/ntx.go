package driver

import (
	"context"
	"fmt"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
	"o.o/api/top/types/etc/shipping_provider"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	"o.o/backend/pkg/integration/shipping"
	ntxclient "o.o/backend/pkg/integration/shipping/ntx/client"
	"time"
)

var _ carriertypes.ShipmentCarrier = &NTXDriver{}

type NTXDriver struct {
	client         *ntxclient.Client
	locationQS     location.QueryBus
	identityQS     identity.QueryBus
	shippingcodeQS shippingcode.QueryBus
}

func New(
	env string,
	cfg ntxclient.Config,
	locationQS location.QueryBus,
	identityQS identity.QueryBus,
	shippingcodeQS shippingcode.QueryBus,
) *NTXDriver {
	client := ntxclient.New(env, cfg)
	return &NTXDriver{
		client:         client,
		locationQS:     locationQS,
		identityQS:     identityQS,
		shippingcodeQS: shippingcodeQS,
	}
}

func (d *NTXDriver) Ping(ctx context.Context) error {
	return nil
}

func (d *NTXDriver) GetAffiliateID() string {
	return ""
}

func (d *NTXDriver) GenerateToken(ctx context.Context) (*carriertypes.GenerateTokenResponse, error) {
	panic("implement me")
}

func (d *NTXDriver) CreateFulfillment(ctx context.Context, fulfillment *shipmodel.Fulfillment, args *carriertypes.GetShippingServicesArgs, service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	panic("implement me")
}

func (d *NTXDriver) RefreshFulfillment(ctx context.Context, fulfillment *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	panic("implement me")
}

func (d *NTXDriver) UpdateFulfillmentInfo(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	panic("implement me")
}

func (d *NTXDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	panic("implement me")
}

func (d *NTXDriver) CancelFulfillment(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	panic("implement me")
}

func (d *NTXDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{
		DistrictCode: args.FromDistrictCode,
		WardCode:     args.FromWardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: args.ToDistrictCode,
		WardCode:     args.ToWardCode,
	}

	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}

	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	cmd := &ntxclient.CalcShippingFeeRequest{
		PartnerID:     d.client.PartnerID,
		CodAmount:     args.CODAmount,
		CargoValue:    args.CODAmount,
		Weight:        args.ChargeableWeight,
		PaymentMethod: 1,
		SProvinceID:   fromProvince.NTXId,
		SDistrictID:   fromDistrict.NTXId,
		RProvinceID:   toProvince.NTXId,
		RDistrictID:   toDistrict.NTXId,
		PackageNo:     1,
		UtmSource:     ntxclient.UTMSource,
	}

	r, err := d.client.CalcShippingFee(ctx, cmd)
	if err != nil {
		return nil, err
	}

	var result []*shippingsharemodel.AvailableShippingService
	arbitraryID := args.AccountID.Int64() + args.ArbitraryID.Int64()
	generator := randgenerator.NewGenerator(arbitraryID)
	for _, service := range r.Data {
		providerServiceID, err := GenerateServiceID(generator, service.ServiceID)
		if err != nil {
			return nil, err
		}

		deliveryAt, err := time.ParseInLocation("02/01/2006 15:04", service.LeadTime, time.Local)
		if err != nil {
			return nil, err
		}

		result = append(result, &shippingsharemodel.AvailableShippingService{
			Name:               service.ServiceName,
			Provider:           shipping_provider.NTX,
			ProviderServiceID:  providerServiceID,
			ShippingFeeMain:    service.MainFee,
			ServiceFee:         service.MainFee,
			ExpectedDeliveryAt: deliveryAt,
		})
	}

	result = shipping.CalcServicesTime(shipping_provider.NTX, fromDistrict, toDistrict, result)
	return result, nil
}

func (d *NTXDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (d *NTXDriver) ParseServiceID(code string) (serviceID string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Missing service id")
		return
	}
	if len(code) <= 3 {
		err = cm.Errorf(cm.InvalidArgument, nil, "NTX invalid service id (code = %v)", code)
		return
	}

	switch code[1] {
	case '1': // Nhanh
		serviceID = "NH"
	case '2': // Chuẩn
		serviceID = "CH"
	default:
		return "", cm.Errorf(cm.Internal, nil, fmt.Sprintf("NTX invalid service id (code = %v)", code))
	}
	return serviceID, nil
}

func (d *NTXDriver) GetMaxValueFreeInsuranceFee() int {
	return 0
}

func (d *NTXDriver) SignIn(ctx context.Context, args *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.Unimplemented, nil, "Không hỗ trợ đăng nhập tài khoản NTX")
}

func (d *NTXDriver) SignUp(ctx context.Context, args *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.Unimplemented, nil, "Không hỗ trợ đăng ký tài khoản NTX")
}

func GenerateServiceID(generator *randgenerator.RandGenerator, ID int) (string, error) {
	code := generator.RandomAlphabet32(8)
	switch ID {
	case 1:
		code[1] = '1'
	case 2:
		code[1] = '2'
	default:
		return "", cm.Errorf(cm.Internal, nil, "NTX invalid service name")
	}

	return string(code), nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) < 6 {
		return "", false
	}
	switch {
	case code[1] == '1': // Nhanh
		return "Nhanh", true
	case code[1] == '2': // Chuẩn
		return "Chuẩn", true
	}
	return "", false
}
