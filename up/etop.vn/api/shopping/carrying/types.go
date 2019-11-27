package carrying

import (
	"time"

	dot "etop.vn/capi/dot"
)

type ShopCarrier struct {
	ID        dot.ID
	ShopID    dot.ID
	FullName  string
	Note      string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
