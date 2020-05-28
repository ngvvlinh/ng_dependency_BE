package shopping

import (
	"github.com/google/wire"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
)

var WireSet = wire.NewSet(
	wire.Struct(new(Shopping), "*"),
)

type Shopping struct {
	LocationQuery     location.QueryBus
	CustomerQuery     customering.QueryBus
	CustomerAggregate customering.CommandBus
	AddressQuery      addressing.QueryBus
	AddressAggregate  addressing.CommandBus
	InventoryQuery    inventory.QueryBus
	CatalogQuery      catalog.QueryBus
	CatalogAggregate  catalog.CommandBus
}
