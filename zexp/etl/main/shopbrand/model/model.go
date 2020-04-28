package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ShopBrand struct {
	ID     dot.ID
	ShopID dot.ID

	BrandName   string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
