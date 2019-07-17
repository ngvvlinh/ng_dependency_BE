package identity

import "time"

type ExternalAccountHaravan struct {
	ID                                int64
	ShopID                            int64
	Subdomain                         string
	ExternalShopID                    int
	AccessToken                       string
	ExternalCarrierServiceID          int
	ExternalConnectedCarrierServiceAt time.Time
	ExpiresAt                         time.Time
	CreatedAt                         time.Time
	UpdatedAt                         time.Time
}
