package shipnow

import (
	"context"
	"strings"

	"o.o/api/main/address"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	com "o.o/backend/com/main"
	shipnowcarrier "o.o/backend/com/main/shipnow/carrier"
	"o.o/backend/com/main/shipnow/convert"
	"o.o/backend/com/main/shipnow/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

var _ shipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	location        location.QueryBus
	identityQuery   identity.QueryBus
	addressQuery    address.QueryBus
	order           ordering.QueryBus
	connectionQuery connectioning.QueryBus

	db             *cmsql.Database
	store          sqlstore.ShipnowStoreFactory
	eventBus       capi.EventBus
	shipnowManager *shipnowcarrier.ShipnowManager
}

func NewAggregate(eventBus capi.EventBus,
	db com.MainDB, location location.QueryBus,
	identityQuery identity.QueryBus,
	addressQuery address.QueryBus,
	connectionQS connectioning.QueryBus,
	order ordering.QueryBus,
	shipnowManager *shipnowcarrier.ShipnowManager,
) *Aggregate {
	return &Aggregate{
		db:       db,
		store:    sqlstore.NewShipnowStore(db),
		eventBus: eventBus,

		location:        location,
		identityQuery:   identityQuery,
		addressQuery:    addressQuery,
		connectionQuery: connectionQS,
		order:           order,
		shipnowManager:  shipnowManager,
	}
}

func AggregateMessageBus(a *Aggregate) shipnow.CommandBus {
	b := bus.New()
	return shipnow.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, args *shipnow.CreateShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	orderIDs := make([]dot.ID, len(args.DeliveryPoints))
	for i, point := range args.DeliveryPoints {
		if cm.IDsContain(orderIDs, point.OrderID) {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Một đơn hàng không được chọn nhiều lần.")
		}
		orderIDs[i] = point.OrderID
	}

	var conn *connectioning.Connection
	if args.ConnectionID == 0 {
		if args.Carrier == 0 {
		}
		switch args.Carrier {
		case carriertypes.Default:
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển")
		case carriertypes.Ahamove:
			args.ConnectionID = connectioning.DefaultTopShipAhamoveConnectionID
		default:
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Nhà vận chuyển không hợp lệ")
		}
	}
	queryConn := &connectioning.GetConnectionByIDQuery{
		ID: args.ConnectionID,
	}
	if err := a.connectionQuery.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}
	conn = queryConn.Result

	if args.ExternalID != "" && !validate.ExternalCode(args.ExternalID) {
		return nil, cm.Error(cm.InvalidArgument, "Mã đơn external_id không hợp lệ", nil)
	}

	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffmID := cm.NewID()
		// ShipnowOrderReservationEvent
		event := &shipnow.ShipnowOrderReservationEvent{
			EventMeta:            meta.NewEvent(),
			OrderIDs:             orderIDs,
			ShipnowFulfillmentID: ffmID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		pickupAddress, err := a.PreparePickupAddress(ctx, args.ShopID, args.PickupAddress)
		if err != nil {
			return err
		}
		points, weightInfo, valueInfo, err := a.PrepareDeliveryPointsV2(ctx, args.DeliveryPoints)
		if err != nil {
			return err
		}
		shipnowFfm := &shipnow.ShipnowFulfillment{
			ID:                  ffmID,
			ShopID:              args.ShopID,
			PickupAddress:       pickupAddress,
			DeliveryPoints:      points,
			Carrier:             args.Carrier,
			ShippingServiceCode: args.ShippingServiceCode,
			ShippingServiceFee:  args.ShippingServiceFee,
			WeightInfo:          weightInfo,
			ValueInfo:           valueInfo,
			ShippingNote:        args.ShippingNote,
			ConnectionID:        conn.ID,
			ConnectionMethod:    conn.ConnectionMethod,
			ExternalID:          args.ExternalID,
			Coupon:              args.Coupon,
		}

		if err := a.store(ctx).Create(shipnowFfm); err != nil {
			if xerr, ok := err.(*xerrors.APIError); ok && xerr.Err != nil {
				msg := xerr.Err.Error()
				switch {
				case strings.Contains(msg, "shipnow_fulfillment_partner_external_id_idx"), strings.Contains(msg, "shipnow_fulfillment_shop_external_id_idx"):
					newErr := cm.Errorf(cm.AlreadyExists, nil, "Mã đơn external_id đã tồn tại. Vui lòng kiểm tra lại").WithMeta("duplicated", "external_id")
					return newErr
				}
			}
			return err
		}
		_result = shipnowFfm
		return nil
	})
	return _result, err
}

func (a *Aggregate) UpdateShipnowFulfillment(ctx context.Context, args *shipnow.UpdateShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetShipnow()
		if err != nil {
			return err
		}
		if ffm.ConfirmStatus != status3.Z || ffm.ShippingCode != "" {
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật đơn giao hàng này.")
		}

		orderIDs := make([]dot.ID, len(args.DeliveryPoints))
		for i, point := range args.DeliveryPoints {
			if cm.IDsContain(orderIDs, point.OrderID) {
				return cm.Errorf(cm.InvalidArgument, nil, "Một đơn hàng không được chọn nhiều lần.")
			}
			orderIDs[i] = point.OrderID
		}

		updateArgs := sqlstore.UpdateInfoArgs{
			ID:                  args.ID,
			PickupAddress:       args.PickupAddress,
			Carrier:             args.Carrier,
			ShippingServiceCode: args.ShippingServiceCode,
			ShippingServiceFee:  args.ShippingServiceFee,
			ShippingNote:        args.ShippingNote,
			RequestPickupAt:     args.RequestPickupAt,
			Coupon:              args.Coupon,
		}

		if len(orderIDs) > 0 {
			// ShipnowOrderChangedEvent
			event := &shipnow.ShipnowOrderChangedEvent{
				EventMeta:            meta.NewEvent(),
				ShipnowFulfillmentID: ffm.ID,
				OldOrderIDs:          ffm.OrderIDs,
				OrderIDs:             orderIDs,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return nil
			}
			points, weightInfo, valueInfo, err := a.PrepareDeliveryPointsV2(ctx, args.DeliveryPoints)
			if err != nil {
				return err
			}
			updateArgs.DeliveryPoints = points
			updateArgs.WeightInfo = weightInfo
			updateArgs.ValueInfo = valueInfo
		}

		result, err := a.store(ctx).UpdateInfo(updateArgs)
		if err != nil {
			return err
		}
		_result = result
		return nil
	})
	return _result, err
}

func (a *Aggregate) CancelShipnowFulfillment(ctx context.Context, cmd *shipnow.CancelShipnowFulfillmentArgs) (*meta.Empty, error) {
	if cmd.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu shop_id")
	}
	if cmd.ID == 0 && cmd.ShippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp id hoặc shipping_code")
	}
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm, err := a.store(ctx).ShopID(cmd.ShopID).
			OptionalID(cmd.ID).
			OptionalShippingCode(cmd.ShippingCode).
			GetShipnow()
		if err != nil {
			return err
		}

		switch ffm.Status {
		case status5.P, status5.N, status5.NS:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển không thể hủy")
		}

		switch ffm.ShippingState {
		case shipnow_state.StateCancelled:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã bị hủy")
		case shipnow_state.StateDelivering:
			return cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đang giao. Không thể hủy đơn.")
		case shipnow_state.StateDelivered,
			shipnow_state.StateReturning, shipnow_state.StateReturned:
			return cm.Errorf(cm.FailedPrecondition, nil, "Không thể hủy đơn.")
		}

		event := &shipnow.ShipnowCancelledEvent{
			EventMeta:            meta.NewEvent(),
			ShipnowFulfillmentID: ffm.ID,
			OrderIDs:             ffm.OrderIDs,
			CancelReason:         cmd.CancelReason,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		updateArgs := sqlstore.UpdateCancelArgs{
			ID:            ffm.ID,
			ShippingState: shipnow_state.StateCancelled,
			Status:        status5.N,
			ConfirmStatus: status3.N,
			CancelReason:  cmd.CancelReason,
		}
		ffm, err = a.store(ctx).UpdateCancelled(updateArgs)
		if err != nil {
			return err
		}
		return nil
	})
	return &meta.Empty{}, err
}

func (a *Aggregate) ConfirmShipnowFulfillment(ctx context.Context, cmd *shipnow.ConfirmShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _err error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffm, err := a.store(ctx).ID(cmd.ID).ShopID(cmd.ShopID).GetShipnow()
		if err != nil {
			return err
		}
		if err := ValidateConfirmFulfillment(ffm); err != nil {
			return err
		}

		event := &shipnow.ShipnowValidateConfirmedEvent{
			EventMeta:            meta.NewEvent(),
			ShipnowFulfillmentID: ffm.ID,
			OrderIDs:             ffm.OrderIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		event2 := &shipnow.ShipnowExternalCreatedEvent{
			ShipnowFulfillmentID: ffm.ID,
		}
		if err := a.eventBus.Publish(ctx, event2); err != nil {
			return err
		}

		event3 := &shipnow.ShipnowCreatedEvent{
			ShipnowFulfillmentID: ffm.ID,
		}
		if err := a.eventBus.Publish(ctx, event3); err != nil {
			return err
		}

		update := sqlstore.UpdateStateArgs{
			ID:             cmd.ID,
			ConfirmStatus:  status3.P,
			ShippingStatus: status5.S,
			Status:         status5.S,
		}
		shipnowFfm, err := a.store(ctx).UpdateSyncState(update)
		if err != nil {
			return err
		}
		_result = shipnowFfm
		return nil
	})
	return _result, err
}

func (a *Aggregate) PreparePickupAddress(ctx context.Context, shopID dot.ID, pickupAddress *ordertypes.Address) (*ordertypes.Address, error) {
	if pickupAddress != nil {
		return a.ValidateAddress(ctx, pickupAddress)
	}
	query := &identity.GetShopByIDQuery{ID: shopID}
	if err := a.identityQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	shop := query.Result
	shopAddressID := shop.ShipFromAddressID
	if shopAddressID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
	}
	query2 := &address.GetAddressByIDQuery{
		ID: shopAddressID,
	}
	if err := a.addressQuery.Dispatch(ctx, query2); err != nil {
		return nil, err
	}
	shopAddress := query2.Result
	pickupAddress = shopAddress.ToOrderAddress()
	return pickupAddress, nil
}

func (a *Aggregate) PrepareDeliveryPoints(ctx context.Context, orderIDs []dot.ID) (points []*shipnow.DeliveryPoint, weightInfo shippingtypes.WeightInfo, valueinfo shippingtypes.ValueInfo, _err error) {
	query := &ordering.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := a.order.Dispatch(ctx, query); err != nil {
		_err = err
		return
	}
	orders := query.Result.Orders

	// Note: Không thay đổi thứ tự đơn hàng vì nó ảnh hưởng tới giá
	mapOrders := make(map[dot.ID]*ordering.Order)
	for _, order := range orders {
		if order.Shipping == nil {
			_err = cm.Errorf(cm.InvalidArgument, nil, "Đơn hàng thiếu thông tin giao hàng: khối lượng, COD,...")
			return
		}
		mapOrders[order.ID] = order
	}
	for _, orderID := range orderIDs {
		order := mapOrders[orderID]
		points = append(points, convert.OrderToDeliveryPoint(order))
	}
	weightInfo = convert.GetWeightInfo(orders)
	valueinfo = convert.GetValueInfo(orders)
	return
}

func (a *Aggregate) PrepareDeliveryPointsV2(ctx context.Context, points []*shipnow.OrderShippingInfo) (deliveryPoints []*shipnow.DeliveryPoint, weightInfo shippingtypes.WeightInfo, valueInfo shippingtypes.ValueInfo, _err error) {
	orderIDs := make([]dot.ID, len(points))
	for i, point := range points {
		orderIDs[i] = point.OrderID
	}
	query := &ordering.GetOrdersQuery{
		IDs: orderIDs,
	}
	if err := a.order.Dispatch(ctx, query); err != nil {
		_err = err
		return
	}
	orders := query.Result.Orders

	// Note: Không thay đổi thứ tự đơn hàng vì nó ảnh hưởng tới giá
	mapOrders := make(map[dot.ID]*ordering.Order)
	for _, order := range orders {
		mapOrders[order.ID] = order
	}

	for _, point := range points {
		if point.ShippingAddress == nil {
			_err = cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp địa chỉ giao hàng")
		}
		point.ShippingAddress, _err = a.ValidateAddress(ctx, point.ShippingAddress)
		if _err != nil {
			return
		}
		order := mapOrders[point.OrderID]
		p := &shipnowtypes.DeliveryPoint{
			ShippingAddress: point.ShippingAddress,
			Lines:           order.Lines,
			ShippingNote:    cm.Coalesce(point.ShippingNote, order.ShippingNote),
			OrderId:         order.ID,
			OrderCode:       order.Code,
			WeightInfo: shippingtypes.WeightInfo{
				GrossWeight:      point.GrossWeight,
				ChargeableWeight: point.ChargeableWeight,
				Length:           point.Length,
				Width:            point.Width,
				Height:           point.Height,
			},
			ValueInfo: shippingtypes.ValueInfo{
				BasketValue:      cm.CoalesceInt(point.BasketValue, order.BasketValue),
				CODAmount:        point.CODAmount,
				IncludeInsurance: point.IncludeInsurance,
			},
			TryOn: point.TryOn,
		}

		deliveryPoints = append(deliveryPoints, p)
		weightInfo.ChargeableWeight += point.ChargeableWeight
		weightInfo.GrossWeight += point.GrossWeight
		valueInfo.BasketValue += cm.CoalesceInt(point.BasketValue, order.BasketValue)
		valueInfo.CODAmount += point.CODAmount
	}
	return
}

func ValidateConfirmFulfillment(ffm *shipnow.ShipnowFulfillment) error {
	switch ffm.ConfirmStatus {
	case status3.N:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy")
	case status3.P:
		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã xác nhận")
	}
	if ffm.Status == status5.N || ffm.Status == status5.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể xác nhận đơn giao hàng này")
	}

	if len(ffm.DeliveryPoints) == 0 || len(ffm.OrderIDs) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số điểm giao hàng không hợp lệ")
	}
	return nil
}

func (a *Aggregate) UpdateShipnowFulfillmentCarrierInfo(ctx context.Context, args *shipnow.UpdateShipnowFulfillmentCarrierInfoArgs) (*shipnow.ShipnowFulfillment, error) {
	updateArgs := sqlstore.UpdateCarrierInfoArgs{
		ID:                  args.ID,
		FeeLines:            args.FeeLines,
		CarrierFeeLines:     args.CarrierFeeLines,
		ShippingCode:        args.ShippingCode,
		ShippingCreatedAt:   args.ShippingCreatedAt,
		ShippingState:       args.ShippingState,
		ShippingStatus:      args.ShippingStatus,
		EtopPaymentStatus:   args.EtopPaymentStatus,
		CODEtopTransferedAt: args.CodEtopTransferedAt,
		Status:              args.Status,

		ShippingPickingAt:          args.ShippingPickingAt,
		ShippingDeliveringAt:       args.ShippingDeliveringAt,
		ShippingDeliveredAt:        args.ShippingDeliveredAt,
		ShippingCancelledAt:        args.ShippingCancelledAt,
		ShippingServiceName:        args.ShippingServiceName,
		ShippingServiceDescription: args.ShippingServiceDescription,
		CancelReason:               args.CancelReason,
		ShippingSharedLink:         args.ShippingSharedLink,
		DeliveryPoints:             args.DeliveryPoints,
	}
	updateArgs.TotalFee = shippingtypes.GetTotalShippingFee(args.FeeLines)
	ffm, err := a.store(ctx).UpdateCarrierInfo(updateArgs)
	return ffm, err
}

func (a *Aggregate) UpdateShipnowFulfillmentState(ctx context.Context, args *shipnow.UpdateShipnowFulfillmentStateArgs) (*shipnow.ShipnowFulfillment, error) {
	updateArgs := sqlstore.UpdateStateArgs{
		ID:             args.Id,
		SyncStatus:     args.SyncStatus,
		Status:         args.Status,
		ShippingState:  args.ShippingState,
		SyncStates:     args.SyncStates,
		ConfirmStatus:  args.ConfirmStatus,
		ShippingStatus: args.ShippingStatus,
	}
	ffm, err := a.store(ctx).UpdateSyncState(updateArgs)
	return ffm, err
}

func (a *Aggregate) GetShipnowServices(ctx context.Context, args *shipnow.GetShipnowServicesArgs) (*shipnow.GetShipnowServicesResult, error) {
	if len(args.OrderIds) == 0 && len(args.DeliveryPoints) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp địa chỉ giao hàng")
	}

	pickupAddress, err := a.PreparePickupAddress(ctx, args.ShopId, args.PickupAddress)
	if err != nil {
		return nil, cm.Errorf(cm.ErrorCode(err), err, "Địa chỉ lấy hàng không hợp lệ: %v", err.Error())
	}

	var points = args.DeliveryPoints
	if len(args.OrderIds) > 0 {
		points, _, _, err = a.PrepareDeliveryPoints(ctx, args.OrderIds)
		if err != nil {
			return nil, err
		}
	}
	for _, p := range points {
		p.ShippingAddress, err = a.ValidateAddress(ctx, p.ShippingAddress)
		if err != nil {
			return nil, cm.Errorf(cm.ErrorCode(err), err, "Địa chỉ giao hàng không hợp lệ: %v", err.Error())
		}
	}

	cmd := &carrier.GetExternalShipnowServicesCommand{
		ShopID:         args.ShopId,
		PickupAddress:  pickupAddress,
		DeliveryPoints: points,
		ConnectionIDs:  args.ConnectionIDs,
		Coupon:         args.Coupon,
	}

	services, err := a.shipnowManager.GetExternalShipnowServices(ctx, cmd)
	return &shipnow.GetShipnowServicesResult{
		Services: services,
	}, nil
}

func (a *Aggregate) ValidateAddress(ctx context.Context, addr *ordertypes.Address) (*ordertypes.Address, error) {
	if addr == nil {
		return nil, nil
	}
	locationQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: addr.ProvinceCode,
		DistrictCode: addr.DistrictCode,
		WardCode:     addr.WardCode,
		Province:     addr.Province,
		District:     addr.District,
		Ward:         addr.Ward,
	}
	if err := a.location.Dispatch(ctx, locationQuery); err != nil {
		return nil, err
	}
	loc := locationQuery.Result
	if loc.Province == nil || loc.District == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "cần cung cấp thông tin tỉnh/thành phố và quận/huyện")
	}

	addr.Province = loc.Province.Name
	addr.ProvinceCode = loc.Province.Code
	addr.District = loc.District.Name
	addr.DistrictCode = loc.District.Code
	if loc.Ward != nil {
		addr.Ward = loc.Ward.Name
		addr.WardCode = loc.Ward.Code
	}
	return addr, nil
}
