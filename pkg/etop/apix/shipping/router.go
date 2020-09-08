package shipping

import (
	"o.o/api/main/location"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipping"
	"o.o/api/top/external/types"
	com "o.o/backend/com/main"
	ordersqlstore "o.o/backend/com/main/ordering/sqlstore"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shipsqlstore "o.o/backend/com/main/shipping/sqlstore"
	orderS "o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/etop/sqlstore"
)

var locationList *types.LocationResponse

type Shipping struct {
	LocationBus      location.QueryBus
	LocationList     *types.LocationResponse
	OrderStore       ordersqlstore.OrderStoreFactory
	FulfillmentStore shipsqlstore.FulfillmentStoreFactory
	ShipmentManager  *shippingcarrier.ShipmentManager
	ShippingAggr     shipping.CommandBus
	ShippingQuery    shipping.QueryBus
	OrderLogic       *orderS.OrderLogic
	ShipnowAggr      shipnow.CommandBus
	ShipnowQuery     shipnow.QueryBus

	OrderStoreIface sqlstore.OrderStoreInterface
}

func New(
	locationBus location.QueryBus,
	db com.MainDB,
	shipmentM *shippingcarrier.ShipmentManager,
	shippingA shipping.CommandBus,
	shippingQ shipping.QueryBus,
	orderLogic *orderS.OrderLogic,
	shipnowA shipnow.CommandBus,
	shipnowQ shipnow.QueryBus,
) *Shipping {
	s := &Shipping{
		LocationBus:      locationBus,
		OrderStore:       ordersqlstore.NewOrderStore(db),
		FulfillmentStore: shipsqlstore.NewFulfillmentStore(db),
		ShipmentManager:  shipmentM,
		ShippingAggr:     shippingA,
		ShippingQuery:    shippingQ,
		OrderLogic:       orderLogic,
		ShipnowAggr:      shipnowA,
		ShipnowQuery:     shipnowQ,
	}
	locationList = buildLocationList(locationBus)
	return s
}
