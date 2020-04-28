package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ShopCollection struct {
	ID        dot.ID
	ShopID    dot.ID
	PartnerID dot.ID

	ExternalID string

	Name        string
	Description string
	DescHTML    string
	ShortDesc   string

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
