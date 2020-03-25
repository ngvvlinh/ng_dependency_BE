package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopProductCollection struct {
	ShopID dot.ID

	ProductID    dot.ID
	CollectionID dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
