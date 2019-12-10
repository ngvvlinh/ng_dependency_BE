package vtpost

import (
	"context"
	"time"

	"etop.vn/api/main/location"
	"etop.vn/api/top/types/etc/shipping"
	shipping_provider2 "etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/authorize/login"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	shippingprovider "etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	vtpostclient "etop.vn/backend/pkg/integration/shipping/vtpost/client"
	"etop.vn/capi/dot"
)

var _ shippingprovider.ShippingProvider = &Carrier{}

type Carrier struct {
	clients  map[byte]vtpostclient.Client
	location location.QueryBus
}

func New(cfg Config, locationBus location.QueryBus) *Carrier {
	clientDefault := vtpostclient.New(cfg.Env, cfg.AccountDefault)
	clients := map[byte]vtpostclient.Client{
		VTPostCodePublic: clientDefault,
	}

	return &Carrier{
		clients:  clients,
		location: locationBus,
	}
}

func (c *Carrier) InitAllClients(ctx context.Context) error {
	for code, client := range c.clients {
		if err := client.Ping(); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "VTPost: can not init client").
				WithMetap("client", code)
		}
		err := CreateShippingSource(code, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateShippingSource(code byte, client vtpostclient.Client) error {
	ctx := context.Background()
	generator := newServiceIDGenerator(SecretCode)
	// generate a default clientName to save to db
	clientName, err := generator.GenerateServiceID(code, vtpostclient.OrderServiceCodeSCOD)
	if err != nil {
		return err
	}

	cmd := &model.CreateShippingSource{
		Name:     clientName,
		Type:     shipping_provider2.VTPost,
		Username: client.GetUserName(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	states := shippingprovider.Convert_model_ShippingSourceInternal_To_vtpost_ClientStates(cmd.Result.ShippingSourceInternal)
	client.InitFromSavedStates(*states)

	changed, err := client.AutoLoginAndRefreshToken(ctx)
	if err != nil {
		return err
	}
	if changed {
		newStates := client.GetStatesForSerialization()
		updateCmd := &model.UpdateOrCreateShippingSourceInternal{
			ID:          cmd.Result.ShippingSource.ID,
			AccessToken: newStates.AccessToken,
			ExpiresAt:   newStates.ExpiresAt,
			LastSyncAt:  newStates.AccessTokenCreatedAt,
			Secret: &model.ShippingSourceSecret{
				CustomerID: newStates.CustomerID,
				Username:   cmd.Username,
				Password:   login.EncodePassword(cmd.Password),
			},
		}
		if err := bus.Dispatch(ctx, updateCmd); err != nil {
			return err
		}
	}

	return nil
}

func (c *Carrier) CreateFulfillment(ctx context.Context, order *ordermodel.Order, ffm *shipmodel.Fulfillment, args shippingprovider.GetShippingServicesArgs, service *model.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _err error) {
	if ffm.AddressReturn != nil {
		// vtpost does not support address_return
		return nil, cm.Error(cm.ExternalServiceError, "VTPost không hỗ trợ địa chỉ trả hàng. Vui lòng để trống thông tin này.", nil)
	}

	note := shippingprovider.GetShippingProviderNote(order, ffm)
	weight := order.TotalWeight
	if weight == 0 {
		weight = 100
	}
	providerServiceID := cm.Coalesce(order.ShopShipping.ProviderServiceID, order.ShopShipping.ExternalServiceID)

	// add total COD amount when get shipping service (it use for calc: phí thu hộ)
	// this is different with another provider, we must get providerResponse when calc shipping fee
	maxValueFreeInsuranceFee := c.GetMaxValueFreeInsuranceFee()
	valueInsurance := args.GetInsuranceAmount(maxValueFreeInsuranceFee)

	if ffm.AddressFrom.WardCode == "" || ffm.AddressTo.WardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ViettelPost yêu cầu thông tin phường xã hợp lệ để giao hàng")
	}

	fromQuery := &location.GetLocationQuery{
		DistrictCode: ffm.AddressFrom.DistrictCode,
		WardCode:     ffm.AddressFrom.WardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: ffm.AddressTo.DistrictCode,
		WardCode:     ffm.AddressTo.WardCode,
	}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromWard, fromDistrict, fromProvince := fromQuery.Result.Ward, fromQuery.Result.District, fromQuery.Result.Province
	toWard, toDistrict, toProvince := toQuery.Result.Ward, toQuery.Result.District, toQuery.Result.Province

	if fromWard.VtpostId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ViettelPost không thể lấy hàng tại địa chỉ này (%v, %v, %v)", fromWard.Name, fromDistrict.Name, fromProvince.Name)
	}
	if toWard.VtpostId == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ViettelPost không thể giao hàng tại địa chỉ này (%v, %v, %v)", toWard.Name, toDistrict.Name, toProvince.Name)
	}

	deliveryDate := time.Now()
	deliveryDate.Add(30 * time.Minute)

	// prepare products for vtpost
	var products []*vtpostclient.Product
	var productName string
	for _, line := range order.Lines {
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

	vtpostCmd := &CreateOrderArgs{
		ServiceID: providerServiceID,
		Request: &vtpostclient.CreateOrderRequest{
			OrderNumber: "", // will be filled later
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
			ProductPrice:       valueInsurance,
			ProductWeight:      weight,
			OrderNote:          note,
			MoneyCollection:    ffm.TotalCODAmount,
			MoneyTotalFee:      service.ServiceFee,
			ListItem:           products,
			ProductName:        productName,
			ProductDescription: productName,
		},
	}

	shippingCode, err := sqlstore.GenerateVtpostShippingCode()
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can not generate shipping code for ffm.")
	}
	vtpostCmd.Request.OrderNumber = shippingCode

	if vtpostErr := c.createOrder(ctx, vtpostCmd); vtpostErr != nil {
		return nil, vtpostErr
	}
	r := vtpostCmd.Result

	now := time.Now()
	updateFfm := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		Status:                    status5.S, // Now processing
		ShippingStatus:            status5.S,
		ShippingFeeCustomer:       order.ShopShippingFee,
		ShippingServiceFee:        r.Data.MoneyTotal,
		ShippingFeeShop:           r.Data.MoneyTotal,
		ShippingCode:              r.Data.OrderNumber,
		ExternalShippingName:      "", // TODO
		ExternalShippingID:        r.Data.OrderNumber,
		ExternalShippingCode:      r.Data.OrderNumber,
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		// ShippingFeeShopLines:      shopShippingFeeLines,
		// ProviderShippingFeeLines:  providerResponse.ShippingFeeLines,
		ExpectedPickAt:     service.ExpectedPickAt,
		ExpectedDeliveryAt: service.ExpectedDeliveryAt,
		ShippingState:      shipping.Created,
		SyncStatus:         status4.P,
		SyncStates: &model.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
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
		updateFfm.ShippingFeeShopLines = model.GetShippingFeeShopLines(lines, false, dot.NullInt{})
	}

	return updateFfm, nil
}

func (c *Carrier) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment, action model.FfmAction) error {
	code := ffm.ExternalShippingCode
	cmd := &CancelOrderCommand{
		ServiceID: ffm.ProviderServiceID,
		Request: &vtpostclient.CancelOrderRequest{
			OrderNumber: code,
		},
	}
	return c.cancelOrder(ctx, cmd)
}

func (c *Carrier) GetShippingServices(ctx context.Context, args shippingprovider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {

	fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}

	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	maxValueFreeInsuranceFee := c.GetMaxValueFreeInsuranceFee()

	cmd := &CalcShippingFeeAllServicesArgs{
		ArbitraryID:  args.AccountID,
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
	err := c.CalcShippingFee(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}

func (c *Carrier) GetAllShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	return c.GetShippingServices(ctx, args)
}

func (c *Carrier) GetMaxValueFreeInsuranceFee() int {
	// Follow the policy of provider
	return 0
}
