package model

import (
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopCarrier(&ShopCarrier{})

type ShopCarrier struct {
	ID        dot.ID
	ShopID    dot.ID
	FullName  string
	Note      string
	Status    int
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}
