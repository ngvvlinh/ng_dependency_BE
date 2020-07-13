package driverv2

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/ghn_payment_type"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	carrierutil "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
)

var _ carriertypes.ShipmentCarrier = &GHNDriver{}

type GHNDriver struct {
	client     *ghnclient.Client
	locationQS location.QueryBus
}

func New(
	env string, cfg ghnclient.GHNAccountCfg,
	locationQS location.QueryBus,
) *GHNDriver {
	client := ghnclient.New(env, cfg)
	return &GHNDriver{
		client:     client,
		locationQS: locationQS,
	}
}

func (d *GHNDriver) Ping(ctx context.Context) error {
	if err := d.client.Ping(); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Can not init GHN client")
	}
	return nil
}

func (d *GHNDriver) GetAffiliateID() string {
	return d.client.GetAffiliateID()
}

func (d *GHNDriver) CreateFulfillment(
	ctx context.Context, ffm *shipmodel.Fulfillment,
	args *carriertypes.GetShippingServicesArgs, service *shippingsharemodel.AvailableShippingService,
) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	if args.FromWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Yêu cầu nhập phường/xã của địa chỉ lấy hàng")
	}
	if args.ToWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Yêu cầu nhập phường/xã của địa chỉ nhận hàng")
	}
	note := carrierutil.GetShippingProviderNote(ffm)
	noteCode := ffm.TryOn
	if noteCode == 0 {
		// harcode
		noteCode = try_on.None
	}
	ghnNoteCode := etopmodel.GHNNoteCodeFromTryOn(noteCode)

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
	fromDistrict, fromWard := fromQuery.Result.District, fromQuery.Result.Ward
	toDistrict, toWard := toQuery.Result.District, toQuery.Result.Ward
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()

	serviceID, err := d.parseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)
	cmd := &ghnclient.CreateOrderRequest{
		FromName:        ffm.AddressFrom.GetFullName(),
		FromPhone:       ffm.AddressFrom.Phone,
		FromAddress:     ffm.AddressFrom.GetFullAddress(),
		FromWardCode:    fromWard.GhnCode,
		FromDistrictID:  fromDistrict.GhnId,
		ToName:          ffm.AddressTo.GetFullName(),
		ToPhone:         ffm.AddressTo.Phone,
		ToAddress:       ffm.AddressTo.GetFullAddress(),
		ToWardCode:      toWard.GhnCode,
		ToDistrictID:    toDistrict.GhnId,
		ClientOrderCode: ffm.ID.String(),
		CODAmount:       ffm.TotalCODAmount,
		Weight:          args.ChargeableWeight,
		Length:          cm.CoalesceInt(args.Length, 10),
		Width:           cm.CoalesceInt(args.Width, 10),
		Height:          cm.CoalesceInt(args.Height, 10),
		InsuranceValue:  insuranceValue,
		Coupon:          args.Coupon,
		ServiceID:       serviceID,
		// người bán trả tiền ship (hardcode)
		PaymentTypeID: ghn_payment_type.BUYER.Enum(),
		Note:          note,
		RequiredNote:  ghnNoteCode.String(),
	}

	if ffm.AddressReturn != nil {
		returnQuery := &location.GetLocationQuery{
			DistrictCode: ffm.AddressReturn.DistrictCode,
			WardCode:     ffm.AddressReturn.WardCode,
		}
		if err := d.locationQS.Dispatch(ctx, returnQuery); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ trả hàng không hợp lệ: %v", err)
		}
		returnDistrict, returnWard := returnQuery.Result.District, returnQuery.Result.Ward

		cmd.ReturnPhone = ffm.AddressReturn.Phone
		cmd.ReturnAddress = ffm.AddressReturn.GetFullAddress()
		cmd.ReturnDistrictID = returnDistrict.GhnId
		cmd.ReturnWardCode = returnWard.GhnCode
	}

	r, err := d.client.CreateOrder(ctx, cmd)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	updateFfm := &shipmodel.Fulfillment{
		ID:                ffm.ID,
		ProviderServiceID: service.ProviderServiceID,
		Status:            status5.S, // Now processing

		ShippingFeeShop: ffm.ShippingServiceFee,

		ShippingCode:              r.OrderCode.String(),
		ExternalShippingName:      service.Name,
		ExternalShippingCode:      r.OrderCode.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       int(r.TotalFee),
		ShippingState:             shippingstate.Created,
		SyncStatus:                status4.P,
		SyncStates: &shippingsharemodel.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedPickAt:     service.ExpectedPickAt,
		ExpectedDeliveryAt: service.ExpectedDeliveryAt,
		InsuranceValue:     insuranceValue,
	}
	// Calc expected delivery at
	// add some rules
	expectedDeliveryAt := shipping.CalcDeliveryTime(shipping_provider.GHN, toDistrict, r.ExpectedDeliveryTime.ToTime())
	updateFfm.ExpectedDeliveryAt = expectedDeliveryAt

	updateFfm.ProviderShippingFeeLines = r.Fee.ToFeeLines()
	return updateFfm, nil
}

func (d *GHNDriver) RefreshFulfillment(ctx context.Context, fulfillment *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *GHNDriver) UpdateFulfillmentInfo(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	addressFrom := ffm.AddressFrom
	addressTo := ffm.AddressTo
	if addressFrom.WardCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Yêu cầu nhập phường/xã của địa chỉ lấy hàng")
	}
	if addressTo.WardCode == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Yêu cầu nhập phường/xã của địa chỉ nhận hàng")
	}
	note := carrierutil.GetShippingProviderNote(ffm)
	noteCode := ffm.TryOn
	if noteCode == 0 {
		// harcode
		noteCode = try_on.None
	}
	ghnNoteCode := etopmodel.GHNNoteCodeFromTryOn(noteCode)

	fromQuery := &location.GetLocationQuery{
		DistrictCode: addressFrom.DistrictCode,
		WardCode:     addressFrom.WardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: addressTo.DistrictCode,
		WardCode:     addressTo.WardCode,
	}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return err
	}
	fromDistrict, fromWard := fromQuery.Result.District, fromQuery.Result.Ward
	toDistrict, toWard := toQuery.Result.District, toQuery.Result.Ward

	cmd := &ghnclient.UpdateOrderRequest{
		OrderCode:      ffm.ShippingCode,
		FromName:       ffm.AddressFrom.GetFullName(),
		FromPhone:      ffm.AddressFrom.GetPhone(),
		FromAddress:    ffm.AddressFrom.GetFullAddress(),
		FromWardCode:   fromWard.GhnCode,
		FromDistrictID: fromDistrict.GhnId,
		ToName:         ffm.AddressTo.GetFullName(),
		ToPhone:        ffm.AddressTo.GetPhone(),
		ToAddress:      ffm.AddressTo.GetFullAddress(),
		ToWardCode:     toWard.GhnCode,
		ToDistrictID:   toDistrict.GhnId,
		Weight:         ffm.GrossWeight,
		InsuranceValue: &ffm.InsuranceValue,
		Note:           note,
		RequiredNote:   ghnNoteCode.String(),
	}
	if err := d.client.UpdateOrder(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (d *GHNDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	code := fulfillment.ExternalShippingCode
	cmd := &ghnclient.UpdateOrderCODRequest{
		OrderCode: code,
		CODAmount: fulfillment.TotalCODAmount,
	}
	return d.client.UpdateOrderCOD(ctx, cmd)
}

func (d *GHNDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	code := ffm.ExternalShippingCode
	cmd := &ghnclient.CancelOrderRequest{OrderCodes: []string{code}}
	return d.client.CancelOrder(ctx, cmd)
}

func (d *GHNDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	if args.FromWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ gửi hàng - phường/xã không được để trống!")
	}
	if args.ToWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ nhận hàng - phường/xã không được để trống!")
	}
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
	fromProvince, fromDistrict, fromWard := fromQuery.Result.Province, fromQuery.Result.District, fromQuery.Result.Ward
	toProvince, toDistrict, toWard := toQuery.Result.Province, toQuery.Result.District, toQuery.Result.Ward
	if fromDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ gửi hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", fromDistrict.Name)
	}
	if args.FromWardCode != "" && fromWard.GhnCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ gửi hàng %v, %v không được hỗ trợ bởi đơn vị vận chuyển!", fromDistrict.Name, fromWard.Name)
	}
	if toDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ nhận hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", toDistrict.Name)
	}
	if args.ToWardCode != "" && toWard.GhnCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ nhận hàng %v, %v không được hỗ trợ bởi đơn vị vận chuyển!", toDistrict.Name, toWard.Name)
	}
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()

	arbitraryID := args.AccountID.Int64() + args.ArbitraryID.Int64()
	cmd := &CalcShippingFeeArgs{
		ArbitraryID:  arbitraryID,
		FromProvince: fromProvince,
		FromDistrict: fromDistrict,
		ToProvince:   toProvince,
		ToDistrict:   toDistrict,
		Request: &ghnclient.FindAvailableServicesRequest{
			Weight:         args.ChargeableWeight,
			Length:         args.Length,
			Width:          args.Width,
			Height:         args.Height,
			FromDistrictID: fromDistrict.GhnId,
			ToDistrictID:   toDistrict.GhnId,
			InsuranceFee:   args.GetInsuranceAmount(maxValueFreeInsuranceFee),
			Coupon:         args.Coupon,
		},
	}
	if fromWard != nil {
		cmd.Request.FromWardCode = fromWard.GhnCode
	}
	if toWard != nil {
		cmd.Request.ToWardCode = toWard.GhnCode
	}

	carrierServices, err := d.CalcShippingFee(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return carrierServices, nil
}

func (d *GHNDriver) CalcShippingFee(ctx context.Context, args *CalcShippingFeeArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	resp, err := d.client.FindAvailableServices(ctx, args.Request)
	if err != nil {
		return nil, err
	}

	res := resp.AvailableServices

	var result []*shippingsharemodel.AvailableShippingService
	generator := randgenerator.NewGenerator(args.ArbitraryID)
	for _, s := range res {
		providerServiceID, err := GenerateServiceID(generator, s.Name.String(), s.ServiceID.String())
		if err != nil {
			continue
		}
		result = append(result, s.ToShippingService(providerServiceID))
	}
	result = shipping.CalcServicesTime(shipping_provider.GHN, args.FromDistrict, args.ToDistrict, result)

	return result, nil
}

func (d *GHNDriver) parseServiceID(code string) (serviceID int, err error) {
	res, err := d.ParseServiceID(code)
	if err != nil {
		return 0, err
	}
	serviceID, err = strconv.Atoi(res)
	if err != nil {
		err = cm.Errorf(cm.InvalidArgument, nil, "Invalid Service ID: %v", code)
	}
	return
}

func GenerateServiceID(generator *randgenerator.RandGenerator, serviceName string, serviceID string) (string, error) {
	if serviceName == "" {
		return "", cm.Errorf(cm.Internal, nil, "Service Name can not be empty").WithMeta("GHN", "func GenerateServiceID")
	}
	if serviceID == "" {
		return "", cm.Errorf(cm.Internal, nil, "ServiceID can not be empty").WithMeta("GHN", "func GenerateServiceID")
	}

	// backward compatible
	// old id: the first character is the client code
	clientCode := ghn.GHNCodeDefault
	shortCode := strings.ToUpper(string(serviceName[0]))
	return string(clientCode) + shortCode + serviceID, nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) < 6 {
		return "", false
	}
	switch {
	case code[1] == 'C': // Chuẩn
		return etopmodel.ShippingServiceNameStandard, true
	case code[1] == 'N': // Nhanh
		return etopmodel.ShippingServiceNameFaster, true
	}
	return "", false
}

func (d *GHNDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (d *GHNDriver) ParseServiceID(code string) (serviceID string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Missing service id")
		return
	}
	if len(code) <= 3 {
		err = cm.Errorf(cm.InvalidArgument, nil, "GHN invalid service id (code = %v)", code)
		return
	}

	// old service id format: "DC123123"
	// Thống nhất service id cho tất cả NVC, sau đó parse tương ứng

	serviceID = code[2:]
	return serviceID, nil
}

func (d *GHNDriver) GetMaxValueFreeInsuranceFee() int {
	return 3000000
}

func (d *GHNDriver) SignIn(
	ctx context.Context, args *carriertypes.SignInArgs,
) (*carriertypes.AccountResponse, error) {
	if args.OTP == "" {
		sendOTPShopAffiliateRequest := &ghnclient.SendOTPShopAffiliateRequest{
			Phone: args.Identifier,
		}
		if _, err := d.client.SendOTPShopToAffiliateAccount(ctx, sendOTPShopAffiliateRequest); err != nil {
			return nil, err
		}
		return &carriertypes.AccountResponse{
			IsRequiredOTP: true,
		}, nil
	}

	createShopAffiliateRequest := &ghnclient.CreateShopAffiliateRequest{
		Phone: args.Identifier,
		OTP:   args.OTP,
	}
	resp, err := d.client.CreateShopByAffiliateAccount(ctx, createShopAffiliateRequest)
	if err != nil {
		return nil, err
	}

	clientID := fmt.Sprintf("%d", resp.ClientID)
	shopID := fmt.Sprintf("%d", resp.ShopID)

	return &carriertypes.AccountResponse{
		UserID: clientID,
		ShopID: shopID,
		Token:  d.client.GetToken(),
	}, nil
}

func (d *GHNDriver) SignUp(
	ctx context.Context, args *carriertypes.SignUpArgs,
) (*carriertypes.AccountResponse, error) {
	panic("implement me")
}
