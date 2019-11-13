package model

import "time"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopTrader(&ShopTrader{})

type ShopTrader struct {
	ID        int64
	ShopID    int64
	Type      string
	DeletedAt time.Time
}
