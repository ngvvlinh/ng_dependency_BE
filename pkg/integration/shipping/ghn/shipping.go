package ghn

import (
	"context"
	"strconv"
	"time"

	"etop.vn/api/main/location"
	"etop.vn/api/top/types/etc/ghn_note_code"
	typesshipping "etop.vn/api/top/types/etc/shipping"
	typeshippingprovider "etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
	ghnupdate "etop.vn/backend/pkg/integration/shipping/ghn/update"
	"etop.vn/capi/dot"
)

var _ shipping_provider.ShippingCarrier = &Carrier{}

type Carrier struct {
	clients  map[ClientType]*ghnclient.Client
	location location.QueryBus
}

func New(cfg Config, location location.QueryBus) *Carrier {
	accountCfg := ghnclient.GHNAccountCfg{
		ClientID: cfg.AccountDefault.AccountID,
		Token:    cfg.AccountDefault.Token,
	}
	clientDefault := ghnclient.New(cfg.Env, accountCfg)
	clients := map[ClientType]*ghnclient.Client{
		GHNCodeDefault: clientDefault,
	}

	return &Carrier{
		clients:  clients,
		location: location,
	}
}

func (c *Carrier) InitAllClients(ctx context.Context) error {
	for cName, c := range c.clients {
		if err := c.Ping(); err != nil {
			return cm.Error(cm.ExternalServiceError, "can not init client", err).
				WithMetap("client", cName)
		}
		cmd := &model.CreateShippingSource{
			Name: GetShippingSourceName(cName, c.ClientID()),
			Type: typeshippingprovider.GHN,
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (c *Carrier) CreateFulfillment(
	ctx context.Context,
	order *ordermodel.Order,
	ffm *shipmodel.Fulfillment,
	args shipping_provider.GetShippingServicesArgs,
	service *model.AvailableShippingService,
) (ffmToUpdate *shipmodel.Fulfillment, _err error) {

	note := shipping_provider.GetShippingProviderNote(order, ffm)
	noteCode := order.GhnNoteCode
	if noteCode == 0 {
		// harcode
		noteCode = ghn_note_code.CHOXEMHANGKHONGTHU
	}

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	toDistrict := toQuery.Result.District
	maxValueFreeInsuranceFee := c.GetMaxValueFreeInsuranceFee()

	ghnCmd := &RequestCreateOrderCommand{
		ServiceID: service.ProviderServiceID,
		Request: &ghnclient.CreateOrderRequest{
			FromDistrictID:     fromQuery.Result.District.GhnId,
			ToDistrictID:       toQuery.Result.District.GhnId,
			Note:               note,
			ExternalCode:       ffm.ID.String(),
			ClientContactName:  ffm.AddressFrom.GetFullName(),
			ClientContactPhone: ffm.AddressFrom.Phone,
			ClientAddress:      ffm.AddressFrom.GetFullAddress(),
			CustomerName:       ffm.AddressTo.GetFullName(),
			CustomerPhone:      ffm.AddressTo.Phone,
			ShippingAddress:    ffm.AddressTo.GetFullAddress(),
			CoDAmount:          ffm.TotalCODAmount,
			NoteCode:           noteCode.String(), // TODO: convert to try on code
			Weight:             args.ChargeableWeight,
			Length:             args.Length,
			Width:              args.Width,
			Height:             args.Height,
			InsuranceFee:       args.GetInsuranceAmount(maxValueFreeInsuranceFee),
		},
	}

	if ffm.AddressReturn != nil {
		returnQuery := &location.GetLocationQuery{DistrictCode: ffm.AddressReturn.DistrictCode}
		if err := c.location.Dispatch(ctx, returnQuery); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ trả hàng không hợp lệ: %v", err)
		}
		returnDistrict := returnQuery.Result.District

		ghnCmd.Request.ReturnContactName = ffm.AddressReturn.GetFullName()
		ghnCmd.Request.ReturnContactPhone = ffm.AddressReturn.Phone
		ghnCmd.Request.ReturnAddress = ffm.AddressReturn.GetFullAddress()
		ghnCmd.Request.ReturnDistrictID = returnDistrict.GhnId

		// ExternalReturnCode is required, we generate a random code here
		ghnCmd.Request.ExternalReturnCode = cm.IDToDec(cm.NewID())
	}

	if ghnErr := c.CreateOrder(ctx, ghnCmd); ghnErr != nil {
		return nil, ghnErr
	}

	r := ghnCmd.Result

	now := time.Now()
	expectedDeliveryAt := shipping.CalcDeliveryTime(typeshippingprovider.GHN, toDistrict, r.ExpectedDeliveryTime.ToTime())
	updateFfm := &shipmodel.Fulfillment{
		ID:                ffm.ID,
		ProviderServiceID: service.ProviderServiceID,
		Status:            status5.S, // Now processing

		ShippingFeeCustomer: order.ShopShippingFee,
		ShippingFeeShop:     order.ShopShipping.ExternalShippingFee,

		ShippingCode:              r.OrderCode.String(),
		ExternalShippingName:      service.Name,
		ExternalShippingID:        r.OrderID.String(),
		ExternalShippingCode:      r.OrderCode.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       int(r.TotalServiceFee),
		ShippingState:             typesshipping.Created,
		SyncStatus:                status4.P,
		SyncStates: &model.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedPickAt:     service.ExpectedPickAt,
		ExpectedDeliveryAt: expectedDeliveryAt,
	}
	// Get order GHN to update ProviderShippingFeeLine
	ghnGetOrderCmd := &RequestGetOrderCommand{
		ServiceID: service.ProviderServiceID,
		Request: &ghnclient.OrderCodeRequest{
			OrderCode: r.OrderCode.String(),
		},
	}
	if err := c.GetOrder(ctx, ghnGetOrderCmd); err == nil {
		updateFfm.ProviderShippingFeeLines = ghnclient.CalcAndConvertShippingFeeLines(ghnGetOrderCmd.Result.ShippingOrderCosts)
	}
	updateFfm.CreatedBy = order.CreatedBy
	updateFfm.ShippingFeeShopLines = model.GetShippingFeeShopLines(updateFfm.ProviderShippingFeeLines, updateFfm.EtopPriceRule, dot.Int(updateFfm.EtopAdjustedShippingFeeMain))
	return updateFfm, nil
}

func (c *Carrier) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment, action model.FfmAction) error {
	code := ffm.ExternalShippingCode
	var ghnErr error
	providerServiceID := ffm.ProviderServiceID
	switch action {
	case model.FfmActionCancel:
		ghnCmd := &RequestCancelOrderCommand{
			ServiceID: providerServiceID,
			Request:   &ghnclient.OrderCodeRequest{OrderCode: code},
		}
		ghnErr = c.CancelOrder(ctx, ghnCmd)

	case model.FfmActionReturn:
		ghnCmd := &RequestReturnOrderCommand{
			ServiceID: providerServiceID,
			Request:   &ghnclient.OrderCodeRequest{OrderCode: code},
		}
		ghnErr = c.ReturnOrder(ctx, ghnCmd)

	default:
		panic("expected")
	}
	return ghnErr
}

func (c *Carrier) GetShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, toDistrict := fromQuery.Result.District, toQuery.Result.District
	if fromDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ gửi hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", fromDistrict.Name)
	}
	if toDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ nhận hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", toDistrict.Name)
	}
	maxValueFreeInsuranceFee := c.GetMaxValueFreeInsuranceFee()

	cmd := &RequestFindAvailableServicesCommand{
		FromDistrict: fromDistrict,
		ToDistrict:   toDistrict,
		Request: &ghnclient.FindAvailableServicesRequest{
			Connection:     ghnclient.Connection{},
			Weight:         args.ChargeableWeight,
			Length:         args.Length,
			Width:          args.Width,
			Height:         args.Height,
			FromDistrictID: fromDistrict.GhnId,
			ToDistrictID:   toDistrict.GhnId,
			InsuranceFee:   args.GetInsuranceAmount(maxValueFreeInsuranceFee),
		},
	}
	err := c.FindAvailableServices(ctx, cmd)
	return cmd.Result, err
}

func (c *Carrier) GetAllShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province

	if fromDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ gửi hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", fromDistrict.Name)
	}
	if toDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ nhận hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", toDistrict.Name)
	}
	maxValueFreeInsuranceFee := c.GetMaxValueFreeInsuranceFee()

	cmd := &RequestFindAvailableServicesCommand{
		FromDistrict: fromDistrict,
		ToDistrict:   toDistrict,
		Request: &ghnclient.FindAvailableServicesRequest{
			Connection:     ghnclient.Connection{},
			Weight:         args.ChargeableWeight,
			Length:         args.Length,
			Width:          args.Width,
			Height:         args.Height,
			FromDistrictID: fromDistrict.GhnId,
			ToDistrictID:   toDistrict.GhnId,
			InsuranceFee:   args.GetInsuranceAmount(maxValueFreeInsuranceFee),
		},
	}
	err := c.FindAvailableServices(ctx, cmd)
	if err != nil {
		return nil, err
	}
	providerServices := cmd.Result

	// get Etop services
	etopServiceArgs := &etop_shipping_price.GetEtopShippingServicesArgs{
		ArbitraryID:  args.AccountID,
		Carrier:      typeshippingprovider.GHN,
		FromProvince: fromProvince,
		ToProvince:   toProvince,
		ToDistrict:   toDistrict,
		Weight:       args.ChargeableWeight,
	}
	etopServices := etop_shipping_price.GetEtopShippingServices(etopServiceArgs)
	etopServices, _ = etop_shipping_price.FillInfoEtopServices(providerServices, etopServices)

	allServices := append(providerServices, etopServices...)
	return allServices, nil
}

func (c *Carrier) GetShippingService(ffm *shipmodel.Fulfillment, order *ordermodel.Order, weight int, valueInsurance int) (providerService *model.AvailableShippingService, etopService *model.AvailableShippingService, err error) {
	return nil, nil, cm.ErrTODO
}

func (c *Carrier) CalcRefreshFulfillmentInfo(ctx context.Context, ffm *shipmodel.Fulfillment, orderGHN *ghnclient.Order) (*shipmodel.Fulfillment, error) {
	update, err := ghnupdate.CalcRefreshFulfillmentInfo(ffm, orderGHN)
	if err != nil {
		return nil, err
	}

	// Always update shipping address because we don't know whether it was changed
	addressQuery := &location.GetLocationQuery{
		DistrictCode:     strconv.Itoa(int(orderGHN.ToDistrictID)),
		LocationCodeType: location.LocCodeTypeGHN,
	}
	if err := c.location.Dispatch(ctx, addressQuery); err != nil {
		// ignore this error
		return update, nil
	}
	province := addressQuery.Result.Province
	district := addressQuery.Result.District
	addressTo := update.AddressTo
	addressTo.ProvinceCode = province.Code
	addressTo.Province = province.Name
	addressTo.DistrictCode = district.Code
	addressTo.District = district.Name
	addressTo.Address1 = orderGHN.ShippingAddress.String()
	// Reset ward: GHN does not require ward
	// so when address change, we don't have information about the ward and actually we don't need it.
	addressTo.WardCode = ""
	addressTo.Ward = ""
	return update, nil
}

func (c *Carrier) GetMaxValueFreeInsuranceFee() int {
	return 1000000
}
