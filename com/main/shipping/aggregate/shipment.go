package aggregate

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/location"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipping"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	shipstate "etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	addressconvert "etop.vn/backend/com/main/address/convert"
	orderconvert "etop.vn/backend/com/main/ordering/convert"
	"etop.vn/backend/com/main/shipping/carrier"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/com/main/shipping/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/validate"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var _ shipping.Aggregate = &Aggregate{}
var ll = l.New()

type Aggregate struct {
	db             cmsql.Transactioner
	locationQS     location.QueryBus
	orderQS        ordering.QueryBus
	shimentManager *carrier.ShipmentManager
	connectionQS   connectioning.QueryBus
	ffmStore       sqlstore.FulfillmentStoreFactory
	eventBus       capi.EventBus
}

func NewAggregate(db *cmsql.Database, locationQS location.QueryBus, orderQS ordering.QueryBus, shipmentManager *carrier.ShipmentManager, connectionQS connectioning.QueryBus, eventB capi.EventBus) *Aggregate {
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

func (a *Aggregate) MessageBus() shipping.CommandBus {
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
	if args.ShippingState == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "shipping_state không được để trống.")
	}
	query := a.ffmStore(ctx).OptionalPartnerID(args.PartnerID)
	if args.FulfillmentID != 0 {
		query = query.ID(args.FulfillmentID)
	}
	if args.ShippingCode != "" {
		query = query.ShippingCode(args.ShippingCode)
	}
	if len(args.ConnectionIDs) > 0 {
		query = query.ConnectionIDs(args.ConnectionIDs...)
	}
	ffm, err := query.GetFulfillment()
	if err != nil {
		return 0, err
	}
	event := &shipping.FulfillmentUpdatingEvent{
		EventMeta:     meta.NewEvent(),
		FulfillmentID: ffm.ID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
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

	args.FulfillmentID = ffm.ID
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
	event := &shipping.FulfillmentUpdatingEvent{
		EventMeta:         meta.NewEvent(),
		FulfillmentID:     ffm.ID,
		MoneyTxShippingID: ffm.MoneyTransactionID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return 0, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err := a.ffmStore(ctx).UpdateFulfillmentShippingFees(args); err != nil {
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
	ffm, err := a.ffmStore(ctx).ID(args.FulfillmentID).GetFfmDB()
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
		if err := a.ffmStore(ctx).CancelFulfillment(args); err != nil {
			return err
		}

		// backward compatible
		if ffm.ShippingType == 0 {
			switch ffm.ShippingProvider {
			case shipping_provider.GHN:
				ffm.ConnectionID = connectioning.DefaultTopshipGHNConnectionID
			case shipping_provider.GHTK:
				ffm.ConnectionID = connectioning.DefaultTopshipGHTKConnectionID
			case shipping_provider.VTPost:
				ffm.ConnectionID = connectioning.DefaultTopshipVTPostConnectionID
			default:

			}
		}

		// case shipment: cancel ffm from carrier
		if ffm.ConnectionID != 0 {
			return a.shimentManager.CancelFulfillment(ctx, ffm)
		}
		return nil
	})

	return err
}
