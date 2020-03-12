package shipmentservice

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

type ShipmentService struct {
	ID           dot.ID         `json:"id"`
	ConnectionID dot.ID         `json:"connection_id"`
	Name         string         `json:"name"`
	EdCode       string         `json:"ed_code"`
	ServiceIDs   []string       `json:"service_ids"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    time.Time      `json:"-"`
	WLPartnerID  dot.ID         `json:"-"`
	ImageURL     string         `json:"image_url"`
	Status       status3.Status `json:"status"`
}
