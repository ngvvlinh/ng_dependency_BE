package identity

import (
	"time"

	"etop.vn/capi/dot"
)

type ExternalAccountHaravan struct {
	ID                                dot.ID
	ShopID                            dot.ID
	Subdomain                         string
	ExternalShopID                    int
	AccessToken                       string
	ExternalCarrierServiceID          int
	ExternalConnectedCarrierServiceAt time.Time
	ExpiresAt                         time.Time
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
}
