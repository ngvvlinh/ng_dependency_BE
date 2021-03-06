package driver

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"o.o/api/main/location"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	"o.o/backend/pkg/etc/typeutil"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	v1 "o.o/backend/pkg/integration/shipping/ghn/update/v1"
	"o.o/capi/dot"
)

var _ carriertypes.ShipmentCarrier = &GHNDriver{}

var (
	defaultDrivers = []string{
		"shipping/shipment/builtin/ghn",
		"shipping/shipment/direct/ghn",
	}
)

type GHNDriver struct {
	client          *ghnclient.Client
	locationQS      location.QueryBus
	webhookEndpoint string
}

func New(env string, cfg ghnclient.GHNAccountCfg, locationQS location.QueryBus, webhookEndpoint string) *GHNDriver {
	client := ghnclient.New(env, cfg)
	return &GHNDriver{
		client:          client,
		locationQS:      locationQS,
		webhookEndpoint: webhookEndpoint,
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

func (d *GHNDriver) GenerateToken(ctx context.Context) (*carriertypes.GenerateTokenResponse, error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *GHNDriver) CreateFulfillment(
	ctx context.Context,
	ffm *shipmodel.Fulfillment,
	args *carriertypes.GetShippingServicesArgs, service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	// note := carrierutil.GetShippingProviderNote(ffm)
	note := ffm.ShippingNote
	noteCode := ffm.TryOn
	if noteCode == 0 {
		// harcode
		noteCode = try_on.None
	}
	ghnNoteCode := typeutil.GHNNoteCodeFromTryOn(noteCode)

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict := fromQuery.Result.District
	toDistrict := toQuery.Result.District
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()

	serviceID, err := d.parseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)
	cmd := &ghnclient.CreateOrderRequest{
		FromDistrictID:     fromDistrict.GhnId,
		ToDistrictID:       toDistrict.GhnId,
		Note:               note,
		ExternalCode:       ffm.ID.String(),
		ClientContactName:  ffm.AddressFrom.GetFullName(),
		ClientContactPhone: ffm.AddressFrom.Phone,
		ClientAddress:      ffm.AddressFrom.GetFullAddress(),
		CustomerName:       ffm.AddressTo.GetFullName(),
		CustomerPhone:      ffm.AddressTo.Phone,
		ShippingAddress:    ffm.AddressTo.GetFullAddress(),
		CoDAmount:          ffm.TotalCODAmount,
		NoteCode:           ghnNoteCode.String(),
		Weight:             args.ChargeableWeight,
		Length:             cm.CoalesceInt(args.Length, 10),
		Width:              cm.CoalesceInt(args.Width, 10),
		Height:             cm.CoalesceInt(args.Height, 10),
		InsuranceFee:       insuranceValue,
		ServiceID:          serviceID,
	}

	if ffm.AddressReturn != nil {
		returnQuery := &location.GetLocationQuery{DistrictCode: ffm.AddressReturn.DistrictCode}
		if err := d.locationQS.Dispatch(ctx, returnQuery); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "?????a ch??? tr??? h??ng kh??ng h???p l???: %v", err)
		}
		returnDistrict := returnQuery.Result.District

		cmd.ReturnContactName = ffm.AddressReturn.GetFullName()
		cmd.ReturnContactPhone = ffm.AddressReturn.Phone
		cmd.ReturnAddress = ffm.AddressReturn.GetFullAddress()
		cmd.ReturnDistrictID = int(returnDistrict.GhnId)

		// ExternalReturnCode is required, we generate a random code here
		cmd.ExternalReturnCode = cm.IDToDec(cm.NewID())
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
		ExternalShippingID:        r.OrderID.String(),
		ExternalShippingCode:      r.OrderCode.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       int(r.TotalServiceFee),
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
	// Calc expected delivery at
	// add some rules
	expectedDeliveryAt := shipping.CalcDeliveryTime(shipping_provider.GHN, toDistrict, r.ExpectedDeliveryTime.ToTime())
	updateFfm.ExpectedDeliveryAt = expectedDeliveryAt

	// Get order GHN to update ProviderShippingFeeLine
	ghnGetOrderCmd := &ghnclient.OrderCodeRequest{
		OrderCode: r.OrderCode.String(),
	}
	if res, err := d.client.GetOrderInfo(ctx, ghnGetOrderCmd); err == nil {
		updateFfm.ProviderShippingFeeLines = ghnclient.CalcAndConvertShippingFeeLines(res.ShippingOrderCosts)
	}

	updateFfm.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(updateFfm.ProviderShippingFeeLines, updateFfm.EtopPriceRule, dot.Int(updateFfm.EtopAdjustedShippingFeeMain))
	return updateFfm, nil
}

func (d *GHNDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	code := ffm.ExternalShippingCode
	cmd := &ghnclient.OrderCodeRequest{OrderCode: code}
	return d.client.CancelOrder(ctx, cmd)
}

func (d *GHNDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	if fromDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: ?????a ch??? g???i h??ng %v kh??ng ???????c h??? tr??? b???i ????n v??? v???n chuy???n!", fromDistrict.Name)
	}
	if toDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: ?????a ch??? nh???n h??ng %v kh??ng ???????c h??? tr??? b???i ????n v??? v???n chuy???n!", toDistrict.Name)
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
		},
	}

	return d.CalcShippingFee(ctx, cmd)
}

func (d *GHNDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (d *GHNDriver) GetMaxValueFreeInsuranceFee() int {
	return 1000000
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
	// Th???ng nh???t service id cho t???t c??? NVC, sau ???? parse t????ng ???ng

	serviceID = code[2:]
	return serviceID, nil
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

func (d *GHNDriver) CalcShippingFee(ctx context.Context, args *CalcShippingFeeArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	req := args.Request
	resp, err := d.client.FindAvailableServices(ctx, args.Request)
	if err != nil {
		return nil, err
	}

	res := resp.AvailableServices
	if len(res) > 0 && req.InsuranceFee > 0 {
		// Include insurance:
		//
		// "FindAvailableServices" does not include insurance in calculation,
		// therefore we must call "CalculateFee" to get insurance fee, then add it to all previous calculation.
		calcFeeCmd := &ghnclient.CalculateFeeRequest{
			Weight:         req.Weight,
			Length:         req.Length,
			Width:          req.Width,
			Height:         req.Height,
			FromDistrictID: req.FromDistrictID,
			ToDistrictID:   req.ToDistrictID,
			InsuranceFee:   req.InsuranceFee,
			ServiceID:      int(res[0].ServiceID),
		}
		calcFeeResult, err := d.client.CalculateFee(ctx, calcFeeCmd)
		if err != nil {
			return nil, cm.Errorf(cm.ExternalServiceError, err, "L???i t??? GHN: Kh??ng th??? t??nh ???????c ph?? giao h??ng (%v)", err)
		}
		insuranceFee := ghnclient.GetInsuranceFee(calcFeeResult.OrderCosts)
		for _, shippingservice := range res {
			shippingservice.ServiceFeeMain = shippingservice.ServiceFee
			shippingservice.ServiceFee += ghnclient.Int(insuranceFee)
		}
	}
	// Sort result for stable service id generating. This must run before generating service id
	sort.Slice(res, func(i, j int) bool {
		return res[i].ServiceID < res[j].ServiceID
	})

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
	case code[1] == 'C': // Chu???n
		return etopmodel.ShippingServiceNameStandard, true
	case code[1] == 'N': // Nhanh
		return etopmodel.ShippingServiceNameFaster, true
	}
	return "", false
}

func (d *GHNDriver) SignIn(ctx context.Context, args *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	if d.webhookEndpoint == "" {
		return nil, cm.Errorf(cm.Internal, nil, "GHN ????ng nh???p c???n cung c???p webhook endpoint.")
	}
	cmd := &ghnclient.SignInRequest{
		Email:    args.Identifier,
		Password: args.Password,
	}
	resp, err := d.client.SignIn(ctx, cmd)
	if err != nil {
		return nil, err
	}
	userID := strconv.Itoa(int(resp.ClientID))

	// register webhook
	args2 := &RegisterWebhookForClientArgs{
		TokenClients: []string{resp.Token},
		URLCallback:  d.webhookEndpoint,
	}
	if err := d.RegisterWebhook(ctx, args2); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "L???i ????ng k?? webhook: %v", err)
	}

	return &carriertypes.AccountResponse{
		Token:  resp.Token,
		UserID: userID,
	}, nil
}

func (d *GHNDriver) SignUp(ctx context.Context, args *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	if d.webhookEndpoint == "" {
		return nil, cm.Errorf(cm.Internal, nil, "GHN ????ng k?? c???n cung c???p webhook endpoint")
	}
	cmd := &ghnclient.SignUpRequest{
		Email:        args.Email,
		Password:     args.Password,
		ContactPhone: args.Phone,
		ContactName:  args.Name,
	}
	resp, err := d.client.SignUp(ctx, cmd)
	if err != nil {
		return nil, err
	}
	userID := strconv.Itoa(int(resp.ClientID))

	// register webhook
	args2 := &RegisterWebhookForClientArgs{
		TokenClients: []string{resp.Token},
		URLCallback:  d.webhookEndpoint,
	}
	if err := d.RegisterWebhook(ctx, args2); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "L???i ????ng k?? webhook: %v", err)
	}

	return &carriertypes.AccountResponse{
		Token:  resp.Token,
		UserID: userID,
	}, nil
}

func (d *GHNDriver) RegisterWebhook(ctx context.Context, args *RegisterWebhookForClientArgs) error {
	cmd := &ghnclient.RegisterWebhookForClientRequest{
		TokenClient: args.TokenClients,
		URLCallback: args.URLCallback,
	}
	return d.client.RegisterWebhookForClient(ctx, cmd)
}

func (d *GHNDriver) RefreshFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	cmd := &ghnclient.OrderCodeRequest{
		OrderCode: ffm.ExternalShippingCode,
	}
	externalOrder, err := d.client.GetOrderInfo(ctx, cmd)
	if err != nil {
		return nil, err
	}

	ffmToUpdate, err = v1.CalcRefreshFulfillmentInfo(ffm, externalOrder)
	return
}

func (d *GHNDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *GHNDriver) UpdateFulfillmentInfo(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}
