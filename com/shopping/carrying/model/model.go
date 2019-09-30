package model

import (
	"time"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopCarrier(&ShopCarrier{})

type ShopCarrier struct {
	ID        int64
	ShopID    int64
	FullName  string
	Note      string
	Status    int32
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}
