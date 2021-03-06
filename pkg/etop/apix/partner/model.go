package partner

import "o.o/capi/dot"

type AuthMode string
type AuthType string

const (
	AuthModeDefault AuthMode = "default"
	AuthModeManual  AuthMode = "manual"

	AuthTypeUserKey AuthType = "user_key"
	AuthTypeShopKey AuthType = "shop_key"
)

type PartnerShopToken struct {
	PartnerID dot.ID `json:"partner_id"`

	ShopID         dot.ID `json:"shop_id"`
	ShopName       string `json:"shop_name"`
	ShopOwnerEmail string `json:"shop_owner_email"`
	ShopOwnerPhone string `json:"shop_owner_phone"`
	ShopOwnerName  string `json:"shop_owner_name"`

	// if there is a shop with this external_shop_id in the partner_relation,
	// use the shop, otherwise automatically attach this external_shop_id to the
	// partner_relation record.
	ExternalShopID string `json:"external_shop_id"`
	ExternalUserID string `json:"external_user_id"`
	ExtraToken     string `json:"extra_token"`

	AuthMode AuthMode `json:"auth_mode"`
	AuthType AuthType `json:"auth_key"`

	// if this is set, client must send the same email/phone to request_login
	RetainCurrentInfo bool   `json:"retain,omitempty"`
	RedirectURL       string `json:"redirect_url"`
	// Config: "shipment,..."
	Config string `json:"config"`
}
