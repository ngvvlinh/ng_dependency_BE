package pm

import (
	"context"

	"etop.vn/api/main/address"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/ordering"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	cm "etop.vn/backend/pkg/common"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
)

type ProcessManager struct {
	shipnow shipnow.Aggregate
	order   ordering.Aggregate

	identityQuery  identity.QueryService
	orderQuery     ordering.QueryService
	addressQuery   address.QueryService
	carrierManager carrier.Manager
}

func New(
	shipnowAggr shipnow.Aggregate,
	orderAggr ordering.Aggregate,
	identityQuery identity.QueryService,
	orderQuery ordering.QueryService,
	addressQuery address.QueryService,
	carrierManager carrier.Manager,
) *ProcessManager {
	return &ProcessManager{
		shipnow:        shipnowAggr,
		order:          orderAggr,
		identityQuery:  identityQuery,
		orderQuery:     orderQuery,
		addressQuery:   addressQuery,
		carrierManager: carrierManager,
	}
}

func (m *ProcessManager) GetActiveShipnowFulfillments(ctx context.Context, cmd *shipnow.GetActiveShipnowFulfillmentsCommand) ([]*shipnow.ShipnowFulfillment, error) {
	return m.shipnow.GetActiveShipnowFulfillments(ctx, cmd)
}

func (m *ProcessManager) HandleShipnowCreation(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentCommand) (*shipnow.ShipnowFulfillment, error) {
	if err := m.validateOrders(ctx, cmd.OrderIds, 0); err != nil {
		return nil, err
	}
	cmd3 := &ordering.GetOrdersArgs{
		ShopID: cmd.ShopId,
		IDs:    cmd.OrderIds,
	}
	orders, err := m.order.GetOrders(ctx, cmd3)
	if err != nil {
		return nil, err
	}

	var pickupAddress *ordertypes.Address
	if cmd.PickupAddress != nil {
		pickupAddress = cmd.PickupAddress
	} else {
		// prepare pickup address from shopinfo
		shop, err := m.identityQuery.GetShopByID(ctx, &identity.GetShopByIDQueryArgs{ID: cmd.ShopId})
		if err != nil {
			return nil, err
		}
		shopAddressID := shop.ShipFromAddressID
		if shopAddressID == 0 {
			return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
		}
		shopAddress, err := m.addressQuery.GetAddressByID(ctx, &address.GetAddressByIDQueryArgs{ID: shopAddressID})
		if err != nil {
			return nil, err
		}
		pickupAddress = shopAddress.ToOrderAddress()
	}

	var deliveryPoints []*shipnow.DeliveryPoint
	for _, order := range orders {
		deliveryPoints = append(deliveryPoints, shipnowconvert.OrderToDeliveryPoint(order))
	}

	_ = pickupAddress

	// TODO: implement it
	//
	// shipnowFfm := &shipnow.ShipnowFulfillment{
	// 	Id:                  cm.NewID(),
	// 	ShopId:              cmd.ShopId,
	// 	PickupAddress:       pickupAddress,
	// 	DeliveryPoints:      deliveryPoints,
	// 	Carrier:             cmd.Carrier,
	// 	ShippingServiceCode: cmd.ShippingServiceCode,
	// 	ShippingServiceFee:  cmd.ShippingServiceFee,
	// 	WeightInfo:          shipnowconvert.GetWeightInfo(orders),
	// 	ValueInfo:           shipnowconvert.GetValueInfo(orders),
	// 	ShippingNote:        cmd.ShippingNote,
	// 	RequestPickupAt:     nil,
	// }
	// if err := m.carrierManager.CreateExternalShipping(ctx, nil); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (m *ProcessManager) HandleShipnowCancellation(ctx context.Context, cmd *shipnow.CancelShipnowFulfillmentCommand) error {

	// TODO
	return cm.ErrTODO

	// ffm, err := a.s.WithContext(ctx).GetByID(shipnowmodel.GetByIDArgs{
	// 	ID: cmd.Id,
	// })
	// if err != nil {
	// 	return &meta.Empty{}, err
	// }
	// switch ffm.Status {
	// case etoptypes.S5Positive, etoptypes.S5Negative, etoptypes.S5NegSuper:
	// 	return &meta.Empty{}, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển không thể hủy")
	// }
	// switch ffm.ShippingState {
	// case shipnowtypes.StateCancelled:
	// 	return &meta.Empty{}, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đã bị hủy")
	// case shipnowtypes.StateDelivering:
	// 	return &meta.Empty{}, cm.Errorf(cm.FailedPrecondition, nil, "Đơn vận chuyển đang giao. Không thể hủy đơn.")
	// case shipnowtypes.StateDelivered,
	// 	shipnowtypes.StateReturning, shipnowtypes.StateReturned:
	// 	return &meta.Empty{}, cm.Errorf(cm.FailedPrecondition, nil, "Không thể hủy đơn.")
	// }
	//
	// if err := ctrl.CancelExternalShipping(ctx, ffm); err != nil {
	// 	return &meta.Empty{}, err
	// }
	//
	// updateArgs := sqlstore.UpdateSyncStateArgs{
	// 	ID:         ffm.Id,
	// 	SyncStatus: etoptypes.S4Negative,
	// 	State:      shipnowtypes.StateCancelled,
	// 	SyncStates: &model.FulfillmentSyncStates{
	// 		TrySyncAt:         time.Now(),
	// 		NextShippingState: model.StateCreated,
	// 	},
	// }
	// ffm, err = a.s.WithContext(ctx).UpdateSyncState(updateArgs)
	// if err != nil {
	// 	return &meta.Empty{}, err
	// }
	// return &meta.Empty{}, nil
	//
	// return m.carrierManager.CancelExternalShipping(ctx, nil)
}

func (m *ProcessManager) validateOrders(ctx context.Context, orderIDs []int64, shipnowFfmID int64) error {
	if len(orderIDs) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn đơn hàng")
	}
	cmd := &ordering.ValidateOrdersForShippingCommand{
		OrderIDs: orderIDs,
	}
	if err := m.order.ValidateOrders(ctx, cmd); err != nil {
		return err
	}
	// for _, id := range orderIDs {
	// 	cmd1 := &shipnowmodel.GetActiveShipnowFulfillmentsByOrderIDArgs{
	// 		OrderID:                     id,
	// 		ExcludeShipnowFulfillmentID: shipnowFfmID,
	// 	}
	// 	ffms, err := m.s.WithContext(ctx).GetActiveShipnowFulfillmentsByOrderID(cmd1)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(ffms) > 0 {
	// 		return cm.Errorf(cm.FailedPrecondition, nil, "Đơn hàng %v đã thuộc đơn vận chuyển khác (%v)", id, ffms[0].Id)
	// 	}
	// }
	return nil
}
