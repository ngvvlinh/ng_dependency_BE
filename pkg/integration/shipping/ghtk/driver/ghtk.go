package driver

import (
	"context"
	"sort"
	"sync"
	"time"

	"o.o/api/main/location"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	carrierutil "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/randgenerator"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	ghtkupdate "o.o/backend/pkg/integration/shipping/ghtk/update"
	"o.o/capi/dot"
)

var _ carriertypes.ShipmentCarrier = &GHTKDriver{}

var (
	defaultDrivers = []string{
		"shipping/shipment/builtin/ghtk",
		"shipping/shipment/direct/ghtk",
	}
)

type GHTKDriver struct {
	client     *ghtkclient.Client
	locationQS location.QueryBus
}

func New(env string, cfg ghtkclient.GhtkAccount, locationQS location.QueryBus) *GHTKDriver {
	client := ghtkclient.New(env, cfg)
	return &GHTKDriver{
		client:     client,
		locationQS: locationQS,
	}
}

func (d *GHTKDriver) Ping(context.Context) error {
	if err := d.client.Ping(); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Can not init GHTK client")
	}
	return nil
}

func (d *GHTKDriver) GetAffiliateID() string {
	return d.client.GetAffiliateID()
}

func (d *GHTKDriver) CreateFulfillment(
	ctx context.Context,
	ffm *shipmodel.Fulfillment,
	args *carriertypes.GetShippingServicesArgs,
	service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	note := carrierutil.GetShippingProviderNote(ffm)

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()

	// prepare products for ghtk
	var products []*ghtkclient.ProductRequest
	for _, line := range ffm.Lines {
		products = append(products, &ghtkclient.ProductRequest{
			Name:     line.ProductName,
			Price:    line.ListPrice,
			Quantity: line.Quantity,
		})
	}

	transport, err := d.parseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)
	cmd := &ghtkclient.CreateOrderRequest{
		Products: products,
		Order: &ghtkclient.OrderRequest{
			ID:           ffm.ID.String(),
			PickName:     ffm.AddressFrom.GetFullName(),
			PickMoney:    ffm.TotalCODAmount,
			PickAddress:  cm.Coalesce(ffm.AddressFrom.Address1, ffm.AddressFrom.Address2),
			PickProvince: fromProvince.Name,
			PickDistrict: fromDistrict.Name,
			PickWard:     ffm.AddressFrom.Ward,
			PickTel:      ffm.AddressFrom.Phone,
			Name:         ffm.AddressTo.GetFullName(),
			Address:      cm.Coalesce(ffm.AddressTo.Address1, ffm.AddressTo.Address2),
			Province:     toProvince.Name,
			District:     toDistrict.Name,
			Ward:         ffm.AddressTo.Ward,
			Tel:          ffm.AddressTo.Phone,
			Note:         note,
			WeightOption: "gram",
			Value:        insuranceValue,
			TotalWeight:  float32(args.ChargeableWeight),
			Transport:    transport,
		},
	}
	if ffm.AddressReturn != nil {
		returnQuery := &location.GetLocationQuery{DistrictCode: ffm.AddressReturn.DistrictCode}
		if err := d.locationQS.Dispatch(ctx, returnQuery); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ trả hàng không hợp lệ: %v", err)
		}
		returnProvince, returnDistrict := returnQuery.Result.Province, returnQuery.Result.District

		cmd.Order.UseReturnAddress = 1
		cmd.Order.ReturnName = ffm.AddressReturn.GetFullName()
		cmd.Order.ReturnAddress = cm.Coalesce(ffm.AddressReturn.Address1, ffm.AddressReturn.Address2)
		cmd.Order.ReturnProvince = returnProvince.Name
		cmd.Order.ReturnDistrict = returnDistrict.Name
		cmd.Order.ReturnWard = ffm.AddressReturn.Ward
		cmd.Order.ReturnTel = ffm.AddressReturn.Phone
		returnEmail := ffm.AddressReturn.Email
		if returnEmail == "" {
			returnEmail = wl.X(ctx).CSEmail
		}
		// ReturnEmail can not empty
		cmd.Order.ReturnEmail = returnEmail
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

		ShippingCode:              ghtk.NormalizeGHTKCode(r.Order.Label.String()),
		ExternalShippingName:      service.Name,
		ExternalShippingID:        r.Order.TrackingID.String(),
		ExternalShippingCode:      r.Order.Label.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       int(r.Order.Fee),
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
	if expectedDeliveryAt, err := shipping.ParseDateTimeShipping(r.Order.EstimatedDeliverTime.String()); err == nil {
		updateFfm.ExpectedDeliveryAt = shipping.CalcDeliveryTime(shipping_provider.GHTK, toDistrict, *expectedDeliveryAt)
	}

	// prepare info to calc providerShippingFeeLines
	orderInfo := &ghtkclient.OrderInfo{
		LabelID:   r.Order.Label,
		ShipMoney: r.Order.Fee,
		Insurance: r.Order.InsuranceFee,
	}
	updateFfm.ProviderShippingFeeLines = ghtk.CalcAndConvertShippingFeeLines(orderInfo)
	updateFfm.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(updateFfm.ProviderShippingFeeLines, updateFfm.EtopPriceRule, dot.Int(updateFfm.EtopAdjustedShippingFeeMain))

	return updateFfm, nil
}

func (d *GHTKDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	code := ffm.ExternalShippingCode
	_, err := d.client.CancelOrder(ctx, code, "")
	return err
}

func (d *GHTKDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
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
		Request: &ghtkclient.CalcShippingFeeRequest{
			Weight:          args.ChargeableWeight,
			Value:           args.GetInsuranceAmount(maxValueFreeInsuranceFee),
			PickingProvince: fromProvince.Name,
			PickingDistrict: fromDistrict.Name,
			Province:        toProvince.Name,
			District:        toDistrict.Name,
		},
	}
	return d.CalcShippingFee(ctx, cmd)
}

func (d *GHTKDriver) CalcShippingFee(ctx context.Context, args *CalcShippingFeeArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	type Result struct {
		Transport ghtkclient.TransportType
		Result    *ghtkclient.CalcShippingFeeResponse
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex
	wg.Add(2)
	go func() {
		defer wg.Done()
		// clone the request to prevent race condition
		req := *args.Request
		req.Transport = ghtkclient.TransportRoad
		resp, err := d.client.CalcShippingFee(ctx, &req)
		if err != nil {
			return
		}
		m.Lock()
		results = append(results, Result{ghtkclient.TransportRoad, resp})
		m.Unlock()
	}()
	go func() {
		defer wg.Done()
		// trường hợp nội tỉnh: có gói nhanh
		// trường hợp nội vùng: bỏ qua gói nhanh
		if args.FromProvince.Code != args.ToProvince.Code && args.FromProvince.Region == args.ToProvince.Region {
			return
		}

		req := *args.Request
		req.Transport = ghtkclient.TransportFly
		resp, err := d.client.CalcShippingFee(ctx, &req)
		if err != nil {
			return
		}
		m.Lock()
		results = append(results, Result{ghtkclient.TransportFly, resp})
		m.Unlock()
	}()
	wg.Wait()
	if len(results) == 0 {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Giaohangtietkiem: không thể lấy thông tin gói cước dịch vụ")
	}
	// Sort result for stable service id generating. This must run before generating service id
	sort.Slice(results, func(i, j int) bool {
		return results[i].Transport < results[j].Transport
	})

	// Calc expectedPictAt & expectedDeliveryAt
	now := time.Now()
	expectedPickAt := shipping.CalcPickTime(shipping_provider.GHTK, now)
	generator := randgenerator.NewGenerator(args.ArbitraryID)
	var res []*shippingsharemodel.AvailableShippingService
	for _, result := range results {
		providerServiceID, err := d.GenerateServiceID(generator, result.Transport)
		if err != nil {
			continue
		}
		if !result.Result.Fee.Delivery {
			continue
		}

		expectedDeliveryDuration := ghtk.CalcDeliveryDuration(result.Transport, args.FromProvince, args.ToProvince)
		expectedDeliveryAt := expectedPickAt.Add(expectedDeliveryDuration)
		service := result.Result.Fee.ToShippingService(providerServiceID, result.Transport, expectedPickAt, expectedDeliveryAt)
		res = append(res, service)
	}
	res = shipping.CalcServicesTime(shipping_provider.GHTK, args.FromDistrict, args.ToDistrict, res)
	return res, nil
}

func (d *GHTKDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	switch {
	case code[5] == 'R': // road
		return etopmodel.ShippingServiceNameStandard, true
	case code[6] == 'F': // fly
		return etopmodel.ShippingServiceNameFaster, true
	}
	return "", false
}

func (d *GHTKDriver) GetMaxValueFreeInsuranceFee() int {
	return 3000000
}

func (d *GHTKDriver) SignIn(ctx context.Context, args *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	cmd := &ghtkclient.SignInRequest{
		Email:    args.Identifier,
		Password: args.Password,
	}
	resp, err := d.client.SignIn(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &carriertypes.AccountResponse{
		Token:  resp.Data.Token,
		UserID: resp.Data.Code,
	}, nil
}

func (d *GHTKDriver) SignUp(ctx context.Context, args *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	cmd := &ghtkclient.SignUpRequest{
		Name:         args.Name,
		FirstAddress: args.Address,
		Province:     args.Province,
		District:     args.District,
		Tel:          args.Phone,
		Email:        args.Email,
	}
	resp, err := d.client.SignUp(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &carriertypes.AccountResponse{
		Token:  resp.Data.Token,
		UserID: resp.Data.Code,
	}, nil
}

func (d *GHTKDriver) parseServiceID(code string) (transport ghtkclient.TransportType, err error) {
	res, err := d.ParseServiceID(code)
	if err != nil {
		return "", err
	}
	return ghtkclient.TransportType(res), nil
}

func (d *GHTKDriver) ParseServiceID(code string) (serviceID string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Missing service id")
		return
	}
	if len(code) != 8 {
		err = cm.Errorf(cm.InvalidArgument, nil, "GHTK invalid service id (code = %v)", code)
		return
	}

	if code[5] == 'R' {
		serviceID = string(ghtkclient.TransportRoad)
	}
	if code[6] == 'F' {
		if serviceID != "" {
			err = cm.Errorf(cm.InvalidArgument, nil, "GHTK invalid service id. Transport is invalid (code = %v)", code)
		}
		serviceID = string(ghtkclient.TransportFly)
	}
	if serviceID == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "GHTK invalid service id (code = %v)", code)
	}
	return
}

func (d *GHTKDriver) GenerateServiceID(generator *randgenerator.RandGenerator, transport ghtkclient.TransportType) (string, error) {
	if transport == "" {
		return "", cm.Errorf(cm.Internal, nil, "GHTK transport can not be empty").WithMeta("GHTK", "func GenerateServiceID")
	}

	code := generator.RandomAlphabet32(8)
	switch transport {
	case ghtkclient.TransportRoad:
		code[5] = 'R'
		code[6] = carrierutil.Blacklist(code[6], 'J', 'R', 'F')
	case ghtkclient.TransportFly:
		code[6] = 'F'
		code[5] = carrierutil.Blacklist(code[5], 'J', 'R', 'F')
	default:
		return "", cm.Errorf(cm.Internal, nil, "GHTK invalid transport")
	}

	// backward compatible
	// old id: the fourth character is the client code
	var clientCode byte = 'T' // default client code
	code[3] = clientCode
	return string(code), nil
}

func (d *GHTKDriver) RefreshFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, err error) {
	externalOrder, err := d.client.GetOrder(ctx, ffm.ShippingCode, "")
	if err != nil {
		return nil, err
	}

	ffmToUpdate, err = ghtkupdate.CalcRefreshFulfillmentInfo(ffm, &externalOrder.Order)
	return
}

func (d *GHTKDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *GHTKDriver) UpdateFulfillmentInfo(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}
