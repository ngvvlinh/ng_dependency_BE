package model

import (
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopBrand(&ShopBrand{})

type ShopBrand struct {
	ID         dot.ID
	ShopID     dot.ID
	ExternalID string
	PartnerID  dot.ID

	BrandName   string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Rid dot.ID
}
