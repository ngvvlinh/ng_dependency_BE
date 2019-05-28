package shipnow

import (
	"context"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/location"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	shipnowmodel "etop.vn/backend/pkg/services/shipnow/model"
	"etop.vn/backend/pkg/services/shipnow/pm"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

var _ shipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	location location.Bus
	s        *sqlstore.ShipnowStore
	pm       *pm.ProcessManager
}

func NewAggregate(db cmsql.Database, location location.Bus) *Aggregate {
	return &Aggregate{
		location: location,
		s:        sqlstore.NewShipnowStore(db),
	}
}

func (a *Aggregate) WithPM(pm *pm.ProcessManager) *Aggregate {
	a.pm = pm
	return a
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentArgs) (*shipnow.ShipnowFulfillment, error) {
	shipnowFfm, err := a.HandleCreation(ctx, cmd)
	if err != nil {
		return nil, err
	}
	if err := a.s.WithContext(ctx).Create(shipnowFfm); err != nil {
		return nil, err
	}
	return shipnowFfm, err
}

func (a *Aggregate) ConfirmShipnowFulfillment(ctx context.Context, cmd *shipnow.ConfirmShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	query1 := shipnowmodel.GetByIDArgs{
		ID:     cmd.Id,
		ShopID: cmd.ShopId,
	}
	shipnowFfm, err = a.s.WithContext(ctx).GetByID(query1)
	if err != nil {
		return nil, err
	}
	if err := a.ValidateConfirm(ctx, shipnowFfm); err != nil {
		return nil, err
	}
	shipnowFfmUpdate := &shipnow.ShipnowFulfillment{
		Id:            cmd.Id,
		ConfirmStatus: etoptypes.S3Positive,
	}
	shipnowFfm, err = a.s.WithContext(ctx).Update(shipnowFfmUpdate)
	if err != nil {
		return nil, err
	}

	// if err := a.shipnowManagerCtrl.CreateExternalShipping(ctx, ffm); err != nil {
	// 	return &meta.Empty{}, err
	// }

	return shipnowFfm, nil
}

func (a *Aggregate) CancelShipnowFulfillment(ctx context.Context, cmd *shipnow.CancelShipnowFulfillmentArgs) (*meta.Empty, error) {
	err := a.pm.HandleShipnowCancellation(ctx, cmd)
	return &meta.Empty{}, err
}

func (a *Aggregate) HandleCreation(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	return a.pm.HandleShipnowCreation(ctx, cmd)
}

func (a *Aggregate) ValidateConfirm(ctx context.Context, ffm *shipnow.ShipnowFulfillment) error {
	return cm.ErrTODO

	// switch ffm.ConfirmStatus {
	// case etoptypes.S3Negative:
	// 	return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã hủy")
	// case etoptypes.S3Positive:
	// 	return cm.Errorf(cm.FailedPrecondition, nil, "Đơn giao hàng đã xác nhận")
	// }
	// if ffm.Status == etoptypes.S5Negative || ffm.Status == etoptypes.S5Positive {
	// 	return cm.Errorf(cm.FailedPrecondition, nil, "Không thể xác nhận đơn giao hàng này")
	// }
	//
	// if len(ffm.DeliveryPoints) == 0 {
	// 	return cm.Errorf(cm.FailedPrecondition, nil, "Số điểm giao hàng không hợp lệ.")
	// }
	// var orderIDs []int64
	// for _, point := range ffm.DeliveryPoints {
	// 	if point.OrderId == 0 {
	// 		continue
	// 	}
	// 	orderIDs = append(orderIDs, point.OrderId)
	// }
	// if err := a.ValidateOrders(ctx, orderIDs, ffm.Id); err != nil {
	// 	return err
	// }
	//
	// ffm.ConfirmStatus = etoptypes.S3Positive
	// return nil
}

func (a *Aggregate) UpdateShipnowFulfillment(ctx context.Context, cmd *shipnow.UpdateShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	shipnowFfm, err = a.HandleUpdate(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := a.s.WithContext(ctx).Update(shipnowFfm)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (a *Aggregate) HandleUpdate(ctx context.Context, cmd *shipnow.UpdateShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	return nil, cm.ErrTODO

	// query1 := shipnowmodel.GetByIDArgs{
	// 	ID:     cmd.Id,
	// 	ShopID: cmd.ShopId,
	// }
	// dbFfm, err := a.s.WithContext(ctx).GetByID(query1)
	// if err != nil {
	// 	return nil, err
	// }
	// ffm := shipnowconvert.Shipnow(dbFfm)
	// if ffm.ConfirmStatus != etoptypes.S3Zero || ffm.ShippingCode != "" {
	// 	return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể cập nhật đơn giao hàng này.")
	// }
	//
	// shipnowFfm = &shipnow.ShipnowFulfillment{
	// 	Id:                  cmd.Id,
	// 	ShopId:              cmd.ShopId,
	// 	PickupAddress:       cmd.PickupAddress,
	// 	Carrier:             cmd.Carrier,
	// 	ShippingServiceCode: cmd.ShippingServiceCode,
	// 	ShippingServiceFee:  cmd.ShippingServiceFee,
	// 	ShippingNote:        cmd.ShippingNote,
	// 	RequestPickupAt:     nil,
	// }
	//
	// if len(cmd.OrderIds) > 0 {
	// 	if err := a.ValidateOrders(ctx, cmd.OrderIds, cmd.Id); err != nil {
	// 		return nil, err
	// 	}
	// 	cmd3 := &ordering.GetOrdersArgs{
	// 		ShopID: cmd.ShopId,
	// 		IDs:    cmd.OrderIds,
	// 	}
	// 	orders, err := a.orderPM.GetOrders(ctx, cmd3)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	shipnowFfm.WeightInfo = shipnowconvert.GetWeightInfo(orders)
	// 	shipnowFfm.ValueInfo = shipnowconvert.GetValueInfo(orders)
	// 	var deliveryPoints []*shipnow.DeliveryPoint
	// 	for _, order := range orders {
	// 		deliveryPoints = append(deliveryPoints, shipnowconvert.OrderToDeliveryPoint(order))
	// 	}
	// 	shipnowFfm.DeliveryPoints = deliveryPoints
	// }
	// return shipnowFfm, nil
}
