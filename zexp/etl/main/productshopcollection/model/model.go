package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ProductShopCollection struct {
	CollectionID dot.ID
	ProductID    dot.ID
	ShopID       dot.ID
	Status       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Rid          dot.ID
}
