package aggregate

import (
	"context"
	"fmt"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	shipstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	com "o.o/backend/com/main"
	addressconvert "o.o/backend/com/main/address/convert"
	orderconvert "o.o/backend/com/main/ordering/convert"
	"o.o/backend/com/main/shipping/carrier"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/backend/com/main/shipping/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ shipping.Aggregate = &Aggregate{}
var ll = l.New()
var scheme = conversion.Build(shippingconvert.RegisterConversions)

type Aggregate struct {
	db             *cmsql.Database
	locationQS     location.QueryBus
	orderQS        ordering.QueryBus
	shimentManager *carrier.ShipmentManager
	connectionQS   connectioning.QueryBus
	ffmStore       sqlstore.FulfillmentStoreFactory
	eventBus       capi.EventBus
}

func NewAggregate(db com.MainDB, eventB capi.EventBus, locationQS location.QueryBus, orderQS ordering.QueryBus, shipmentManager *carrier.ShipmentManager, connectionQS connectioning.QueryBus) *Aggregate {
	return &Aggregate{
		db:             db,
		locationQS:     locationQS,
		orderQS:        orderQS,
		shimentManager: shipmentManager,
		connectionQS:   connectionQS,
		ffmStore:       sqlstore.NewFulfillmentStore(db),
		eventBus:       eventB,
	}
}

func AggregateMessageBus(a *Aggregate) shipping.CommandBus {
	b := bus.New()
	return shipping.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateFulfillments(ctx context.Context, args *shipping.CreateFulfillmentsArgs) (fulfillmentIDs []dot.ID, _ error) {
	query := &ordering.GetOrderByIDQuery{
		ID: args.OrderID,
	}
	if err := a.orderQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	order := query.Result
	switch order.Status {
	case status5.N:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã hủy")
	case status5.NS:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng đã trả hàng")
	}

	if order.ConfirmStatus != status3.P || order.ShopConfirm != status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng xác nhận đơn hàng trước khi tạo đơn giao hàng")
	}

	var ffms []*shipmodel.Fulfillment
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		event := &shipping.FulfillmentsCreatingEvent{
			EventMeta: meta.NewEvent(),
			ShopID:    args.ShopID,
			OrderID:   args.OrderID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		_, _, err := a.getAndVerifyAddress(ctx, args.ShippingAddress)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ giao hàng không hợp lệ: %v", err)
		}
		_, _, err = a.getAndVerifyAddress(ctx, args.PickupAddress)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Địa chỉ lấy hàng không hợp lệ: %v", err)
		}

		oldFulfillments, err := a.ffmStore(ctx).OrderID(args.OrderID).ListFfmsDB()
		if err != nil {
			return err
		}

		ffm, err := a.prepareFulfillmentFromOrder(ctx, order, args)
		if err != nil {
			return err
		}

		creates, updates, err := CompareFulfillments(oldFulfillments, ffm)
		if err != nil {
			return err
		}
		if creates != nil {
			if err := a.ffmStore(ctx).CreateFulfillmentsDB(creates); err != nil {
				return err
			}
		}
		if updates != nil {
			if err := a.ffmStore(ctx).StatusNotIn(status5.N, status5.NS, status5.P).UpdateFulfillmentsDB(updates); err != nil {
				return err
			}
		}
		ll.S.Infof("Compare fulfillments: create %v update %v", len(creates), len(updates))
		totalChanges := len(creates) + len(updates)
		if totalChanges == 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "Không tạo được fulfillment")
		}

		ffms = append(creates, updates...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := a.shimentManager.CreateFulfillments(ctx, order, ffms); err != nil {
		return nil, err
	}

	var res []dot.ID
	for _, ffm := range ffms {
		res = append(res, ffm.ID)
	}

	// TODO: kiểm tra trường hợp tự giao hoặc giao qua NVC
	// rollback khi gặp lỗi. VD: hủy ffm bên NVC
	event2 := &shipping.FulfillmentsCreatedEvent{
		EventMeta:      meta.NewEvent(),
		FulfillmentIDs: res,
		ShippingType:   args.ShippingType,
		OrderID:        order.ID,
	}
	if err := a.eventBus.Publish(ctx, event2); err != nil {
		return nil, err
	}
	return res, nil
}

func (a *Aggregate) getAndVerifyAddress(ctx context.Context, address *ordertypes.Address) (*location.Province, *location.District, error) {
	if address == nil {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin địa chỉ")
	}
	if address.Phone == "" {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin số điện thoại")
	}
	if _, ok := validate.NormalizePhone(address.Phone); !ok {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ (%v).", address.Phone)
	}

	query := &location.FindOrGetLocationQuery{
		Province:     address.Province,
		District:     address.District,
		Ward:         address.Ward,
		ProvinceCode: address.ProvinceCode,
		DistrictCode: address.DistrictCode,
		WardCode:     address.WardCode,
	}
	if err := a.locationQS.Dispatch(ctx, query); err != nil {
		return nil, nil, cm.Errorf(cm.InvalidArgument, err, "Địa chỉ không hợp lệ %v", err.Error())
	}
	district, province := query.Result.District, query.Result.Province

	address.Location = ordertypes.Location{
		ProvinceCode: province.Code,
		Province:     province.Name,
		DistrictCode: district.Code,
		District:     district.Name,
		Coordinates:  address.Coordinates,
	}
	ward := query.Result.Ward
	if ward != nil {
		address.Location.Ward = ward.Name
		address.Location.WardCode = ward.Code
	}
	return province, district, nil
}

func (a *Aggregate) prepareFulfillmentFromOrder(ctx context.Context, order *ordering.Order, args *shipping.CreateFulfillmentsArgs) (*shipmodel.Fulfillment, error) {
	shippingType := args.ShippingType
	var connectionMethod connection_type.ConnectionMethod
	var conn *connectioning.Connection

	switch shippingType {
	case ordertypes.ShippingTypeManual:
		if args.ShopCarrierID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển tự quản lý")
		}
	case ordertypes.ShippingTypeShipment:
		if args.ConnectionID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển (connection_id)")
		}
		queryConn := &connectioning.GetConnectionByIDQuery{
			ID: args.ConnectionID,
		}
		if err := a.connectionQS.Dispatch(ctx, queryConn); err != nil {
			return nil, err
		}
		conn = queryConn.Result

		if conn.ConnectionProvider == connection_type.ConnectionProviderGHN {
			if args.TryOn.String() == "" {
				return nil, cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng chọn ghi chú xem hàng!")
			}
		}
		connectionMethod = conn.ConnectionMethod
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phương thức vận chuyển không hợp lệ. Chỉ hỗ trợ: manual và shipment.")
	}

	tryOn := args.TryOn
	if tryOn == 0 {
		tryOn = try_on.None
	}

	variantIDs := []dot.ID{}
	totalItems, basketValue, totalAmount := 0, 0, 0

	if len(order.Lines) != 0 {
		for _, line := range order.Lines {
			variantIDs = append(variantIDs, line.VariantID)
		}
	}
	totalItems = order.TotalItems
	basketValue = order.BasketValue
	totalAmount = order.TotalAmount
	variantIDs = variantIDs

	typeFrom := etopmodel.FFShop
	typeTo := etopmodel.FFCustomer

	ffm := &shipmodel.Fulfillment{
		ID:                cm.NewID(),
		OrderID:           order.ID,
		ShopID:            order.ShopID,
		PartnerID:         order.PartnerID,
		ShopConfirm:       status3.P, // Always set shop_confirm to 1
		ConfirmStatus:     0,
		TotalItems:        totalItems,
		TotalWeight:       cm.CoalesceInt(args.GrossWeight, args.ChargeableWeight, order.TotalWeight),
		BasketValue:       basketValue,
		TotalDiscount:     0,
		TotalAmount:       totalAmount,
		TotalCODAmount:    args.CODAmount,
		OriginalCODAmount: order.ShopCOD,
		EtopDiscount:      0,
		EtopFeeAdjustment: 0,
		VariantIDs:        variantIDs,
		Lines:             orderconvert.OrderLinesToModel(order.Lines),
		// after this
		TypeFrom:            typeFrom,
		TypeTo:              typeTo,
		AddressFrom:         addressconvert.OrderAddressToModel(args.PickupAddress),
		AddressTo:           addressconvert.OrderAddressToModel(args.ShippingAddress),
		AddressReturn:       addressconvert.OrderAddressToModel(args.ReturnAddress),
		ShippingFeeCustomer: order.ShopShippingFee,
		ProviderServiceID:   args.ShippingServiceCode,
		ShippingServiceFee:  args.ShippingServiceFee,
		ShippingServiceName: args.ShippingServiceName,
		ShippingNote:        args.ShippingNote,
		TryOn:               tryOn,
		IncludeInsurance:    args.IncludeInsurance,
		ConnectionID:        args.ConnectionID,
		ConnectionMethod:    connectionMethod,
		ShopCarrierID:       args.ShopCarrierID,
		GrossWeight:         args.GrossWeight,
		ChargeableWeight:    args.ChargeableWeight,
		Width:               args.Width,
		Height:              args.Height,
		Length:              args.Length,
		ShippingState:       shipstate.Default,
		ShippingType:        args.ShippingType,
	}
	if conn != nil {
		// backward compatible
		if shippingProvider, ok := shipping_provider.ParseShippingProvider(conn.ConnectionProvider.String()); ok {
			ffm.ShippingProvider = shippingProvider
		}
	}
	return ffm, nil
}

// Compare the old fulfillments and expected fulfillments
// - Missing/cancelled fulfillments: create
// - Error fulfillments: update
// - Processing fulfillments: ignore
func CompareFulfillments(olds []*shipmodel.Fulfillment, ffm *shipmodel.Fulfillment) (creates, updates []*shipmodel.Fulfillment, err error) {
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

func (a *Aggregate) UpdateFulfillmentShippingState(ctx context.Context, args *shipping.UpdateFulfillmentShippingStateArgs) (updated int, _ error) {
	if args.FulfillmentID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Fulfillment ID không được để trống.")
	}
	if args.ShippingState == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "shipping_state không được để trống.")
	}
	query := a.ffmStore(ctx).OptionalPartnerID(args.PartnerID).ID(args.FulfillmentID)

	ffm, err := query.GetFulfillment()
	if err != nil {
		return 0, err
	}
	if ffm.MoneyTransactionID != 0 {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật trạng thái giao hàng. Đơn đã nằm trong phiên đối soát.")
	}
	if args.ActualCompensationAmount.Valid {
		if ffm.ShippingState != shipstate.Undeliverable && args.ShippingState != shipstate.Undeliverable {
			return 0, cm.Errorf(cm.FailedPrecondition, nil, "Chỉ cập nhật phí hoàn hàng khi đơn vận chuyển không giao được hàng")
		}
	}

	if err := a.ffmStore(ctx).UpdateFulfillmentShippingState(args, ffm.TotalCODAmount); err != nil {
		return 0, nil
	}
	return 1, nil
}

func (a *Aggregate) UpdateFulfillmentShippingFees(ctx context.Context, args *shipping.UpdateFulfillmentShippingFeesArgs) (updated int, _ error) {
	if args.FulfillmentID == 0 && args.ShippingCode == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	ffm, err := a.ffmStore(ctx).OptionalID(args.FulfillmentID).OptionalShippingCode(args.ShippingCode).GetFulfillment()
	if err != nil {
		return 0, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var lines []*shippingsharemodel.ShippingFeeLine
		var providerShippingFeeLines []*shippingsharemodel.ShippingFeeLine
		if err := scheme.Convert(args.ProviderShippingFeeLines, &providerShippingFeeLines); err != nil {
			return err
		}
		if args.ShippingFeeLines != nil {
			if err := scheme.Convert(args.ShippingFeeLines, &lines); err != nil {
				return err
			}
		} else {
			lines = shippingsharemodel.GetShippingFeeShopLines(providerShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		}

		totalShippingFeeShop := shippingsharemodel.GetTotalShippingFee(lines)
		update := &shipmodel.Fulfillment{
			ProviderShippingFeeLines: providerShippingFeeLines,
			ShippingFeeShopLines:     lines,
			ShippingFeeShop:          shipping.CalcShopShippingFee(totalShippingFeeShop, ffm),
		}
		if err := a.ffmStore(ctx).ID(args.FulfillmentID).UpdateFulfillmentDB(update); err != nil {
			return err
		}
		if err := a.ffmStore(ctx).UpdateFulfillmentPriceListInfo(args); err != nil {
			return err
		}

		eventChanged := &shipping.FulfillmentShippingFeeChangedEvent{
			EventMeta:         meta.NewEvent(),
			FulfillmentID:     ffm.ID,
			MoneyTxShippingID: ffm.MoneyTransactionID,
		}
		if err := a.eventBus.Publish(ctx, eventChanged); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *Aggregate) UpdateFulfillmentsMoneyTxID(ctx context.Context, args *shipping.UpdateFulfillmentsMoneyTxIDArgs) (updated int, _ error) {
	if len(args.FulfillmentIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillmentIDs").WithMetap("function", "UpdateFulfillmentsMoneyTxID")
	}

	if args.MoneyTxShippingExternalID == 0 && args.MoneyTxShippingID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing MoneyTxID").WithMetap("function", "UpdateFulfillmentsMoneyTxID")
	}

	return a.ffmStore(ctx).UpdateFulfillmentsMoneyTxID(args)
}

func (a *Aggregate) UpdateFulfillmentsCODTransferedAt(ctx context.Context, args *shipping.UpdateFulfillmentsCODTransferedAtArgs) error {
	if len(args.MoneyTxShippingIDs) == 0 &&
		len(args.FulfillmentIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing required fields").WithMetap("func", "UpdateFulfillmentsCODTransferedAt")
	}
	if args.CODTransferedAt.IsZero() {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing cod_transfered_at").WithMetap("func", "UpdateFulfillmentsCODTransferedAt")
	}
	query := a.ffmStore(ctx)
	if len(args.MoneyTxShippingIDs) > 0 {
		query = query.MoneyTxShippingIDs(args.MoneyTxShippingIDs...)
	}
	if len(args.FulfillmentIDs) > 0 {
		query = query.IDs(args.FulfillmentIDs...)
	}
	update := &shipmodel.Fulfillment{
		CODEtopTransferedAt: args.CODTransferedAt,
	}
	return query.UpdateFulfillmentDB(update)
}

func (a *Aggregate) RemoveFulfillmentsMoneyTxID(ctx context.Context, args *shipping.RemoveFulfillmentsMoneyTxIDArgs) (updated int, _ error) {
	if len(args.FulfillmentIDs) == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillmentIDs").WithMetap("function", "RemoveFulfillmentsMoneyTxShippingExternalID")
	}
	return a.ffmStore(ctx).RemoveFulfillmentsMoneyTxID(args)
}

func (a *Aggregate) UpdateFulfillmentsStatus(ctx context.Context, args *shipping.UpdateFulfillmentsStatusArgs) error {
	if len(args.FulfillmentIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing FulfillmentIDs").WithMetap("function", "UpdateFulfillmentsStatus")
	}

	return a.ffmStore(ctx).UpdateFulfillmentsStatus(args)
}

func (a *Aggregate) CancelFulfillment(ctx context.Context, args *shipping.CancelFulfillmentArgs) error {
	ffm, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFulfillment()
	if err != nil {
		return err
	}
	switch ffm.Status {
	case status5.P, status5.NS:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hoàn thành. Không thể hủy")
	case status5.N:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy.")
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// backward compatible
		ffm.ConnectionID = shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)

		var ffmDB shipmodel.Fulfillment
		if err := scheme.Convert(ffm, &ffmDB); err != nil {
			return err
		}
		// case shipment: cancel ffm from carrier
		if ffm.ConnectionID != 0 {
			if err := a.shimentManager.CancelFulfillment(ctx, &ffmDB); err != nil {
				return err
			}
		}
		if err := a.ffmStore(ctx).CancelFulfillment(args); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (a *Aggregate) UpdateFulfillmentExternalShippingInfo(ctx context.Context, args *shipping.UpdateFfmExternalShippingInfoArgs) (updated int, err error) {
	if args.FulfillmentID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing Fulfillment ID")
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		update := &shipmodel.Fulfillment{
			ShippingState:             args.ShippingState,
			ShippingStatus:            args.ShippingStatus,
			ExternalShippingData:      args.ExternalShippingData,
			ExternalShippingState:     args.ExternalShippingState,
			ExternalShippingStatus:    args.ExternalShippingStatus,
			ExternalShippingUpdatedAt: args.ExternalShippingUpdatedAt,
			ExternalShippingStateCode: args.ExternalShippingStateCode,
			ClosedAt:                  args.ClosedAt,
			LastSyncAt:                args.LastSyncAt,
			ShippingCreatedAt:         args.ShippingCreatedAt,
			ShippingPickingAt:         args.ShippingPickingAt,
			ShippingDeliveringAt:      args.ShippingDeliveringAt,
			ShippingDeliveredAt:       args.ShippingDeliveredAt,
			ShippingReturningAt:       args.ShippingReturningAt,
			ShippingReturnedAt:        args.ShippingReturnedAt,
			ShippingCancelledAt:       args.ShippingCancelledAt,
		}
		if args.Weight != 0 {
			update.TotalWeight = args.Weight
			update.GrossWeight = args.Weight
		}
		if args.ExternalShippingLogs != nil {
			logs := shippingconvert.Convert_shipping_ExternalShippingLogs_shippingmodel_ExternalShippingLogs(args.ExternalShippingLogs)
			update.ExternalShippingLogs = logs
		}

		if err := a.ffmStore(ctx).ID(args.FulfillmentID).UpdateFulfillmentDB(update); err != nil {
			return err
		}

		args2 := &sqlstore.ForceUpdateExternalShippingInfoArgs{
			FulfillmentID:            args.FulfillmentID,
			ExternalShippingNote:     args.ExternalShippingNote,
			ExternalShippingSubState: args.ExternalShippingSubState,
		}
		if err := a.ffmStore(ctx).ForceUpdateExternalShippingInfo(args2); err != nil {
			return err
		}

		updated = 1
		return nil
	})
	return
}

// UpdateFulfillmentShippingFeesFromWebhook
//
// Cập nhật giá vận chuyển từ webhook
// Luôn cập nhật ProviderShippingFeeLines
// Kiểm tra nếu có thay đổi về khối lượng, sẽ tính lại giá theo bảng giá hiện tại của TOPSHIP. Nếu không có giá TOPSHIP sẽ giữ nguyên giá từ NVC => cập nhật lại giá mới hoặc thông báo qua telegram nếu đơn đã nằm trong phiên thanh toán với shop
func (a *Aggregate) UpdateFulfillmentShippingFeesFromWebhook(ctx context.Context, args *shipping.UpdateFulfillmentShippingFeesFromWebhookArgs) (_err error) {
	providerFeeLines := args.ProviderFeeLines
	if providerFeeLines == nil || len(providerFeeLines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing providerFeeLines")
	}
	ffm, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFulfillment()
	if err != nil {
		return err
	}

	update := &shipping.Fulfillment{
		ID:                       ffm.ID,
		ProviderShippingFeeLines: args.ProviderFeeLines,
	}
	defer func() error {
		// always update shipping fee even if error occurred
		if update.ShippingFeeShopLines != nil {
			totalFee := shipping.GetTotalShippingFee(update.ShippingFeeShopLines)
			shippingFeeShop := shipping.CalcShopShippingFee(totalFee, ffm)
			update.ShippingFeeShop = shippingFeeShop
			if shippingFeeShop != ffm.ShippingFeeShop {
				// Giá thay đổi
				// check money_transaction_shipping
				if ffm.MoneyTransactionID != 0 {
					// Đơn đã nằm trong phiên
					// Giá cước đơn thay đổi
					// Không cập nhật + bắn noti telegram để follow
					connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
					if connectionID == 0 {
						return cm.Errorf(cm.FailedPrecondition, nil, "ConnectionID can not be empty")
					}
					connection, _ := a.shimentManager.GetConnectionByID(ctx, connectionID)
					str := "–––\n👹 %v: đơn %v có thay đổi về giá nhưng đã nằm trong phiên thanh toán. Không thể cập nhật, vui lòng kiểm tra lại. 👹 \n- Giá hiện tại: %v \n- Giá mới: %v\n–––"
					ll.SendMessage(fmt.Sprintf(str, connection.Name, ffm.ShippingCode, ffm.ShippingFeeShop, shippingFeeShop))
					// shop shipping fee does not change
					update.ShippingFeeShopLines = nil
					update.ShippingFeeShop = 0
				}
			}
		}

		if err := a.ffmStore(ctx).ID(ffm.ID).UpdateFulfillment(update); err != nil {
			return err
		}
		// keep origin error
		return _err
	}()

	if !a.shimentManager.FlagApplyShipmentPrice || !ffm.EtopPriceRule {
		update.ShippingFeeShopLines = shipping.GetShippingFeeShopLines(providerFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		return nil
	}

	// Trường hợp có áp dụng bảng giá
	// Các trường hợp cần tính lại cước phí đơn
	//   - Thay đổi khối lượng
	//   - Đơn trả hàng => tính phí trả hàng
	var feeLines []*shipping.ShippingFeeLine
	shippingFeeShopLines := ffm.ShippingFeeShopLines
	if args.NewWeight != ffm.TotalWeight {
		feeLines, err = a.shimentManager.CalcMakeupShippingFeesByFfm(ctx, ffm, args.NewWeight, args.NewState)
		if err != nil {
			return err
		}
		mainFeeLine := shipping.GetShippingFeeLine(feeLines, shipping_fee_type.Main)
		shippingFeeShopLines = shipping.ApplyShippingFeeLine(shippingFeeShopLines, mainFeeLine)

		// Remove field EtopAdjustedShippingFeeMain if not use
		update.ShippingFeeShopLines = shippingFeeShopLines
		update.EtopAdjustedShippingFeeMain = mainFeeLine.Cost
	}

	if shipping.IsStateReturn(args.NewState) {
		if feeLines == nil {
			feeLines, err = a.shimentManager.CalcMakeupShippingFeesByFfm(ctx, ffm, args.NewWeight, args.NewState)
			if err != nil {
				return err
			}
		}
		returnFeeLine := shipping.GetShippingFeeLine(feeLines, shipping_fee_type.Return)
		shippingFeeShopLines = shipping.ApplyShippingFeeLine(shippingFeeShopLines, returnFeeLine)
		update.ShippingFeeShopLines = shippingFeeShopLines
	}

	return nil
}
