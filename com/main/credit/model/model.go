package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/capi/dot"
)

// +sqlgen
type Credit struct {
	ID        dot.ID
	Amount    int
	ShopID    dot.ID
	Type      string
	Status    status3.Status
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	PaidAt    time.Time
}

// +sqlgen:           Credit as c
// +sqlgen:left-join: Shop   as s on s.id = c.shop_id
type CreditExtended struct {
	*Credit
	Shop *identitymodel.Shop
}
