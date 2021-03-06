package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ProductShopCollection struct {
	CollectionID dot.ID
	ProductID    dot.ID
	ShopID       dot.ID
	Status       int `sql_gen:"int2"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Rid          dot.ID
}
