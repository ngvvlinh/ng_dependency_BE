package shipnow

import (
	"context"

	"etop.vn/api/main/ordering"

	"etop.vn/api/main/address"
	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/location"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
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

	db       cmsql.Transactioner
	store    sqlstore.ShipnowStoreFactory
	eventBus meta.EventBus
	pm       *pm.ProcessManager // nil
}

func NewAggregate(eventBus meta.EventBus, db cmsql.Database, location location.Bus, identityQuery identity.QueryService, addressQuery address.QueryService) *Aggregate {
	return &Aggregate{
		db:       db,
		store:    sqlstore.NewShipnowStore(db),
		eventBus: eventBus,

		location:      location,
		identityQuery: identityQuery,
		addressQuery:  addressQuery,
	}
}

func (a *Aggregate) MessageBus() shipnow.AggregateBus {
	b := bus.New()
	shipnow.NewAggregateHandler(a).RegisterHandlers(b)
	return shipnow.AggregateBus{b}
}

func (a *Aggregate) CreateShipnowFulfillment(ctx context.Context, cmd *shipnow.CreateShipnowFulfillmentArgs) (_result *shipnow.ShipnowFulfillment, _ error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		ffmID := cm.NewID()
		// ShipnowOrderReservationEvent
		event := &shipnow.ShipnowOrderReservationEvent{
			OrderIds:             cmd.OrderIds,
			ShipnowFulfillmentId: ffmID,
		}
		err := a.pm.ShipnowOrderReservation(ctx, event)
		if err != nil {
			return err
		}
		var orders []*ordering.Order // TODO: fix it

		var pickupAddress *ordertypes.Address
		if cmd.PickupAddress != nil {
			pickupAddress = cmd.PickupAddress
		} else {
			// prepare pickup address from shopinfo
			shopResult, err := a.identityQuery.GetShopByID(ctx, &identity.GetShopByIDQueryArgs{ID: cmd.ShopId})
			if err != nil {
				return err
			}
			shop := shopResult.Shop
			shopAddressID := shop.ShipFromAddressID
			if shopAddressID == 0 {
				return cm.Error(cm.InvalidArgument, "Bán hàng: Cần cung cấp thông tin địa chỉ lấy hàng trong đơn hàng hoặc tại thông tin cửa hàng. Vui lòng cập nhật.", nil)
			}
			shopAddress, err := a.addressQuery.GetAddressByID(ctx, &address.GetAddressByIDQueryArgs{ID: shopAddressID})
			if err != nil {
				return err
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
			return err
		}
		_result = shipnowFfm
		return nil
	})
	return _result, err
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