package model

import "time"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenExternalAccountHaravan(&ExternalAccountHaravan{})

type ExternalAccountHaravan struct {
	ID                                int64
	ShopID                            int64
	Subdomain                         string
	AccessToken                       string
	ExternalShopID                    int
	ExternalCarrierServiceID          int
	ExternalConnectedCarrierServiceAt time.Time
	ExpiresAt                         time.Time
	CreatedAt                         time.Time `sq:"create"`
	UpdatedAt                         time.Time `sq:"update"`
}
