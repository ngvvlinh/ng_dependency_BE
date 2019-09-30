package carrying

import "time"

type ShopCarrier struct {
	ID        int64
	ShopID    int64
	FullName  string
	Note      string
	Status    int32
	CreatedAt time.Time
	UpdatedAt time.Time
}
