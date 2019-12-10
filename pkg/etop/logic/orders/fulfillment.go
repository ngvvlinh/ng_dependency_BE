package orderS

import (
	"context"
	"sync"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/location"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	apishop "etop.vn/api/top/int/shop"
	"etop.vn/api/top/int/types"
	"etop.vn/api/top/types/etc/shipping"
	shipping_provider2 "etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ctrl               *shipping_provider.ProviderManager
	catalogQuery       catalog.QueryBus
	orderAggr          ordering.CommandBus
	customerAggr       customering.CommandBus
	customerQuery      customering.QueryBus
	traderAddressAggr  addressing.CommandBus
	traderAddressQuery addressing.QueryBus
	locationQuery      location.QueryBus
	eventBus           capi.EventBus
)

func Init(shippingProviderCtrl *shipping_provider.ProviderManager,
	catalogQueryBus catalog.QueryBus,
	orderAggregate ordering.CommandBus,
	customerAggregate customering.CommandBus,
	customerQueryBus customering.QueryBus,
	traderAddressAggregate addressing.CommandBus,
	traderAddressQueryBus addressing.QueryBus,
	locationQueryBus location.QueryBus,
	eventB capi.EventBus) {
	ctrl = shippingProviderCtrl
	catalogQuery = catalogQueryBus
	orderAggr = orderAggregate
	customerAggr = customerAggregate
	customerQuery = customerQueryBus
	traderAddressAggr = traderAddressAggregate
	traderAddressQuery = traderAddressQueryBus
	locationQuery = locationQueryBus
	eventBus = eventB
}

func ConfirmOrder(ctx context.Context, shop *model.Shop, r *apishop.ConfirmOrderRequest) (resp *types.Order, _err error) {
	autoInventoryVoucher := inventory.AutoInventoryVoucher(r.AutoInventoryVoucher.Apply(""))
	if !autoInventoryVoucher.ValidateAutoInventoryVoucher() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "AutoInventoryVoucher không hợp lệ, vui lòng kiểm tra lại. Giá trị hợp lệ: create | confirm")
	}
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

	if err := RaiseOrderConfirmingEvent(ctx, shop, autoInventoryVoucher, order); err != nil {
		return nil, err
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
			return resp, err
		}
		order.ConfirmStatus = status3.P
		order.ShopConfirm = status3.P

		event := &ordering.OrderConfirmedEvent{
			OrderID:              order.ID,
			AutoInventoryVoucher: autoInventoryVoucher,
			ShopID:               shop.ID,
			InventoryOverStock:   shop.InventoryOverstock.Apply(true),
		}
		if err := eventBus.Publish(ctx, event); err != nil {
			ll.Error("RaiseOrderConfirmedEvent", l.Error(err))
		}
	}
	resp = convertpb.PbOrder(order, nil, model.TagShop)
	if autoCreateFfm {
		req := &apishop.OrderIDRequest{
			OrderId: r.OrderId,
		}
		_res, err := ConfirmOrderAndCreateFulfillments(ctx, shop, 0, req)
		if err != nil {
			return nil, err
		}
		resp = _res.Order
	}
	return resp, nil
}

func ConfirmOrderAndCreateFulfillments(ctx context.Context, shop *model.Shop, partnerID dot.ID, r *apishop.OrderIDRequest) (resp *types.OrderWithErrorsResponse, _err error) {
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
	ffm, err := prepareFulfillmentFromOrder(ctx, order, shop)
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
	if cm.Env() == cm.EnvSandbox {
		if order.ShopShipping != nil &&
			order.ShopShipping.ShippingProvider != shipping_provider2.GHTK {
			go func() {
				time.Sleep(5 * time.Minute)
				_, err := CancelOrder(ctx, shop.ID, partnerID, order.ID, "Đơn hàng TEST, tự động huỷ", "")
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
		Fulfill:    ordertypes.Fulfill(ordermodel.FulfillShipment),
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

func RaiseOrderConfirmingEvent(ctx context.Context, shop *model.Shop, autoInventoryVoucher inventory.AutoInventoryVoucher, order *ordermodel.Order) error {
	orderLines := []*ordertypes.ItemLine{}
	for _, line := range order.Lines {
		if line.VariantID != 0 {
			_line := &ordertypes.ItemLine{
				OrderId:    line.OrderID,
				Quantity:   line.Quantity,
				ProductId:  line.ProductID,
				VariantId:  line.VariantID,
				IsOutside:  line.IsOutsideEtop,
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

func prepareFulfillmentFromOrder(ctx context.Context, order *ordermodel.Order, shop *model.Shop) (*shipmodel.Fulfillment, error) {
	if order.ShopShipping != nil && order.ShopShipping.ShippingProvider == shipping_provider2.GHN {
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
	var shopAddress *model.Address
	if order.ShopShipping != nil && order.ShopShipping.ShopAddress != nil {
		shopAddress, err = orderAddressToShippingAddress(order.ShopShipping.ShopAddress)
		if err != nil {
			return nil, cm.Error(cm.InvalidArgument, "Thông tin địa chỉ cửa hàng trong đơn hàng: "+err.Error()+"  Vui lòng cập nhật và thử lại.", err)
		}

	} else {
		if shop.ShipFromAddressID == 0 {
			return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
		}
		addressQuery := &model.GetAddressQuery{AddressID: shop.ShipFromAddressID}
		if err := bus.Dispatch(ctx, addressQuery); err != nil {
			return nil, cm.Error(cm.Internal, "Lỗi khi kiểm tra thông tin địa chỉ của cửa hàng: "+err.Error(), err)
		}
		shopAddress = addressQuery.Result
	}
	_, _, err = ctrl.VerifyDistrictCode(shopAddress)
	if err != nil {
		return nil, cm.Error(cm.FailedPrecondition, "Thông tin địa chỉ cửa hàng trong cấu hình cửa hàng: "+err.Error()+" Vui lòng cập nhật và thử lại.", nil)
	}

	if err := blockRachGiaDistrict(shopAddress); err != nil {
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
	return ffm, nil
}

// block create ffm from Rach Gia District: district_code = 899
func blockRachGiaDistrict(shopAddress *model.Address) error {
	if shopAddress.DistrictCode == "899" {
		return cm.Errorf(cm.InvalidArgument, nil, "Không thể lấy hàng tại địa chỉ này %v (%v)", shopAddress.District, shopAddress.Province)
	}
	return nil
}

func prepareSingleFulfillment(order *ordermodel.Order, shop *model.Shop, lines []*ordermodel.OrderLine, addressTo *model.Address) *shipmodel.Fulfillment {

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

	var addressReturn *model.Address
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

func orderAddressToShippingAddress(orderAddr *ordermodel.OrderAddress) (*model.Address, error) {
	if orderAddr == nil || orderAddr.DistrictCode == "" {
		return nil, cm.Error(cm.InvalidArgument, "Thiếu thông tin địa chỉ.", nil)
	}
	if orderAddr.Phone == "" {
		return nil, cm.Error(cm.InvalidArgument, "Thiếu thông tin số điện thoại.", nil)
	}
	if _, ok := validate.NormalizePhone(orderAddr.Phone); !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ (%v).", orderAddr.Phone)
	}

	return &model.Address{
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

func TryCancellingFulfillments(ctx context.Context, order *ordermodel.Order, fulfillments []*shipmodel.Fulfillment) ([]error, error) {
	var ffmToCancel []*shipmodel.Fulfillment
	ffmSendToProvider := make([]model.FfmAction, len(fulfillments))
	count := 0

	for i, ffm := range fulfillments {
		switch ffm.Status {
		case status5.P, status5.N, status5.NS:
			continue
		}

		if ffm.ShopConfirm != status3.N {
			ffmToCancel = append(ffmToCancel, ffm)
		}

		switch ffm.ShippingState {
		case shipping.Created, shipping.Picking:
			ffmSendToProvider[i] = model.FfmActionCancel
			count++

		case shipping.Holding,
			shipping.Delivering,
			shipping.Undeliverable,
			shipping.Returning:
			ffmSendToProvider[i] = model.FfmActionReturn
			count++

		default:
			ffmSendToProvider[i] = model.FfmActionNothing
		}
	}

	// update shop confirm
	if len(ffmToCancel) > 0 {
		ids := make([]dot.ID, len(ffmToCancel))
		for i, ffm := range ffmToCancel {
			ids[i] = ffm.ID
		}
		updateCmd := &shipmodelx.UpdateFulfillmentsStatusCommand{
			FulfillmentIDs: ids,
			ShopConfirm:    status3.N.Wrap(),
		}
		if err := bus.Dispatch(ctx, updateCmd); err != nil {
			return nil, err
		}
	}

	// MUSTDO: wait at least 10s before retry

	now := time.Now()
	var wg sync.WaitGroup
	var errs []error
	if count <= 0 {
		return nil, nil
	}

	errs = make([]error, len(ffmSendToProvider))
	for i, action := range ffmSendToProvider {
		// https://golang.org/doc/faq#closures_and_goroutines
		i, action, ffm := i, action, fulfillments[i]
		if action == model.FfmActionNothing {
			continue
		}

		wg.Add(1)
		go ignoreError(func() (_err error) {
			defer func() {
				wg.Done()
				errs[i] = _err
			}()

			// UpdateInfo to pending
			update := &shipmodel.Fulfillment{
				ID:         ffm.ID,
				SyncStatus: status4.S,
				SyncStates: &model.FulfillmentSyncStates{
					SyncAt:            now,
					NextShippingState: action.ToShippingState(),
				},
			}
			updateCmd := &shipmodelx.UpdateFulfillmentCommand{Fulfillment: update}
			if err := bus.Dispatch(ctx, updateCmd); err != nil {
				return cm.Errorf(cm.Internal, err, "Lỗi khi cập nhật vận đơn: %v", err.Error())
			}

			// TODO
			var shippingProviderErr error
			switch ffm.ShippingProvider {
			case shipping_provider2.GHN:
				shippingProviderErr = ctrl.GHN.CancelFulfillment(ctx, ffm, action)
			case shipping_provider2.GHTK:
				shippingProviderErr = ctrl.GHTK.CancelFulfillment(ctx, ffm, 0)
			case shipping_provider2.VTPost:
				shippingProviderErr = ctrl.VTPost.CancelFulfillment(ctx, ffm, 0)
			default:
				panic("Shipping provider was not supported.")
			}

			// Send
			if shippingProviderErr != nil {
				// UpdateInfo to error
				update2 := &shipmodel.Fulfillment{
					ID:         ffm.ID,
					SyncStatus: status4.N,
					SyncStates: &model.FulfillmentSyncStates{
						SyncAt: time.Now(),
						Error:  model.ToError(shippingProviderErr),

						NextShippingState: update.SyncStates.NextShippingState,
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
				ShippingState: update.SyncStates.NextShippingState,
				SyncStatus:    status4.P,
				SyncStates: &model.FulfillmentSyncStates{
					SyncAt: time.Now(),
				},
			}
			if update2.ShippingState == shipping.Cancelled {
				update2.Status = status5.N
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
