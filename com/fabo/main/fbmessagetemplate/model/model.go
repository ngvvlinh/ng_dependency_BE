package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type FbMessageTemplate struct {
	ID        dot.ID
	ShopID    dot.ID
	Template  string
	ShortCode string
	CreatedAt time.Time
	UpdatedAt time.Time
}
