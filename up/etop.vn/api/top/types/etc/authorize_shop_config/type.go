package authorize_shop_config

// +enum
type AuthorizeShopConfig int

type NullAuthorizeShopConfig struct {
	Enum  AuthorizeShopConfig
	Valid bool
}

const (
	// +enum=shipment
	Shipment AuthorizeShopConfig = 1

	// +enum=whitelabel
	WhiteLabel AuthorizeShopConfig = 2
)
