package ghn

import (
	"context"
	"strconv"
	"time"

	"etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	ghnclient "etop.vn/backend/pkg/integration/ghn/client"
	"etop.vn/backend/pkg/integration/shipping"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
)

var _ shipping_provider.ShippingProvider = &Carrier{}

type Carrier struct {
	clients  map[ClientType]*ghnclient.Client
	location location.Bus
}

func New(cfg Config, location location.Bus) *Carrier {
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
		ID:     ffm.ID,
		Status: model.S5SuperPos, // Now processing

		// TODO(qv): Handle shipping fee for supplier
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
	updateFfm.ShippingFeeShopLines = []*model.ShippingFeeLine{
		{
			ShippingFeeType:      model.ShippingFeeTypeMain,
			Cost:                 int(r.TotalServiceFee),
			ExternalShippingCode: r.OrderCode.String(),
		},
	}
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
	return c.GetShippingServices(ctx, args)
}

func (p *Carrier) GetShippingService(ffm *shipmodel.Fulfillment, order *ordermodel.Order, weight int, valueInsurance int) (providerService *model.AvailableShippingService, etopService *model.AvailableShippingService, err error) {
	return nil, nil, cm.ErrTODO
}

func calcGHNService(order *ordermodel.Order, ffm *shipmodel.Fulfillment, cmd *RequestFindAvailableServicesCommand) (serviceID string, service *model.AvailableShippingService, err error) {

	// Always choose the fastest possible service
	var minTime time.Time
	services := cmd.Result

	if ffm.SupplierID != 0 {
		if len(services) > 0 {
			service = services[0]
			minTime = service.ExpectedDeliveryAt
		}
		for _, s := range services {
			if t := s.ExpectedDeliveryAt; t.Before(minTime) {
				minTime = t
				service = s
			}
		}
		if service == nil {
			return "", nil, cm.Errorf(cm.ExternalServiceError, nil,
				"Lỗi từ Giao Hàng Nhanh: Không thể chọn được gói dịch vụ giao hàng (từ %v, %v đến %v, %v).",
				ffm.AddressFrom.District, ffm.AddressFrom.Province,
				ffm.AddressTo.District, ffm.AddressTo.Province,
			)
		}
		return service.ProviderServiceID, service, err
	}
	if order.ShopShipping != nil {
		providerServiceID := cm.Coalesce(order.ShopShipping.ProviderServiceID, order.ShopShipping.ExternalServiceID)
		for _, s := range services {
			if s.ProviderServiceID == providerServiceID {
				service = s
				minTime = service.ExpectedDeliveryAt
				break
			}
		}
		if service == nil {
			return "", nil, cm.Errorf(cm.InvalidArgument, nil, "Gói dịch vụ giao hàng đã chọn không hợp lệ")
		}
		if order.ShopShipping.ExternalShippingFee != int(service.ServiceFee) {
			return "", nil, cm.Errorf(cm.InvalidArgument, nil,
				"Số tiền phí giao hàng không hợp lệ cho dịch vụ %v: Phí trên đơn hàng %v, phí từ Giao Hàng Nhanh: %v",
				service.Name, order.ShopShipping.ExternalShippingFee, int(service.ServiceFee))
		}
		return providerServiceID, service, nil
	}
	return "", nil, cm.Errorf(cm.InvalidArgument, nil, "Cần chọn gói dịch vụ giao hàng")
}
