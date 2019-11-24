package model

import (
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenExternalAccountHaravan(&ExternalAccountHaravan{})

type ExternalAccountHaravan struct {
	ID                                dot.ID
	ShopID                            dot.ID
	Subdomain                         string
	AccessToken                       string
	ExternalShopID                    int
	ExternalCarrierServiceID          int
	ExternalConnectedCarrierServiceAt time.Time
	ExpiresAt                         time.Time
	CreatedAt                         time.Time `sq:"create"`
	UpdatedAt                         time.Time `sq:"update"`
}
