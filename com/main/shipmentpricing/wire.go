package shipmentpricing

import (
	"github.com/google/wire"
	"o.o/backend/com/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion"
	"o.o/backend/com/main/shipmentpricing/shipmentprice"
	"o.o/backend/com/main/shipmentpricing/shipmentservice"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist"
)

var WireSet = wire.NewSet(
	pricelist.WireSet,
	shipmentprice.WireSet,
	shipmentservice.WireSet,
	shopshipmentpricelist.WireSet,
	pricelistpromotion.WireSet,
)
