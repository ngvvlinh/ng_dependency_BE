package shopping

import (
	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
)

var (
	locationQuery     location.QueryBus
	customerQuery     *customering.QueryBus
	customerAggregate *customering.CommandBus
	addressQuery      *addressing.QueryBus
	addressAggregate  *addressing.CommandBus
	inventoryQuery    *inventory.QueryBus
	catalogQuery      *catalog.QueryBus
	catalogAggregate  *catalog.CommandBus
)

func Init(
	locationQ location.QueryBus,
	customerQ *customering.QueryBus,
	customerA *customering.CommandBus,
	addressQ *addressing.QueryBus,
	addressA *addressing.CommandBus,
	inventoryQ *inventory.QueryBus,
	catalogQ *catalog.QueryBus,
	catalogA *catalog.CommandBus,
) {

	locationQuery = locationQ
	customerQuery = customerQ
	customerAggregate = customerA
	addressQuery = addressQ
	addressAggregate = addressA
	inventoryQuery = inventoryQ
	catalogAggregate = catalogA
	catalogQuery = catalogQ
}
