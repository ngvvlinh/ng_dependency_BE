package shipmentpricing

import (
	"github.com/google/wire"
	"o.o/backend/com/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/com/main/shipmentpricing/shipmentservice"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist"
	"o.o/backend/com/main/shipmentpricing/subpricelist"
)

var WireSet = wire.NewSet(
	pricelist.WireSet,
	shipmentprice.WireSet,
	shipmentservice.WireSet,
	subpricelist.WireSet,
	shopshipmentpricelist.WireSet,
)
