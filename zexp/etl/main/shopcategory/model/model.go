package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ShopCategory struct {
	ID     dot.ID
	ShopID dot.ID

	ParentID dot.ID

	Name string

	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
