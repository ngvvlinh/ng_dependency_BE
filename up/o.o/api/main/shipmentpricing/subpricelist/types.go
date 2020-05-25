package subpricelist

import (
	"time"

	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:event:topic=event/shipmentsubpricelist

type ShipmentSubPriceList struct {
	ID           dot.ID
	Name         string
	Description  string
	Status       status3.Status
	ConnectionID dot.ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	WLPartnerID  dot.ID
}

type ShipmentSubPriceListUpdatedEvent struct {
	meta.EventMeta
	ID dot.ID
}

type ShipmentSubPriceListDeletingEvent struct {
	meta.EventMeta
	ID dot.ID
}

type ShipmentSubPriceListDeletedEvent struct {
	meta.EventMeta
	ID dot.ID
}
