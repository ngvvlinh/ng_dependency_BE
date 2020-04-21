package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/capi/dot"
)

// +sqlgen
type SubscriptionProduct struct {
	ID          dot.ID
	Name        string
	Type        subscription_product_type.ProductSubscriptionType
	Description string
	ImageURL    string
	Status      status3.Status
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
	WLPartnerID dot.ID
}
