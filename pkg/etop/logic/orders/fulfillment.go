package orderS

import (
	"context"
	"sync"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	apishop "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_payment_type"
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
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

type OrderLogic struct {
	CatalogQuery       catalog.QueryBus
	OrderAggr          ordering.CommandBus
	CustomerAggr       customering.CommandBus
	CustomerQuery      customering.QueryBus
	TraderAddressAggr  addressing.CommandBus
	TraderAddressQuery addressing.QueryBus
	LocationQuery      location.QueryBus
	EventBus           capi.EventBus
	ShipmentManager    *carrier.ShipmentManager

	AddressStore sqlstore.AddressStoreInterface
	OrderStore   sqlstore.OrderStoreInterface

	FlagFaboOrderUpdatePaymentSatusConfig FlagFaboOrderAutoConfirmPaymentStatus
}

type FlagFaboOrderAutoConfirmPaymentStatus bool

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
	if err := s.OrderStore.GetOrder(ctx, query); err != nil {
		return resp, err
	}
	order := query.Result.Order
	switch order.Status {
	case status5.N:
		return resp, cm.Errorf(cm.FailedPrecondition, nil, "????n h??ng ???? h???y")
	// case status5.P:
	// 	return resp, cm.Error(cm.FailedPrecondition, "????n h??ng ???? ho??n th??nh.", nil)
	case status5.NS:
		return resp, cm.Error(cm.FailedPrecondition, "????n h??ng ???? tr??? h??ng.", nil)
	}
	if order.ConfirmStatus == status3.N || order.ShopConfirm == status3.N {
		return resp, cm.Errorf(cm.FailedPrecondition, nil, "????n h??ng ???? h???y")
	}

	if err := s.RaiseOrderConfirmingEvent(ctx, shop, r.AutoInventoryVoucher, order); err != nil {
		return nil, err
	}

	if order.ConfirmStatus != status3.P ||
		order.ShopConfirm != status3.P {
		cmd := &ordermodelx.UpdateOrdersStatusCommand{
			OrderIDs:      []dot.ID{r.OrderId},
			ConfirmStatus: status3.P.Wrap(),
			ShopConfirm:   status3.P.Wrap(),
		}
		if s.FlagFaboOrderUpdatePaymentSatusConfig {
			// ????n h??ng c???a Faboshop kh??ng c?? phi???u thu chi ????? ki???m so??t payment_status
			// Tr???ng th??i ????n h??ng th?? ph??? thu???c v??o payment_status ????? ho??n th??nh
			// => Chuy???n payment_status sang Ho??n th??nh v???i c??c ????n h??ng ???????c t???o ra t??? FaboShop
			cmd.PaymentStatus = status4.P.Wrap()
		}
		if err := s.OrderStore.UpdateOrdersStatus(ctx, cmd); err != nil {
			return resp, err
		}
		event := &ordering.OrderConfirmedEvent{
			OrderID:              order.ID,
			AutoInventoryVoucher: r.AutoInventoryVoucher,
			ShopID:               shop.ID,
			InventoryOverStock:   shop.InventoryOverstock.Apply(true),
			UpdatedBy:            userID,
		}
		if err := s.EventBus.Publish(ctx, event); err != nil {
			ll.Error("RaiseOrderConfirmedEvent", l.Error(err))
		}
	}
	if err := s.OrderStore.GetOrder(ctx, query); err != nil {
		return nil, err
	}
	resp = convertpb.PbOrder(query.Result.Order, nil, account_tag.TagShop)
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
	if err := s.OrderStore.GetOrder(ctx, query); err != nil {
		return resp, err
	}
	order := query.Result.Order

	// Verify status
	switch order.Status {
	case status5.N:
		return resp, cm.Error(cm.FailedPrecondition, "????n h??ng ???? hu???.", nil)
	case status5.P:
		return resp, cm.Error(cm.FailedPrecondition, "????n h??ng ???? ho??n th??nh.", nil)
	case status5.NS:
		return resp, cm.Error(cm.FailedPrecondition, "????n h??ng ???? tr??? h??ng.", nil)
	}

	if order.ConfirmStatus == status3.N ||
		order.ShopConfirm == status3.N {
		return resp, cm.Error(cm.FailedPrecondition, "????n h??ng ???? hu???.", nil)
	}

	// Fill response
	fulfillments := query.Result.Fulfillments
	defer func() {
		if _err != nil {
			return
		}

		resp.Order = convertpb.PbOrder(order, fulfillments, account_tag.TagShop)
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
		if err := s.OrderStore.CreateFulfillments(ctx, ffmCmd); err != nil {
			return resp, err
		}
	}
	if updates != nil {
		ffmCmd := &shipmodelx.UpdateFulfillmentsCommand{
			Fulfillments: updates,
		}
		if err := s.OrderStore.UpdateFulfillments(ctx, ffmCmd); err != nil {
			return resp, err
		}
	}

	ll.S.Infof("Compare fulfillments: create %v update %v", len(creates), len(updates))
	totalChanges := len(creates) + len(updates)
	if totalChanges == 0 {
		return resp, nil
	}

	ffms := append(creates, updates...)
	if err := s.ShipmentManager.CreateFulfillments(ctx, ffms); err != nil {
		return resp, err
	}
	// automatically cancel orders on sandbox for ghn and vtpost
	if cmenv.Env() == cmenv.EnvSandbox {
		if order.ShopShipping != nil &&
			order.ShopShipping.ShippingProvider != typeshippingprovider.GHTK {
			go func() {
				time.Sleep(5 * time.Minute)
				_, err := s.CancelOrder(ctx, userID, shop.ID, partnerID, order.ID, "????n h??ng TEST, t??? ?????ng hu???", inventory_auto.Unknown)
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
		if err := s.OrderStore.UpdateOrdersStatus(ctx, cmd); err != nil {
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
	if err = s.OrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	// Get order again
	if err := s.OrderStore.GetOrder(ctx, query); err != nil {
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
	if err := s.EventBus.Publish(ctx, event); err != nil {
		return err
	}
	return nil
}

func (s *OrderLogic) prepareFulfillmentFromOrder(ctx context.Context, order *ordermodel.Order, shop *identitymodel.Shop) (*shipmodel.Fulfillment, error) {
	if order.ShopShipping != nil && order.ShopShipping.ShippingProvider == typeshippingprovider.GHN {
		if order.TryOn == 0 && order.GhnNoteCode == 0 {
			return nil, cm.Error(cm.FailedPrecondition, "Vui l??ng ch???n ghi ch?? xem h??ng!", nil)
		}
	}

	if !model.VerifyOrderSource(order.OrderSourceType) {
		return nil, cm.Error(cm.FailedPrecondition, "Kh??ng th??? x??c ?????nh ngu???n ????n h??ng!", nil)
	}
	addressTo, err := orderAddressToShippingAddress(order.ShippingAddress)
	if err != nil {
		return nil, cm.Error(cm.InvalidArgument, "Th??ng tin ?????a ch??? ng?????i nh???n: "+err.Error()+" Vui l??ng c???p nh???t v?? th??? l???i.", err)
	}
	if _, _, err := s.ShipmentManager.VerifyDistrictCode(addressTo); err != nil {
		return nil, cm.Error(cm.InvalidArgument, "Th??ng tin ?????a ch??? ng?????i nh???n: "+err.Error()+" Vui l??ng c???p nh???t v?? th??? l???i.", nil)
	}

	ffm := prepareSingleFulfillment(order, shop, order.Lines, addressTo)

	// Use shop address from order or from shop default address
	var shopAddress *addressmodel.Address
	if order.ShopShipping != nil && order.ShopShipping.ShopAddress != nil {
		shopAddress, err = orderAddressToShippingAddress(order.ShopShipping.ShopAddress)
		if err != nil {
			return nil, cm.Error(cm.InvalidArgument, "Th??ng tin ?????a ch??? c???a h??ng trong ????n h??ng: "+err.Error()+"  Vui l??ng c???p nh???t v?? th??? l???i.", err)
		}

	} else {
		if shop.ShipFromAddressID == 0 {
			return nil, cm.Error(cm.InvalidArgument, "B??n h??ng: C???n cung c???p th??ng tin ?????a ch??? l???y h??ng trong ????n h??ng ho???c t???i th??ng tin c???a h??ng. Vui l??ng c???p nh???t.", nil)
		}
		addressQuery := &addressmodelx.GetAddressQuery{AddressID: shop.ShipFromAddressID}
		if err := s.AddressStore.GetAddress(ctx, addressQuery); err != nil {
			return nil, cm.Error(cm.Internal, "L???i khi ki???m tra th??ng tin ?????a ch??? c???a c???a h??ng: "+err.Error(), err)
		}
		shopAddress = addressQuery.Result
	}
	_, _, err = s.ShipmentManager.VerifyDistrictCode(shopAddress)
	if err != nil {
		return nil, cm.Error(cm.FailedPrecondition, "Th??ng tin ?????a ch??? c???a h??ng trong c???u h??nh c???a h??ng: "+err.Error()+" Vui l??ng c???p nh???t v?? th??? l???i.", nil)
	}

	if err := s.checkBlockCarrier(shopAddress, order.ShopShipping.ShippingProvider); err != nil {
		return nil, err
	}

	ffm.TotalAmount = order.TotalAmount
	ffm.TotalDiscount = order.TotalDiscount
	ffm.AddressFrom = shopAddress

	if order.ShopCOD < 0 {
		return nil, cm.Error(cm.InvalidArgument, "Th??ng tin ti???n thu h??? (COD) kh??ng h???p l???.", nil)
	}
	if order.ShopCOD > 0 {
		ffm.TotalCODAmount = order.ShopCOD
		ffm.OriginalCODAmount = order.ShopCOD
	}

	if order.ShopShipping == nil {
		return nil, cm.Error(cm.InvalidArgument, "Vui l??ng ch???n d???ch v??? giao h??ng.", nil)
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
		return cm.Errorf(cm.InvalidArgument, nil, "%v kh??ng th??? l???y h??ng t???i ?????a ch??? n??y %v (%v)", provider.Label(), shopAddress.District, shopAddress.Province)
	}
	districts, ok := blockCarrierByDistricts[provider]
	if !ok {
		return nil
	}
	if cm.StringsContain(districts, shopAddress.DistrictCode) {
		return cm.Errorf(cm.InvalidArgument, nil, "%v kh??ng th??? l???y h??ng t???i ?????a ch??? n??y %v (%v)", provider.Label(), shopAddress.District, shopAddress.Province)
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
		TotalCODAmount:    order.ShopCOD,
		OriginalCODAmount: order.ShopCOD,

		ShippingFeeCustomer:      0, // only fill the first fulfillment
		ShippingFeeShop:          0, // after calling GHN
		ShippingFeeShopLines:     nil,
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
		ExternalShippingNote:        dot.NullString{},
		ExternalShippingSubState:    dot.NullString{},
		ExternalShippingData:        nil,
		ShippingState:               shipping.Default,
		ShippingStatus:              0,
		EtopPaymentStatus:           0,
		Status:                      0,
		SyncStatus:                  0,
		SyncStates:                  nil,
		LastSyncAt:                  time.Time{},
		ExternalShippingLogs:        nil,

		// new information
		ShippingPaymentType: shipping_payment_type.Seller,
		ShippingType:        ordertypes.ShippingTypeShipment,
		ConnectionID:        shippingcore.GetConnectionID(0, shippingProvider),
		ConnectionMethod:    connection_type.ConnectionMethodBuiltin,
	}

	if order.ShopShipping != nil {
		shopShipping := order.ShopShipping
		fulfillment.ShippingServiceFee = shopShipping.ExternalShippingFee
		fulfillment.ShippingServiceName = shopShipping.ExternalServiceName
		fulfillment.ProviderServiceID = shopShipping.ProviderServiceID
		fulfillment.IncludeInsurance = dot.Bool(shopShipping.IncludeInsurance)
		fulfillment.GrossWeight = cm.CoalesceInt(shopShipping.GrossWeight, order.TotalWeight)
		fulfillment.ChargeableWeight = cm.CoalesceInt(shopShipping.ChargeableWeight, order.TotalWeight)
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
		return nil, cm.Error(cm.InvalidArgument, "Thi???u th??ng tin ?????a ch???.", nil)
	}
	if orderAddr.Phone == "" {
		return nil, cm.Error(cm.InvalidArgument, "Thi???u th??ng tin s??? ??i???n tho???i.", nil)
	}
	if _, ok := validate.NormalizePhone(orderAddr.Phone); !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "S??? ??i???n tho???i kh??ng h???p l??? (%v).", orderAddr.Phone)
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
		go func() (_err error) {
			defer func() {
				wg.Done()
				errs[i] = _err
			}()

			var shippingProviderErr error
			if ffm.ConnectionID != 0 {
				shippingProviderErr = s.ShipmentManager.CancelFulfillment(ctx, ffm)
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
				if err := s.OrderStore.UpdateFulfillment(ctx, update2Cmd); err != nil {
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
			if err := s.OrderStore.UpdateFulfillment(ctx, update2Cmd); err != nil {
				return err
			}
			return nil
		}()
	}
	wg.Wait()
	return errs, nil
}
