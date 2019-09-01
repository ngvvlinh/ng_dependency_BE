package ghtk

import (
	"context"
	"strconv"
	"time"

	"etop.vn/api/main/location"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	ghtkclient "etop.vn/backend/pkg/integration/shipping/ghtk/client"
)

var _ shipping_provider.ShippingProvider = &Carrier{}

type Carrier struct {
	clients  map[byte]*ghtkclient.Client
	location location.QueryBus
}

func New(cfg Config, locationBus location.QueryBus) *Carrier {
	clientDefault := ghtkclient.New(cfg.Env, cfg.AccountDefault)
	clientSamePrice := ghtkclient.New(cfg.Env, cfg.AccountSamePrice)
	clientSamePrice2 := ghtkclient.New(cfg.Env, cfg.AccountSamePrice2)
	clients := map[byte]*ghtkclient.Client{
		'D': clientDefault,
		'S': clientSamePrice,
		'T': clientSamePrice2,
	}
	return &Carrier{
		clients:  clients,
		location: locationBus,
	}
}

func (c *Carrier) InitAllClients(ctx context.Context) error {
	for name, client := range c.clients {
		if err := client.Ping(); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "can not init client").
				WithMetap("client", name)
		}
	}
	return nil
}

func (p *Carrier) CreateFulfillment(ctx context.Context, order *ordermodel.Order, ffm *shipmodel.Fulfillment, args shipping_provider.GetShippingServicesArgs, service *model.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _err error) {

	note := shipping_provider.GetShippingProviderNote(order, ffm)
	weight := order.TotalWeight
	if weight == 0 {
		weight = 100
	}

	// set default value for GHTK
	// phần bảo hiểm truyền 0 hoặc ko truyền bên ghtk sẽ lấy theo giá trị tiền
	// thu hộ nên nếu ko muốn đóng bảo hiểm truyền ví dụ value=1000
	valueInsurance := 1000
	if order.ShopShipping.IncludeInsurance {
		valueInsurance = order.BasketValue
	}

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province

	// service, etopService, err := p.GetShippingService(ctx, args, order.ShopShipping.ProviderServiceID)
	// if err != nil {
	// 	return err
	// }

	// prepare products for ghtk
	var products []*ghtkclient.ProductRequest
	for _, line := range order.Lines {
		products = append(products, &ghtkclient.ProductRequest{
			Name:     line.ProductName,
			Price:    line.ListPrice,
			Quantity: line.Quantity,
		})
	}

	ghtkCmd := &CreateOrderCommand{
		ServiceID: service.ProviderServiceID,
		Request: &ghtkclient.CreateOrderRequest{
			Products: products,
			Order: &ghtkclient.OrderRequest{
				ID:           strconv.FormatInt(ffm.ID, 10),
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
				TotalWeight:  float32(weight),
				Value:        valueInsurance,
			},
		},
	}
	if ffm.AddressReturn != nil {
		returnQuery := &location.GetLocationQuery{DistrictCode: ffm.AddressReturn.DistrictCode}
		if err := p.location.Dispatch(ctx, returnQuery); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ trả hàng không hợp lệ: %v", err)
		}
		returnProvince, returnDistrict := returnQuery.Result.Province, returnQuery.Result.District

		ghtkCmd.Request.Order.UseReturnAddress = 1
		ghtkCmd.Request.Order.ReturnName = ffm.AddressReturn.GetFullName()
		ghtkCmd.Request.Order.ReturnAddress = cm.Coalesce(ffm.AddressReturn.Address1, ffm.AddressReturn.Address2)
		ghtkCmd.Request.Order.ReturnProvince = returnProvince.Name
		ghtkCmd.Request.Order.ReturnDistrict = returnDistrict.Name
		ghtkCmd.Request.Order.ReturnWard = ffm.AddressReturn.Ward
		ghtkCmd.Request.Order.ReturnTel = ffm.AddressReturn.Phone
		returnEmail := ffm.AddressReturn.Email
		if returnEmail == "" {
			returnEmail = "hotro@etop.vn"
		}
		// ReturnEmail can not empty
		ghtkCmd.Request.Order.ReturnEmail = returnEmail
	}

	if ghtkErr := p.CreateOrder(ctx, ghtkCmd); ghtkErr != nil {
		return nil, ghtkErr
	}
	r := ghtkCmd.Result

	now := time.Now()
	updateFfm := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ProviderServiceID:         service.ProviderServiceID,
		Status:                    model.S5SuperPos, // Now processing
		ShippingStatus:            model.S5SuperPos,
		ShippingFeeCustomer:       order.ShopShippingFee,
		ShippingFeeShop:           int(r.Order.Fee),
		ShippingCode:              NormalizeGHTKCode(r.Order.Label.String()),
		ExternalShippingName:      service.Name,
		ExternalShippingID:        r.Order.TrackingID.String(),
		ExternalShippingCode:      r.Order.Label.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       int(r.Order.Fee),
		ShippingState:             model.StateCreated,
		SyncStatus:                model.S4Positive,
		SyncStates: &model.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedPickAt:     service.ExpectedPickAt,
		ExpectedDeliveryAt: service.ExpectedDeliveryAt,
	}
	// ExpectedDeliveryAt
	expectedDeliveryAt, err := cm.FormatDateTimeEdgeCase(r.Order.EstimatedDeliverTime.String())
	if err == nil {
		updateFfm.ExpectedDeliveryAt = shipping.CalcDeliveryTime(model.TypeGHTK, toDistrict, *expectedDeliveryAt)
	}

	// prepare info to calc providerShippingFeeLines
	orderInfo := &ghtkclient.OrderInfo{
		LabelID:   r.Order.Label,
		ShipMoney: r.Order.Fee,
		Insurance: r.Order.InsuranceFee,
	}
	updateFfm.ProviderShippingFeeLines = CalcAndConvertShippingFeeLines(orderInfo)
	updateFfm.ShippingFeeShopLines = model.GetShippingFeeShopLines(updateFfm.ProviderShippingFeeLines, updateFfm.EtopPriceRule, &updateFfm.EtopAdjustedShippingFeeMain)

	return updateFfm, nil
}

func (p *Carrier) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment, action model.FfmAction) error {
	code := ffm.ExternalShippingCode
	cmd := &CancelOrderCommand{
		ServiceID: ffm.ProviderServiceID,
		LabelID:   code,
	}
	err := p.CancelOrder(ctx, cmd)
	return err
}

func (p *Carrier) GetShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province

	cmd := &CalcShippingFeeCommand{
		ArbitraryID:      args.AccountID,
		FromDistrictCode: args.FromDistrictCode,
		ToDistrictCode:   args.ToDistrictCode,
		Request: &ghtkclient.CalcShippingFeeRequest{
			Weight:          args.ChargeableWeight,
			Value:           args.GetInsuranceAmount(),
			PickingProvince: fromProvince.Name,
			PickingDistrict: fromDistrict.Name,
			Province:        toProvince.Name,
			District:        toDistrict.Name,
		},
	}
	err := p.CalcShippingFee(ctx, cmd)
	services := cmd.Result
	return services, err
}

func (p *Carrier) GetAllShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province

	cmd := &CalcShippingFeeCommand{
		ArbitraryID:      args.AccountID,
		FromDistrictCode: args.FromDistrictCode,
		ToDistrictCode:   args.ToDistrictCode,
		Request: &ghtkclient.CalcShippingFeeRequest{
			Weight:          args.ChargeableWeight,
			Value:           args.GetInsuranceAmount(),
			PickingProvince: fromProvince.Name,
			PickingDistrict: fromDistrict.Name,
			Province:        toProvince.Name,
			District:        toDistrict.Name,
		},
	}
	if err := p.CalcShippingFee(ctx, cmd); err != nil {
		return nil, err
	}
	providerServices := cmd.Result

	// get ETOP services
	etopServicesArgs := &etop_shipping_price.GetEtopShippingServicesArgs{
		ArbitraryID:  args.AccountID,
		Carrier:      model.TypeGHTK,
		FromProvince: fromProvince,
		ToProvince:   toProvince,
		ToDistrict:   toDistrict,
		Weight:       args.ChargeableWeight,
	}
	etopServices := etop_shipping_price.GetEtopShippingServices(etopServicesArgs)
	etopServices, _ = etop_shipping_price.FillInfoEtopServices(providerServices, etopServices)

	allServices := append(providerServices, etopServices...)
	return allServices, nil
}