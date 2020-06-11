package shopshipmentpricelist

import (
	"time"

	"o.o/capi/dot"
)

type ShopShipmentPriceList struct {
	ShopID              dot.ID
	ShipmentPriceListID dot.ID
	ConnectionID        dot.ID
	Note                string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           time.Time
	UpdatedBy           dot.ID
}
