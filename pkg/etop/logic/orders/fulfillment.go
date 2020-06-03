package orderS

import (
	"context"
	"sync"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	apishop "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/shipping"
	typeshippingprovider "o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	addressmodel "o.o/backend/com/main/address/model"
	addressmodelx "o.o/backend/com/main/address/modelx"
	identitymodel "o.o/backend/com/main/identity/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/com/main/shipping/carrier"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

type OrderLogic struct {
}
type FlagFaboOrderAutoConfirmPaymentStatus bool

var (
	ctrl                                  *shipping_provider.CarrierManager
	catalogQuery                          catalog.QueryBus
	orderAggr                             ordering.CommandBus
	customerAggr                          customering.CommandBus
	customerQuery                         customering.QueryBus
	traderAddressAggr                     addressing.CommandBus
	traderAddressQuery                    addressing.QueryBus
	locationQuery                         location.QueryBus
	eventBus                              capi.EventBus
	shipmentManager                       *carrier.ShipmentManager
	flagFaboOrderUpdatePaymentSatusConfig FlagFaboOrderAutoConfirmPaymentStatus
)

func New(shippingProviderCtrl *shipping_provider.CarrierManager,
	catalogQueryBus catalog.QueryBus,
	orderAggregate ordering.CommandBus,
	customerAggregate customering.CommandBus,
	customerQueryBus customering.QueryBus,
	traderAddressAggregate addressing.CommandBus,
	traderAddressQueryBus addressing.QueryBus,
	locationQueryBus location.QueryBus,
	eventB capi.EventBus,
	flagFaboOrderUpdatePaymentSatus FlagFaboOrderAutoConfirmPaymentStatus,
	shipmentCarrierCtrl *carrier.ShipmentManager) *OrderLogic {
	ctrl = shippingProviderCtrl
	catalogQuery = catalogQueryBus
	orderAggr = orderAggregate
	customerAggr = customerAggregate
	customerQuery = customerQueryBus
	traderAddressAggr = traderAddressAggregate
	traderAddressQuery = traderAddressQueryBus
	locationQuery = locationQueryBus
	eventBus = eventB
	flagFaboOrderUpdatePaymentSatusConfig = flagFaboOrderUpdatePaymentSatus
	shipmentManager = shipmentCarrierCtrl
	return &OrderLogic{}
}

var blockCarrierByDistricts = map[typeshippingprovider.ShippingProvider][]string{
	typeshippingprovider.GHN: []string{},
}

var blockCarrierByProvinces = map[typeshippingprovider.ShippingProvider][]string{
	typeshippingprovider.GHN: []string{},
}

func (s *OrderLogic) ConfirmOrder(ctx context.Context, userID dot.ID, shop *identitymodel.Shop, r *apishop.ConfirmOrderRequest) (resp *types.Order, _err error) {
	autoCreateFfm := r.AutoCreateFulfillment

	query := &ordermodelx.GetOrderQuery{
		OrderID: r.OrderId,
		ShopID:  shop.ID,
	}
	resp = &types.Order{}
	if err := bus.Dispatch(ctx, query); err != nil {
		return resp, err
	}
	order := query.Result.Order
	switch order.Status {
	case status5.N:
		return resp, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hủy")
	case status5.P:
		return resp, cm.Error(cm.FailedPrecondition, "Đơn hàng đã hoàn thành.", nil)
	case status5.NS:
		return resp, cm.Error(cm.FailedPrecondition, "Đơn hàng đã trả hàng.", nil)
	}
	if order.ConfirmStatus == status3.N || order.ShopConfirm == status3.N {
		return resp, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hủy")
	}

	if err := s.RaiseOrderConfirmingEvent(ctx, shop, r.AutoInventoryVoucher, order); err != nil {
		return nil, err
	}

	// Only update order status when success.
	// This disallow updating order.
	if order.ConfirmStatus != status3.P ||
		order.ShopConfirm != status3.P {
		paymentStatus := status4.Z
		if flagFaboOrderUpdatePaymentSatusConfig {
			paymentStatus = status4.P
		}
		cmd := &ordermodelx.UpdateOrdersStatusCommand{
			OrderIDs:      []dot.ID{r.OrderId},
			ConfirmStatus: status3.P.Wrap(),
			ShopConfirm:   status3.P.Wrap(),
			PaymentStatus: paymentStatus.Wrap(),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return resp, err
		}
		order.ConfirmStatus = status3.P
		order.ShopConfirm = status3.P
		order.PaymentStatus = paymentStatus
		event := &ordering.OrderConfirmedEvent{
			OrderID:              order.ID,
			AutoInventoryVoucher: r.AutoInventoryVoucher,
			ShopID:               shop.ID,
			InventoryOverStock:   shop.InventoryOverstock.Apply(true),
			UpdatedBy:            userID,
		}
		if err := eventBus.Publish(ctx, event); err != nil {
			ll.Error("RaiseOrderConfirmedEvent", l.Error(err))
		}
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	resp = convertpb.PbOrder(query.Result.Order, nil, model.TagShop)
	resp.ShopName = shop.Name
	if autoCreateFfm {
		req := &apishop.OrderIDRequest{
			OrderId: r.OrderId,
		}
		_res, err := s.ConfirmOrderAndCreateFulfillments(ctx, userID, shop, 0, req)
		if err != nil {
			return nil, err
		}
		resp = _res.Order
	}
	return resp, nil
}

func (s *OrderLogic) ConfirmOrderAndCreateFulfillments(ctx context.Context, userID dot.ID, shop *identitymodel.Shop, partnerID dot.ID, r *apishop.OrderIDRequest) (resp *types.OrderWithErrorsResponse, _err error) {
	shopID := shop.ID
	resp = &types.OrderWithErrorsResponse{}
	query := &ordermodelx.GetOrderQuery{
		ShopID:             shopID,
		PartnerID:          partnerID,
		OrderID:            r.OrderId,
		IncludeFulfillment: true,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return resp, err
	}
	order := query.Result.Order

	// Verify status
	switch order.Status {
	case status5.N:
		return resp, cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil)
	case status5.P:
		return resp, cm.Error(cm.FailedPrecondition, "Đơn hàng đã hoàn thành.", nil)
	case status5.NS:
		return resp, cm.Error(cm.FailedPrecondition, "Đơn hàng đã trả hàng.", nil)
	}

	if order.ConfirmStatus == status3.N ||
		order.ShopConfirm == status3.N {
		return resp, cm.Error(cm.FailedPrecondition, "Đơn hàng đã huỷ.", nil)
	}

	// Fill response
	fulfillments := query.Result.Fulfillments
	defer func() {
		if _err != nil {
			return
		}

		resp.Order = convertpb.PbOrder(order, fulfillments, model.TagShop)
		resp.Order.ShopName = "" // TODO: remove this line
	}()

	// Create fulfillments
	ffm, err := s.prepareFulfillmentFromOrder(ctx, order, shop)
	if err != nil {
		return resp, err
	}
	// Compare fulfillments for retry/update
	creates, updates, err := compareFulfillments(order, query.Result.Fulfillments, ffm)
	if err != nil {
		return resp, err
	}

	if creates != nil {
		ffmCmd := &shipmodelx.CreateFulfillmentsCommand{
			Fulfillments: creates,
		}
		if err := bus.Dispatch(ctx, ffmCmd); err != nil {
			return resp, err
		}
	}
	if updates != nil {
		ffmCmd := &shipmodelx.UpdateFulfillmentsCommand{
			Fulfillments: updates,
		}
		if err := bus.Dispatch(ctx, ffmCmd); err != nil {
			return resp, err
		}
	}

	ll.S.Infof("Compare fulfillments: create %v update %v", len(creates), len(updates))
	totalChanges := len(creates) + len(updates)
	if totalChanges == 0 {
		return resp, nil
	}

	ffms := append(creates, updates...)
	if err := ctrl.CreateExternalShipping(ctx, order, ffms); err != nil {
		return resp, err
	}
	// automatically cancel orders on sandbox for ghn and vtpost
	if cmenv.Env() == cmenv.EnvSandbox {
		if order.ShopShipping != nil &&
			order.ShopShipping.ShippingProvider != typeshippingprovider.GHTK {
			go func() {
				time.Sleep(5 * time.Minute)
				_, err := s.CancelOrder(ctx, userID, shop.ID, partnerID, order.ID, "Đơn hàng TEST, tự động huỷ", inventory_auto.Unknown)
				if err != nil {
					ll.Error("Can not cancel order on sandbox", l.Error(err))
				}
			}()
		}
	}

	// Only update order status when success.
	// This disallow updating order.
	if order.ConfirmStatus != status3.P ||
		order.ShopConfirm != status3.P {
		cmd := &ordermodelx.UpdateOrdersStatusCommand{
			OrderIDs:      []dot.ID{r.OrderId},
			ConfirmStatus: status3.P.Wrap(),
			ShopConfirm:   status3.P.Wrap(),
		}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			_err = err
		}
		order.ConfirmStatus = status3.P
		order.ShopConfirm = status3.P
	}

	// update order fulfillment_type: `shipment`
	var ffmIDs []dot.ID
	for _, _ffm := range ffms {
		ffmIDs = append(ffmIDs, _ffm.ID)
	}
	cmd := &ordering.ReserveOrdersForFfmCommand{
		OrderIDs:   []dot.ID{order.ID},
		Fulfill:    ordertypes.ShippingTypeShipment,
		FulfillIDs: ffmIDs,
	}
	if err = orderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	// Get order again
	if err := bus.Dispatch(ctx, query); err != nil {
		return resp, err
	}
	order = query.Result.Order
	fulfillments = query.Result.Fulfillments

	// Order will be filled by above defer
	return resp, nil
}

func (s *OrderLogic) RaiseOrderConfirmingEvent(ctx context.Context, shop *identitymodel.Shop, autoInventoryVoucher inventory_auto.AutoInventoryVoucher, order *ordermodel.Order) error {
	orderLines := []*ordertypes.ItemLine{}
	for _, line := range order.Lines {
		if line.VariantID != 0 {
			_line := &ordertypes.ItemLine{
				OrderID:    line.OrderID,
				Quantity:   line.Quantity,
				ProductID:  line.ProductID,
				VariantID:  line.VariantID,
				IsOutSide:  line.IsOutsideEtop,
				TotalPrice: line.TotalLineAmount,
			}
			orderLines = append(orderLines, _line)
		}
	}
	event := &ordering.OrderConfirmingEvent{
		OrderID:              order.ID,
		ShopID:               shop.ID,
		InventoryOverStock:   shop.InventoryOverstock.Apply(true),
		Lines:                orderLines,
		AutoInventoryVoucher: autoInventoryVoucher,
	}
	if err := eventBus.Publish(ctx, event); err != nil {
		return err
	}
	return nil
}

func (s *OrderLogic) prepareFulfillmentFromOrder(ctx context.Context, order *ordermodel.Order, shop *identitymodel.Shop) (*shipmodel.Fulfillment, error) {
	if order.ShopShipping != nil && order.ShopShipping.ShippingProvider == typeshippingprovider.GHN {
		if order.TryOn == 0 && order.GhnNoteCode == 0 {
			return nil, cm.Error(cm.FailedPrecondition, "Vui lòng chọn ghi chú xem hàng!", nil)
		}
	}

	if !model.VerifyOrderSource(order.OrderSourceType) {
		return nil, cm.Error(cm.FailedPrecondition, "Không thể xác định nguồn đơn hàng!", nil)
	}
	addressTo, err := orderAddressToShippingAddress(order.ShippingAddress)
	if err != nil {
		return nil, cm.Error(cm.InvalidArgument, "Thông tin địa chỉ người nhận: "+err.Error()+" Vui lòng cập nhật và thử lại.", err)
	}
	if _, _, err := ctrl.VerifyDistrictCode(addressTo); err != nil {
		return nil, cm.Error(cm.InvalidArgument, "Thông tin địa chỉ người nhận: "+err.Error()+" Vui lòng cập nhật và thử lại.", nil)
	}

	ffm := prepareSingleFulfillment(order, shop, order.Lines, addressTo)

	// Use shop address from order or from shop default address
	var shopAddress *addressmodel.Address
	if order.ShopShipping != nil && order.ShopShipping.ShopAddress != nil {
		shopAddress, err = orderAddressToShippingAddress(order.ShopShipping.ShopAddress)
		if err != nil {
			return nil, cm.Error(cm.InvalidArgument, "Thông tin địa chỉ cửa hàng trong đơn hàng: "+err.Error()+"  Vui lòng cập nhật và thử lại.", err)
		}

	} else {
		if shop.ShipFromAddressID == 0 {
			return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
		}
		addressQuery := &addressmodelx.GetAddressQuery{AddressID: shop.ShipFromAddressID}
		if err := bus.Dispatch(ctx, addressQuery); err != nil {
			return nil, cm.Error(cm.Internal, "Lỗi khi kiểm tra thông tin địa chỉ của cửa hàng: "+err.Error(), err)
		}
		shopAddress = addressQuery.Result
	}
	_, _, err = ctrl.VerifyDistrictCode(shopAddress)
	if err != nil {
		return nil, cm.Error(cm.FailedPrecondition, "Thông tin địa chỉ cửa hàng trong cấu hình cửa hàng: "+err.Error()+" Vui lòng cập nhật và thử lại.", nil)
	}

	if err := s.checkBlockCarrier(shopAddress, order.ShopShipping.ShippingProvider); err != nil {
		return nil, err
	}

	ffm.TotalAmount = order.TotalAmount
	ffm.TotalDiscount = order.TotalDiscount
	ffm.AddressFrom = shopAddress

	if order.ShopCOD < 0 {
		return nil, cm.Error(cm.InvalidArgument, "Thông tin tiền thu hộ (COD) không hợp lệ.", nil)
	}
	if order.ShopCOD > 0 {
		ffm.TotalCODAmount = order.ShopCOD
		ffm.OriginalCODAmount = order.ShopCOD
	}

	if order.ShopShipping == nil {
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng chọn dịch vụ giao hàng.", nil)
	}
	ffm.CreatedBy = order.CreatedBy
	return ffm, nil
}

func (s *OrderLogic) checkBlockCarrier(shopAddress *addressmodel.Address, provider typeshippingprovider.ShippingProvider) error {
	if provider == 0 {
		return nil
	}
	provinces := blockCarrierByProvinces[provider]
	if cm.StringsContain(provinces, shopAddress.ProvinceCode) {
		return cm.Errorf(cm.InvalidArgument, nil, "%v không thể lấy hàng tại địa chỉ này %v (%v)", provider.Label(), shopAddress.District, shopAddress.Province)
	}
	districts, ok := blockCarrierByDistricts[provider]
	if !ok {
		return nil
	}
	if cm.StringsContain(districts, shopAddress.DistrictCode) {
		return cm.Errorf(cm.InvalidArgument, nil, "%v không thể lấy hàng tại địa chỉ này %v (%v)", provider.Label(), shopAddress.District, shopAddress.Province)
	}
	return nil
}

func prepareSingleFulfillment(order *ordermodel.Order, shop *identitymodel.Shop, lines []*ordermodel.OrderLine, addressTo *addressmodel.Address) *shipmodel.Fulfillment {

	var variantIDs []dot.ID
	totalItems, totalWeight, basketValue, totalAmount := 0, 0, 0, 0

	if len(order.Lines) != 0 {
		for _, line := range order.Lines {
			variantIDs = append(variantIDs, line.VariantID)
			totalItems += line.Quantity
			totalWeight += line.Weight
			basketValue += line.RetailPrice * line.Quantity
			totalAmount += line.PaymentPrice * line.Quantity
		}
	} else {
		totalItems = order.TotalItems
		totalWeight = order.TotalWeight
		basketValue = order.BasketValue
		totalAmount = order.TotalAmount
		variantIDs = []dot.ID{}
	}

	typeFrom := model.FFShop
	typeTo := model.FFCustomer

	ffmID := cm.NewID()
	shippingProvider := order.ShopShipping.ShippingProvider
	providerServiceID := cm.Coalesce(order.ShopShipping.ProviderServiceID, order.ShopShipping.ExternalServiceID)

	var addressReturn *addressmodel.Address
	if order.ShopShipping.ReturnAddress != nil {
		addressReturn, _ = orderAddressToShippingAddress(order.ShopShipping.ReturnAddress)
	}

	fulfillment := &shipmodel.Fulfillment{
		ID:                ffmID,
		OrderID:           order.ID,
		ShopID:            shop.ID,
		PartnerID:         order.PartnerID,
		ShopConfirm:       status3.P, // Always set shop_confirm to 1
		ConfirmStatus:     0,
		TotalItems:        totalItems,
		TotalWeight:       order.TotalWeight,
		BasketValue:       basketValue,
		TotalDiscount:     0,
		TotalAmount:       totalAmount,
		TotalCODAmount:    0,
		OriginalCODAmount: order.ShopCOD,
		// We only support shop cod
		// TotalCODAmount: totalCODAmount,

		ShippingFeeCustomer:      0, // only fill the first fulfillment
		ShippingFeeShop:          0, // after calling GHN
		ShippingFeeShopLines:     nil,
		ShippingServiceFee:       0,
		ExternalShippingFee:      0, // after calling GHN
		ProviderShippingFeeLines: nil,
		EtopDiscount:             0,
		EtopFeeAdjustment:        0,
		VariantIDs:               variantIDs,
		Lines:                    lines,
		// after this
		TypeFrom:                           typeFrom,
		TypeTo:                             typeTo,
		AddressFrom:                        nil, // will be filled later
		AddressTo:                          addressTo,
		AddressReturn:                      addressReturn,
		AddressToProvinceCode:              "",
		AddressToDistrictCode:              "",
		AddressToWardCode:                  "",
		CreatedAt:                          time.Time{}, // automatically by sq
		UpdatedAt:                          time.Time{}, // automatically by sq
		ShippingCancelledAt:                time.Time{},
		ClosedAt:                           time.Time{},
		ShippingDeliveredAt:                time.Time{},
		ShippingReturnedAt:                 time.Time{},
		ExpectedDeliveryAt:                 time.Time{},
		ExpectedPickAt:                     time.Time{},
		CODEtopTransferedAt:                time.Time{},
		ShippingFeeShopTransferedAt:        time.Time{},
		MoneyTransactionID:                 0,
		MoneyTransactionShippingExternalID: 0,
		CancelReason:                       "",
		ShippingProvider:                   shippingProvider,
		ProviderServiceID:                  providerServiceID,
		ShippingCode:                       "", // after calling GHN
		ShippingNote:                       order.ShippingNote,
		TryOn:                              order.GetTryOn(),
		IncludeInsurance:                   order.ShopShipping.IncludeInsurance,

		// After calling GHN
		ExternalShippingName:        "",
		ExternalShippingID:          "",
		ExternalShippingCode:        "",
		ExternalShippingCreatedAt:   time.Time{},
		ExternalShippingUpdatedAt:   time.Time{},
		ExternalShippingCancelledAt: time.Time{},
		ExternalShippingDeliveredAt: time.Time{},
		ExternalShippingReturnedAt:  time.Time{},
		ExternalShippingClosedAt:    time.Time{},
		ExternalShippingState:       "",
		ExternalShippingStateCode:   "",
		ExternalShippingStatus:      0,
		ExternalShippingNote:        "",
		ExternalShippingSubState:    "",
		ExternalShippingData:        nil,
		ShippingState:               shipping.Default,
		ShippingStatus:              0,
		EtopPaymentStatus:           0,
		Status:                      0,
		SyncStatus:                  0,
		SyncStates:                  nil,
		LastSyncAt:                  time.Time{},
		ExternalShippingLogs:        nil,
	}
	return fulfillment
}

// Compare the old fulfillments and expected fulfillments
// - Missing/cancelled fulfillments: create
// - Error fulfillments: update
// - Processing fulfillments: ignore
func compareFulfillments(order *ordermodel.Order, olds []*shipmodel.Fulfillment, ffm *shipmodel.Fulfillment) (creates, updates []*shipmodel.Fulfillment, err error) {
	// active ffm: Those which are not cancelled
	var old *shipmodel.Fulfillment
	for _, oldFfm := range olds {
		if oldFfm.Status != status5.N && oldFfm.Status != status5.NS &&
			oldFfm.ShopConfirm != status3.N {
			old = oldFfm
		}
	}

	if old == nil {
		return []*shipmodel.Fulfillment{ffm}, nil, nil
	}
	// update error fulfillments
	if old.Status == status5.Z {
		ffm.ID = old.ID
		return nil, []*shipmodel.Fulfillment{ffm}, nil
	}
	// ignore
	return nil, nil, nil
}

func orderAddressToShippingAddress(orderAddr *ordermodel.OrderAddress) (*addressmodel.Address, error) {
	if orderAddr == nil || orderAddr.DistrictCode == "" {
		return nil, cm.Error(cm.InvalidArgument, "Thiếu thông tin địa chỉ.", nil)
	}
	if orderAddr.Phone == "" {
		return nil, cm.Error(cm.InvalidArgument, "Thiếu thông tin số điện thoại.", nil)
	}
	if _, ok := validate.NormalizePhone(orderAddr.Phone); !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ (%v).", orderAddr.Phone)
	}

	return &addressmodel.Address{
		ID:           0,
		FullName:     orderAddr.FullName,
		FirstName:    orderAddr.FirstName,
		LastName:     orderAddr.LastName,
		Phone:        orderAddr.Phone,
		Position:     "",
		Email:        "",
		Country:      orderAddr.Country,
		City:         orderAddr.City,
		Province:     orderAddr.Province,
		District:     orderAddr.District,
		Ward:         orderAddr.Ward,
		Zip:          orderAddr.Zip,
		DistrictCode: orderAddr.DistrictCode,
		ProvinceCode: orderAddr.ProvinceCode,
		WardCode:     orderAddr.WardCode,
		Company:      orderAddr.Company,
		Address1:     orderAddr.Address1,
		Address2:     orderAddr.Address2,
		Type:         "",
		AccountID:    0,
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}, nil
}

func (s *OrderLogic) TryCancellingFulfillments(ctx context.Context, order *ordermodel.Order, fulfillments []*shipmodel.Fulfillment) ([]error, error) {
	var ffmToCancel []*shipmodel.Fulfillment

	for _, ffm := range fulfillments {
		switch ffm.Status {
		case status5.P, status5.N, status5.NS:
			continue
		}

		if ffm.ShopConfirm != status3.N {
			ffmToCancel = append(ffmToCancel, ffm)
		}
	}

	var wg sync.WaitGroup
	var errs []error

	errs = make([]error, len(ffmToCancel))
	for i, ffm := range ffmToCancel {
		i, ffm := i, ffm // https://golang.org/doc/faq#closures_and_goroutines

		wg.Add(1)
		go ignoreError(func() (_err error) {
			defer func() {
				wg.Done()
				errs[i] = _err
			}()

			var shippingProviderErr error
			if ffm.ShippingType == 0 {
				driver := ctrl.GetShippingProviderDriver(ffm.ShippingProvider)
				if driver == nil {
					panic("Shipping provider was not supported.")
				}
				shippingProviderErr = driver.CancelFulfillment(ctx, ffm, model.FfmActionCancel)
			} else if ffm.ConnectionID != 0 {
				shippingProviderErr = shipmentManager.CancelFulfillment(ctx, ffm)
			}

			// Send
			if shippingProviderErr != nil {
				// UpdateInfo to error
				update2 := &shipmodel.Fulfillment{
					ID:         ffm.ID,
					SyncStatus: status4.N,
					SyncStates: &shippingsharemodel.FulfillmentSyncStates{
						SyncAt: time.Now(),
						Error:  model.ToError(shippingProviderErr),

						NextShippingState: shipping.Cancelled,
					},
				}
				update2Cmd := &shipmodelx.UpdateFulfillmentCommand{Fulfillment: update2}
				if err := bus.Dispatch(ctx, update2Cmd); err != nil {
					return err
				}
				return shippingProviderErr
			}

			// UpdateInfo to ok
			update2 := &shipmodel.Fulfillment{
				ID:            ffm.ID,
				ShippingState: shipping.Cancelled,
				SyncStatus:    status4.P,
				SyncStates: &shippingsharemodel.FulfillmentSyncStates{
					SyncAt: time.Now(),
				},
				ShopConfirm:  status3.N,
				Status:       status5.N,
				CancelReason: order.CancelReason,
			}
			update2Cmd := &shipmodelx.UpdateFulfillmentCommand{Fulfillment: update2}
			if err := bus.Dispatch(ctx, update2Cmd); err != nil {
				return err
			}
			return nil
		}())
	}
	wg.Wait()
	return errs, nil
}

func ignoreError(err error) {}
