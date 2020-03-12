package pricelist

import (
	"time"

	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/shipmentpricelist

type ShipmentPriceList struct {
	ID          dot.ID    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WLPartnerID dot.ID    `json:"-"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
}

type ShipmentPriceListActivatedEvent struct {
	meta.EventMeta
	ID dot.ID
}
