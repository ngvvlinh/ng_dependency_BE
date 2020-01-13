package partner_proto

import (
	"etop.vn/api/top/types/etc/authorize_shop_config"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type AuthorizeShopRequest struct {
	ShopId         dot.ID `json:"shop_id"`
	ExternalShopId string `json:"external_shop_id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	RedirectUrl    string `json:"redirect_url"`

	Config []authorize_shop_config.AuthorizeShopConfig `json:"config"`
}

func (m *AuthorizeShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizeShopResponse struct {
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Type      string `json:"type"`
	AuthToken string `json:"auth_token"`
	ExpiresIn int    `json:"expires_in"`
	AuthUrl   string `json:"auth_url"`
}

func (m *AuthorizeShopResponse) String() string { return jsonx.MustMarshalToString(m) }
