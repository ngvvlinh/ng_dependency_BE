package pricelist

import (
	"time"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:event:topic=event/shipmentpricelist

type ShipmentPriceList struct {
	ID                      dot.ID    `json:"id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	WLPartnerID             dot.ID    `json:"-"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"-"`
	UpdatedAt               time.Time `json:"-"`
	DeletedAt               time.Time `json:"-"`
	ShipmentSubPriceListIDs []dot.ID  `json:"-"`
}

type DeleteCachePriceListEvent struct {
	meta.EventMeta
	ShipmentPriceListID dot.ID
}
