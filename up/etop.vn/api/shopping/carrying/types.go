package carrying

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	dot "etop.vn/capi/dot"
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
