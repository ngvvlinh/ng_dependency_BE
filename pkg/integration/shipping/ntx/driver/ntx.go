package driver

import (
	"context"
	"fmt"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	"o.o/backend/pkg/integration/shipping"
	ntxclient "o.o/backend/pkg/integration/shipping/ntx/client"
	"o.o/capi/dot"
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
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NTXDriver) CreateFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment, args *carriertypes.GetShippingServicesArgs, service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	fromQuery := &location.GetLocationQuery{
		DistrictCode: ffm.AddressFrom.DistrictCode,
		WardCode:     ffm.AddressFrom.WardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: ffm.AddressTo.DistrictCode,
		WardCode:     ffm.AddressTo.WardCode,
	}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}

	fromWard, fromDistrict, fromProvince := fromQuery.Result.Ward, fromQuery.Result.District, fromQuery.Result.Province
	toWard, toDistrict, toProvince := toQuery.Result.Ward, toQuery.Result.District, toQuery.Result.Province
	if toWard.VtpostId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "NTX không thể giao hàng tới địa chỉ này (%v, %v, %v)", toWard.Name, toDistrict.Name, toProvince.Name)
	}

	orderService, err := d.ParseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	cmd := &ntxclient.CreateOrderRequest{
		PartnerID:      d.client.PartnerID,
		SName:          ffm.AddressFrom.GetFullName(),
		SPhone:         ffm.AddressFrom.Phone,
		SAddress:       cm.Coalesce(ffm.AddressFrom.Address1, ffm.AddressFrom.Address2),
		SProvinceID:    fromProvince.NTXId,
		SDistrictID:    fromDistrict.NTXId,
		SWardID:        fromWard.NTXId,
		RName:          ffm.AddressTo.GetFullName(),
		RPhone:         ffm.AddressTo.Phone,
		RAddress:       cm.Coalesce(ffm.AddressTo.Address1, ffm.AddressTo.Address2),
		RProvinceID:    toProvince.NTXId,
		RDistrictID:    toDistrict.NTXId,
		RWardID:        toWard.NTXId,
		CodAmount:      ffm.TotalCODAmount,
		PaymentMethod:  d.client.PaymentMethod,
		Weight:         float64(args.ChargeableWeight) / 1000,
		CargoContentID: 8,
		CargoContent:   "Khác",
		Note:           "",
		UtmSource:      ntxclient.UTMSource,
		PackageNo:      1,
	}

	generateShippingCodeQuery := &shippingcode.GenerateShippingCodeQuery{}
	if err := d.shippingcodeQS.Dispatch(ctx, generateShippingCodeQuery); err != nil {
		return nil, err
	}
	cmd.RefCode = generateShippingCodeQuery.Result
	if orderService == "NH" {
		cmd.ServiceID = 1
	}

	if orderService == "CH" {
		cmd.ServiceID = 2
	}

	r, err := d.client.CreateOrder(ctx, cmd)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()
	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)
	updateFfm := &shipmodel.Fulfillment{
		ID:                ffm.ID,
		ProviderServiceID: service.ProviderServiceID,
		Status:            status5.S,

		ShippingFeeShop: ffm.ShippingServiceFee,

		ShippingCode:              r.Order.BillCode,
		ExternalShippingName:      service.Name,
		ExternalShippingID:        r.Order.BillCode,
		ExternalShippingCode:      r.Order.BillCode,
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       r.Order.TotalFee,
		ShippingState:             shippingstate.Created,
		SyncStatus:                status4.P,
		SyncStates: &shippingsharemodel.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedDeliveryAt: service.ExpectedDeliveryAt,
		InsuranceValue:     dot.Int(insuranceValue),
	}

	return updateFfm, nil
}

func (d *NTXDriver) RefreshFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NTXDriver) UpdateFulfillmentInfo(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NTXDriver) UpdateFulfillmentCOD(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NTXDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	cmd := &ntxclient.CancelOrderRequest{
		ListDocode: []string{ffm.ExternalShippingCode},
	}

	_, err := d.client.CancelOrder(ctx, cmd)
	return err
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
		Weight:        float64(args.ChargeableWeight) / 1000,
		PaymentMethod: d.client.PaymentMethod,
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

func (d *NTXDriver) SignIn(context.Context, *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.Unimplemented, nil, "Không hỗ trợ đăng nhập tài khoản NTX")
}

func (d *NTXDriver) SignUp(context.Context, *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
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
