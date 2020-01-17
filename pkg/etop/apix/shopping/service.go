package shopping

import (
	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/location"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
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
