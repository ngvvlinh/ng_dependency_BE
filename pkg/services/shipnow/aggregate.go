package shipnow

import (
	"context"

	"etop.vn/api/main/address"
	"etop.vn/api/main/identity"
	cm "etop.vn/backend/pkg/common"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/location"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
	shipnowmodelx "etop.vn/backend/pkg/services/shipnow/modelx"
	"etop.vn/backend/pkg/services/shipnow/pm"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

var _ shipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	location      location.Bus
	identityQuery identity.QueryService
	addressQuery  address.QueryService

	store sqlstore.ShipnowStoreFactory
	pm    *pm.ProcessManager
}

func NewAggregate(db cmsql.Database, location location.Bus, identityQuery identity.QueryService, addressQuery address.QueryService) *Aggregate {
	return &Aggregate{
		location:      location,
		store:         sqlstore.NewShipnowStore(db),
		identityQuery: identityQuery,
		addressQuery:  addressQuery,
	}
}

func (a *Aggregate) WithPM(pm *pm.ProcessManager) *Aggregate {
	a.pm = pm
	return a
}

func (a *Aggregate) MessageBus() shipnow.AggregateBus {
	b := bus.New()
	shipnow.NewAggregateHandler(a).RegisterHandlers(b)
	return shipnow.AggregateBus{b}
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentArgs) (*shipnow.ShipnowFulfillment, error) {
	tx, err := a.store(ctx).GetDb().Begin()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, meta.KeyTx{}, tx)
	ffmID := cm.NewID()
	// ShipnowOrderReservationEvent
	event := &shipnow.ShipnowOrderReservationEvent{
		OrderIds:             cmd.OrderIds,
		ShipnowFulfillmentId: ffmID,
	}
	orders, err := a.pm.ShipnowOrderReservation(ctx, event)
	if err != nil {
		return nil, err
	}

	var pickupAddress *ordertypes.Address
	if cmd.PickupAddress != nil {
		pickupAddress = cmd.PickupAddress
	} else {
		// prepare pickup address from shopinfo
		shopResult, err := a.identityQuery.GetShopByID(ctx, &identity.GetShopByIDQueryArgs{ID: cmd.ShopId})
		if err != nil {
			return nil, err
		}
		shop := shopResult.Shop
		shopAddressID := shop.ShipFromAddressID
		if shopAddressID == 0 {
			return nil, cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
		}
		shopAddress, err := a.addressQuery.GetAddressByID(ctx, &address.GetAddressByIDQueryArgs{ID: shopAddressID})
		if err != nil {
			return nil, err
		}
		pickupAddress = shopAddress.ToOrderAddress()
	}

	var deliveryPoints []*shipnow.DeliveryPoint
	for _, order := range orders {
		deliveryPoints = append(deliveryPoints, shipnowconvert.OrderToDeliveryPoint(order))
	}

	shipnowFfm := &shipnow.ShipnowFulfillment{
		Id:                  ffmID,
		ShopId:              cmd.ShopId,
		PickupAddress:       pickupAddress,
		DeliveryPoints:      deliveryPoints,
		Carrier:             cmd.Carrier,
		ShippingServiceCode: cmd.ShippingServiceCode,
		ShippingServiceFee:  cmd.ShippingServiceFee,
		WeightInfo:          shipnowconvert.GetWeightInfo(orders),
		ValueInfo:           shipnowconvert.GetValueInfo(orders),
		ShippingNote:        cmd.ShippingNote,
		RequestPickupAt:     nil,
	}

	if err := a.store(ctx).Create(shipnowFfm); err != nil {
		return nil, err
	}

	// tx commit
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	return shipnowFfm, err
}

func (a *Aggregate) ConfirmShipnowFulfillment(ctx context.Context, cmd *shipnow.ConfirmShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	query1 := shipnowmodelx.GetByIDArgs{
		ID:     cmd.Id,
		ShopID: cmd.ShopId,
	}
	shipnowFfm, err = a.store(ctx).GetByID(query1)
	if err != nil {
		return nil, err
	}
	if err := a.pm.ValidateConfirm(ctx, shipnowFfm); err != nil {
		return nil, err
	}
	shipnowFfmUpdate := &shipnow.ShipnowFulfillment{
		Id:            cmd.Id,
		ConfirmStatus: etoptypes.S3Positive,
	}
	shipnowFfm, err = a.store(ctx).Update(shipnowFfmUpdate)
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

func (a *Aggregate) UpdateShipnowFulfillment(ctx context.Context, cmd *shipnow.UpdateShipnowFulfillmentArgs) (shipnowFfm *shipnow.ShipnowFulfillment, err error) {
	shipnowFfm, err = a.pm.HandleUpdate(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result, err := a.store(ctx).Update(shipnowFfm)
	if err != nil {
		return nil, err
	}
	return result, err
}
