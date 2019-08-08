package ghn

import (
	"context"
	"strconv"
	"time"

	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"

	"etop.vn/api/main/location"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	ghnclient "etop.vn/backend/pkg/integration/ghn/client"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/common/bus"
)

var _ shipping_provider.ShippingProvider = &Carrier{}

type Carrier struct {
	clients  map[ClientType]*ghnclient.Client
	location location.QueryBus
}

func New(cfg Config, location location.QueryBus) *Carrier {
	clientDefault := ghnclient.New(cfg.Env, cfg.AccountDefault.AccountID, cfg.AccountDefault.Token)
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
			Type: model.TypeGHN,
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (p *Carrier) CreateFulfillment(
	ctx context.Context,
	order *ordermodel.Order,
	ffm *shipmodel.Fulfillment,
	args shipping_provider.GetShippingServicesArgs,
	service *model.AvailableShippingService,
) (ffmToUpdate *shipmodel.Fulfillment, _err error) {

	note := shipping_provider.GetShippingProviderNote(order, ffm)
	noteCode := order.GhnNoteCode
	if noteCode == "" {
		// harcode
		noteCode = "CHOXEMHANGKHONGTHU"
	}

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	toDistrict := toQuery.Result.District

	ghnCmd := &RequestCreateOrderCommand{
		ServiceID: service.ProviderServiceID,
		Request: &ghnclient.CreateOrderRequest{
			FromDistrictID:     int(fromQuery.Result.District.GhnId),
			ToDistrictID:       int(toQuery.Result.District.GhnId),
			Note:               note,
			ExternalCode:       strconv.FormatInt(ffm.ID, 10),
			ClientContactName:  ffm.AddressFrom.GetFullName(),
			ClientContactPhone: ffm.AddressFrom.Phone,
			ClientAddress:      ffm.AddressFrom.GetFullAddress(),
			CustomerName:       ffm.AddressTo.GetFullName(),
			CustomerPhone:      ffm.AddressTo.Phone,
			ShippingAddress:    ffm.AddressTo.GetFullAddress(),
			CoDAmount:          ffm.TotalCODAmount,
			NoteCode:           noteCode,
			Weight:             args.ChargeableWeight,
			Length:             args.Length,
			Width:              args.Width,
			Height:             args.Height,
			InsuranceFee:       args.GetInsuranceAmount(),
		},
	}

	if ffm.AddressReturn != nil {
		returnQuery := &location.GetLocationQuery{DistrictCode: ffm.AddressReturn.DistrictCode}
		if err := p.location.Dispatch(ctx, returnQuery); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ trả hàng không hợp lệ: %v", err)
		}
		returnDistrict := returnQuery.Result.District

		ghnCmd.Request.ReturnContactName = ffm.AddressReturn.GetFullName()
		ghnCmd.Request.ReturnContactPhone = ffm.AddressReturn.Phone
		ghnCmd.Request.ReturnAddress = ffm.AddressReturn.GetFullAddress()
		ghnCmd.Request.ReturnDistrictID = int(returnDistrict.GhnId)

		// ExternalReturnCode is required, we generate a random code here
		ghnCmd.Request.ExternalReturnCode = cm.IDToDec(cm.NewID())
	}

	if ghnErr := p.CreateOrder(ctx, ghnCmd); ghnErr != nil {
		return nil, ghnErr
	}

	r := ghnCmd.Result

	now := time.Now()
	expectedDeliveryAt := shipping.CalcDeliveryTime(model.TypeGHN, toDistrict, r.ExpectedDeliveryTime.ToTime())
	updateFfm := &shipmodel.Fulfillment{
		ID:                ffm.ID,
		ProviderServiceID: service.ProviderServiceID,
		Status:            model.S5SuperPos, // Now processing

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
		ShippingState:             model.StateCreated,
		SyncStatus:                model.S4Positive,
		SyncStates: &model.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedPickAt:     service.ExpectedPickAt,
		ExpectedDeliveryAt: expectedDeliveryAt,
	}
	// Fake shipping_fee_shop_lines, it will automates update later (when receive webhook)
	updateFfm.ProviderShippingFeeLines = []*model.ShippingFeeLine{
		{
			ShippingFeeType:      model.ShippingFeeTypeMain,
			Cost:                 int(r.TotalServiceFee),
			ExternalShippingCode: r.OrderCode.String(),
		},
	}
	updateFfm.ShippingFeeShopLines = model.GetShippingFeeShopLines(updateFfm.ProviderShippingFeeLines, updateFfm.EtopPriceRule, &updateFfm.EtopAdjustedShippingFeeMain)
	return updateFfm, nil
}

func (p *Carrier) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment, action model.FfmAction) error {
	code := ffm.ExternalShippingCode
	var ghnErr error
	providerServiceID := ffm.ProviderServiceID
	switch action {
	case model.FfmActionCancel:
		ghnCmd := &RequestCancelOrderCommand{
			ServiceID: providerServiceID,
			Request:   &ghnclient.OrderCodeRequest{OrderCode: code},
		}
		ghnErr = p.CancelOrder(ctx, ghnCmd)

	case model.FfmActionReturn:
		ghnCmd := &RequestReturnOrderCommand{
			ServiceID: providerServiceID,
			Request:   &ghnclient.OrderCodeRequest{OrderCode: code},
		}
		ghnErr = p.ReturnOrder(ctx, ghnCmd)

	default:
		panic("expected")
	}
	return ghnErr
}

func (p *Carrier) GetShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, toDistrict := fromQuery.Result.District, toQuery.Result.District
	if fromDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ gửi hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", fromDistrict.Name)
	}
	if toDistrict.GhnId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "GHN: Địa chỉ nhận hàng %v không được hỗ trợ bởi đơn vị vận chuyển!", toDistrict.Name)
	}

	cmd := &RequestFindAvailableServicesCommand{
		FromDistrict: fromDistrict,
		ToDistrict:   toDistrict,
		Request: &ghnclient.FindAvailableServicesRequest{
			Connection:     ghnclient.Connection{},
			Weight:         int(args.ChargeableWeight),
			Length:         int(args.Length),
			Width:          int(args.Width),
			Height:         int(args.Height),
			FromDistrictID: int(fromDistrict.GhnId),
			ToDistrictID:   int(toDistrict.GhnId),
			InsuranceFee:   args.GetInsuranceAmount(),
		},
	}
	err := p.FindAvailableServices(ctx, cmd)
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

	cmd := &RequestFindAvailableServicesCommand{
		FromDistrict: fromDistrict,
		ToDistrict:   toDistrict,
		Request: &ghnclient.FindAvailableServicesRequest{
			Connection:     ghnclient.Connection{},
			Weight:         int(args.ChargeableWeight),
			Length:         int(args.Length),
			Width:          int(args.Width),
			Height:         int(args.Height),
			FromDistrictID: int(fromDistrict.GhnId),
			ToDistrictID:   int(toDistrict.GhnId),
			InsuranceFee:   args.GetInsuranceAmount(),
		},
	}
	err := c.FindAvailableServices(ctx, cmd)
	if err != nil {
		return nil, err
	}
	providerServices := cmd.Result

	// get ETOP services
	etopServiceArgs := &etop_shipping_price.GetEtopShippingServicesArgs{
		ArbitraryID:  args.AccountID,
		Carrier:      model.TypeGHN,
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

func (p *Carrier) GetShippingService(ffm *shipmodel.Fulfillment, order *ordermodel.Order, weight int, valueInsurance int) (providerService *model.AvailableShippingService, etopService *model.AvailableShippingService, err error) {
	return nil, nil, cm.ErrTODO
}
