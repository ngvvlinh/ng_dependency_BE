package partner

import "etop.vn/capi/dot"

type AuthMode string

const (
	AuthModeDefault AuthMode = "default"
	AuthModeManual  AuthMode = "manual"
)

type PartnerShopToken struct {
	PartnerID dot.ID `json:"partner_id"`

	ShopID         dot.ID `json:"shop_id"`
	ShopName       string `json:"shop_name"`
	ShopOwnerEmail string `json:"shop_owner_email"`
	ShopOwnerPhone string `json:"shop_owner_phone"`

	// if there is a shop with this external_shop_id in the partner_relation,
	// use the shop, otherwise automatically attach this external_shop_id to the
	// partner_relation record.
	ExternalShopID string `json:"external_shop_id"`

	AuthMode AuthMode `json:"auth_mode"`

	// if this is set, client must send the same email/phone to request_login
	RetainCurrentInfo bool   `json:"retain,omitempty"`
	RedirectURL       string `json:"redirect_url"`
}
