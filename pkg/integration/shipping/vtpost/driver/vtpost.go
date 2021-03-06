package driver

import (
	"context"
	"strconv"
	"strings"
	"time"

	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	carrierutil "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/vtpost"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
	"o.o/capi/dot"
)

var _ carriertypes.ShipmentCarrier = &VTPostDriver{}

var defaultDrivers = []string{
	"shipping/shipment/builtin/vtpost",
	"shipping/shipment/direct/vtpost",
}

type VTPostDriver struct {
	client         *vtpostclient.ClientImpl
	locationQS     location.QueryBus
	shippingcodeQS shippingcode.QueryBus
}

func New(env string, token string, locationQS location.QueryBus, shippingcodeQS shippingcode.QueryBus) *VTPostDriver {
	client := vtpostclient.NewClientWithToken(env, token)
	return &VTPostDriver{
		client:         client,
		locationQS:     locationQS,
		shippingcodeQS: shippingcodeQS,
	}
}

func (d *VTPostDriver) Ping(context.Context) error {
	if err := d.client.Ping(); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Can not init VTPost client")
	}
	return nil
}

func (d *VTPostDriver) GetAffiliateID() string {
	// vtpost does not support affiliate
	return ""
}

func (d *VTPostDriver) GenerateToken(ctx context.Context) (*carrierutil.GenerateTokenResponse, error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d VTPostDriver) CreateFulfillment(
	ctx context.Context,
	ffm *shipmodel.Fulfillment,
	args *carriertypes.GetShippingServicesArgs,
	service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	if ffm.AddressFrom.WardCode == "" || ffm.AddressTo.WardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "VTPost y??u c???u th??ng tin ph?????ng x?? h???p l??? ????? giao h??ng")
	}

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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "VTPost kh??ng th??? giao h??ng t???i ?????a ch??? n??y (%v, %v, %v)", toWard.Name, toDistrict.Name, toProvince.Name)
	}
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()

	// prepare products for vtpost
	var products []*vtpostclient.Product
	var productName string

	if len(ffm.Lines) > 0 {
		for _, line := range ffm.Lines {
			if productName != "" {
				productName += " + "
			}
			productName += line.ProductName
			products = append(products, &vtpostclient.Product{
				ProductName:     line.ProductName,
				ProductPrice:    line.ListPrice,
				ProductQuantity: line.Quantity,
			})
		}
	} else {
		if ffm.LinesContent != "" {
			productName = ffm.LinesContent
			products = []*vtpostclient.Product{
				{
					ProductName:     ffm.LinesContent,
					ProductPrice:    ffm.BasketValue,
					ProductQuantity: ffm.TotalItems,
				},
			}
		}
	}

	orderService, err := d.parseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	deliveryDate := time.Now()
	deliveryDate.Add(30 * time.Minute)
	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)
	cmd := &vtpostclient.CreateOrderRequest{
		OrderNumber: "", // will be filled later,
		// hard code: 30 mins from now
		DeliveryDate:       deliveryDate.Format("02/01/2006 15:04:05"),
		SenderFullname:     ffm.AddressFrom.GetFullName(),
		SenderAddress:      cm.Coalesce(ffm.AddressFrom.Address1, ffm.AddressFrom.Address2),
		SenderPhone:        ffm.AddressFrom.Phone,
		SenderEmail:        ffm.AddressFrom.Email,
		SenderWard:         fromWard.VtpostId,
		SenderDistrict:     fromDistrict.VtpostId,
		SenderProvince:     fromProvince.VtpostId,
		ReceiverFullname:   ffm.AddressTo.GetFullName(),
		ReceiverAddress:    cm.Coalesce(ffm.AddressTo.Address1, ffm.AddressTo.Address2),
		ReceiverPhone:      ffm.AddressTo.Phone,
		ReceiverEmail:      ffm.AddressTo.Email,
		ReceiverWard:       toWard.VtpostId,
		ReceiverDistrict:   toDistrict.VtpostId,
		ReceiverProvince:   toProvince.VtpostId,
		ProductPrice:       insuranceValue,
		ProductWeight:      args.ChargeableWeight,
		OrderNote:          ffm.ShippingNote,
		MoneyCollection:    ffm.TotalCODAmount,
		MoneyTotalFee:      service.ServiceFee,
		ListItem:           products,
		ProductName:        productName,
		ProductDescription: productName,
		OrderService:       orderService,
	}

	generateShippingCodeQuery := &shippingcode.GenerateShippingCodeQuery{}
	if err := d.shippingcodeQS.Dispatch(ctx, generateShippingCodeQuery); err != nil {
		return nil, err
	}
	cmd.OrderNumber = generateShippingCodeQuery.Result

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

		ShippingCode:              r.Data.OrderNumber,
		ExternalShippingName:      service.Name,
		ExternalShippingID:        r.Data.OrderNumber,
		ExternalShippingCode:      r.Data.OrderNumber,
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       r.Data.MoneyTotal,
		ShippingState:             shippingstate.Created,
		SyncStatus:                status4.P,
		SyncStates: &shippingsharemodel.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedPickAt:     service.ExpectedPickAt,
		ExpectedDeliveryAt: service.ExpectedDeliveryAt,
		InsuranceValue:     dot.Int(insuranceValue),
	}

	// recalculate shipping fee
	shippingFees := &vtpostclient.ShippingFeeData{
		MoneyTotal:         r.Data.MoneyTotal,
		MoneyTotalFee:      r.Data.MoneyTotalFee,
		MoneyFee:           r.Data.MoneyFee,
		MoneyCollectionFee: r.Data.MoneyCollectionFee,
		MoneyOtherFee:      r.Data.MoneyOtherFee,
		MoneyVAT:           r.Data.MoneyFeeVAT,
		KpiHt:              r.Data.KpiHt,
	}
	if lines, err := shippingFees.CalcAndConvertShippingFeeLines(); err == nil {
		updateFfm.ProviderShippingFeeLines = lines
		updateFfm.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(lines, false, dot.NullInt{})
	}

	return updateFfm, nil
}

func (d VTPostDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	cmd := &vtpostclient.CancelOrderRequest{
		OrderNumber: ffm.ExternalShippingCode,
	}
	_, err := d.client.CancelOrder(ctx, cmd)
	return err
}

func (d VTPostDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()

	arbitraryID := args.AccountID.Int64() + args.ArbitraryID.Int64()
	cmd := &CalcShippingFeeArgs{
		ArbitraryID:  arbitraryID,
		FromProvince: fromProvince,
		FromDistrict: fromDistrict,
		ToProvince:   toProvince,
		ToDistrict:   toDistrict,
		Request: &vtpostclient.CalcShippingFeeAllServicesRequest{
			SenderProvince:   fromProvince.VtpostId,
			SenderDistrict:   fromDistrict.VtpostId,
			ReceiverProvince: toProvince.VtpostId,
			ReceiverDistrict: toDistrict.VtpostId,
			ProductWeight:    args.ChargeableWeight,
			ProductPrice:     args.GetInsuranceAmount(maxValueFreeInsuranceFee),
			MoneyCollection:  args.CODAmount,
		},
	}
	return d.CalcShippingFee(ctx, cmd)
}

func (d *VTPostDriver) CalcShippingFee(ctx context.Context, args *CalcShippingFeeArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	req := args.Request
	resp, err := d.client.CalcShippingFeeAllServices(ctx, req)
	if err != nil {
		return nil, err
	}
	var result []*shippingsharemodel.AvailableShippingService
	generator := randgenerator.NewGenerator(args.ArbitraryID)
	now := time.Now()
	expectedPickAt := shipping.CalcPickTime(shipping_provider.VTPost, now)
	for _, s := range resp {
		serviceCode := vtpostclient.VTPostOrderServiceCode(s.MaDVChinh)
		providerServiceID, err := d.GenerateServiceID(generator, serviceCode)
		if err != nil {
			continue
		}

		// ignore this service
		ignoreServices := []string{
			string(vtpostclient.OrderServiceCodeV60),
		}
		if cm.StringsContain(ignoreServices, string(serviceCode)) {
			continue
		}

		// recall get price to get exactly shipping fee for each service
		cmd := &vtpostclient.CalcShippingFeeRequest{
			SenderProvince:   req.SenderProvince,
			SenderDistrict:   req.SenderDistrict,
			ReceiverProvince: req.ReceiverProvince,
			ReceiverDistrict: req.ReceiverDistrict,
			OrderService:     serviceCode,
			ProductWeight:    req.ProductWeight,
			ProductPrice:     req.ProductPrice,
			MoneyCollection:  req.MoneyCollection,
		}
		r, err := d.client.CalcShippingFee(ctx, cmd)
		if err != nil {
			continue
		}
		// T??nh c?????c ch??nh (main fee)
		// D??ng ????? ??p d???ng b???ng gi?? ri??ng v??o c?????c ch??nh
		feeLines, err := r.Data.CalcAndConvertShippingFeeLines()
		feeMain := 0
		if err == nil {
			feeMain = shippingsharemodel.GetMainFee(feeLines)
		}
		s.GiaCuoc = r.Data.MoneyTotal

		thoigian := s.ThoiGian // has format: "12 gi???"
		thoigian = strings.Replace(thoigian, " gi???", "", -1)
		var expectedDeliveryDuration time.Duration
		hours, err := strconv.Atoi(thoigian)
		if err != nil {
			expectedDeliveryDuration = vtpost.CalcDeliveryDuration(serviceCode, args.FromProvince, args.ToProvince, args.FromDistrict, args.ToDistrict)
		} else {
			expectedDeliveryDuration = time.Duration(hours) * time.Hour
		}
		expectedDeliveryAt := expectedPickAt.Add(expectedDeliveryDuration)
		service := s.ToAvailableShippingService(providerServiceID, expectedPickAt, expectedDeliveryAt, feeMain)
		result = append(result, service)
	}
	result = shipping.CalcServicesTime(shipping_provider.VTPost, args.FromDistrict, args.ToDistrict, result)
	return result, nil
}

func (d VTPostDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	switch {
	case code[3] == 'S':
		return etopmodel.ShippingServiceNameStandard, true
	case code[4] == 'F':
		return etopmodel.ShippingServiceNameFaster, true
	}
	return "", false
}

func (d VTPostDriver) GetMaxValueFreeInsuranceFee() int {
	// Follow the policy of provider
	return 0
}

func (d VTPostDriver) SignIn(context.Context, *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.Unimplemented, nil, "Kh??ng h??? tr??? ????ng nh???p t??i kho???n VTPost")
}

func (d VTPostDriver) SignUp(context.Context, *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.Unimplemented, nil, "Kh??ng h??? tr??? ????ng k?? t??i kho???n VTPost")
}

func getLast3Character(code vtpostclient.VTPostOrderServiceCode) string {
	return string(code[len(code)-3:])
}

func (d *VTPostDriver) parseServiceID(code string) (orderService vtpostclient.VTPostOrderServiceCode, err error) {
	res, err := d.ParseServiceID(code)
	if err != nil {
		return "", err
	}
	return vtpostclient.VTPostOrderServiceCode(res), nil
}

func (d *VTPostDriver) ParseServiceID(code string) (orderService string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "missing service id")
		return
	}
	if len(code) != 8 {
		err = cm.Errorf(cm.InvalidArgument, nil, "VTPost invalid service id (code = %v)", code)
		return
	}

	res := code[len(code)-3:]

	// map EtopServiceCode to VTPostServiceCode
	_orderService, ok := vtpostclient.MapVTPostServiceCodes[res]
	if !ok {
		err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
	}
	orderService = _orderService.String()
	return
}

func (d *VTPostDriver) GenerateServiceID(generator *randgenerator.RandGenerator, serviceCode vtpostclient.VTPostOrderServiceCode) (string, error) {
	code := generator.RandomAlphabet32(5)

	switch serviceCode.Name() {
	case etopmodel.ShippingServiceNameStandard:
		code[3] = 'S'
		code[4] = carrierutil.Blacklist(code[4], 'J', 'S', 'F')
	case etopmodel.ShippingServiceNameFaster:
		code[4] = 'F'
		code[3] = carrierutil.Blacklist(code[3], 'J', 'S', 'F')
	default:
		return "", cm.Errorf(cm.Internal, nil, "VTPost invalid service code")
	}

	// map VTPostOrderServiceCode to EtopOrderServiceCode
	_serviceCode, ok := vtpostclient.MapEtopServiceCodes[serviceCode]
	if !ok {
		return "", cm.Errorf(cm.Internal, nil, "VTPost invalid service code")
	}

	// backward compatible
	// old id: the second character is the client code
	var clientCode byte = 'D'
	code[1] = clientCode
	return string(code) + _serviceCode, nil
}

func (d *VTPostDriver) RefreshFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	return nil, cm.Errorf(cm.InvalidArgument, nil, "VTPost does not support API get order info")
}

func (d *VTPostDriver) UpdateFulfillmentInfo(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *VTPostDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}
