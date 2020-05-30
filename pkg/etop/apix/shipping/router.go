package shipping

import (
	"o.o/api/main/location"
	"o.o/api/main/shipping"
	"o.o/api/top/external/types"
	com "o.o/backend/com/main"
	ordersqlstore "o.o/backend/com/main/ordering/sqlstore"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shipsqlstore "o.o/backend/com/main/shipping/sqlstore"
	orderS "o.o/backend/pkg/etop/logic/orders"
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
}

func New(
	locationBus location.QueryBus,
	db com.MainDB,
	shipmentM *shippingcarrier.ShipmentManager,
	shippingA shipping.CommandBus,
	shippingQ shipping.QueryBus,
	orderLogic *orderS.OrderLogic,
) *Shipping {
	s := &Shipping{
		LocationBus:      locationBus,
		OrderStore:       ordersqlstore.NewOrderStore(db),
		FulfillmentStore: shipsqlstore.NewFulfillmentStore(db),
		ShipmentManager:  shipmentM,
		ShippingAggr:     shippingA,
		ShippingQuery:    shippingQ,
		OrderLogic:       orderLogic,
	}
	locationList = buildLocationList(locationBus)
	return s
}
