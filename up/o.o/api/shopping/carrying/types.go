package carrying

import (
	"time"

	"o.o/api/top/types/etc/status3"
	dot "o.o/capi/dot"
)

type ShopCarrier struct {
	ID        dot.ID
	ShopID    dot.ID
	FullName  string
	Note      string
	Status    status3.Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
