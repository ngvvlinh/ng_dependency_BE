package partner_proto

import (
	"time"

	"etop.vn/api/top/types/etc/authorize_shop_config"
	shippingstate "etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type AuthorizeShopRequest struct {
	ShopId         dot.ID `json:"shop_id"`
	ExternalShopID string `json:"external_shop_id"`
	ExternalUserID string `json:"external_user_id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	RedirectUrl    string `json:"redirect_url"`
	ShopName       string `json:"shop_name"`
	ExtraToken     string `json:"extra_token"`

	Config []authorize_shop_config.AuthorizeShopConfig `json:"config"`
}

func (m *AuthorizeShopRequest) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizeShopResponse struct {
	Code      string            `json:"code"`
	Msg       string            `json:"msg"`
	Type      string            `json:"type"`
	AuthToken string            `json:"auth_token"`
	ExpiresIn int               `json:"expires_in"`
	AuthUrl   string            `json:"auth_url"`
	Meta      map[string]string `json:"meta"`
}

func (m *AuthorizeShopResponse) String() string { return jsonx.MustMarshalToString(m) }

type ShipmentConnection struct {
	ID                     dot.ID         `json:"id"`
	Name                   string         `json:"name"`
	Status                 status3.Status `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	ImageURL               string         `json:"image_url"`
	TrackingURL            string         `json:"tracking_url"`
	CreateFulfillmentURL   string         `json:"create_fulfillment_url"`
	GetFulfillmentURL      string         `json:"get_fulfillment_url"`
	CancelFulfillmentURL   string         `json:"cancel_fulfillment_url"`
	GetShippingServicesURL string         `json:"get_shipping_services_url"`
	SignInURL              string         `json:"sign_in_url"`
	SignUpURL              string         `json:"sign_up_url"`
}

func (m *ShipmentConnection) String() string { return jsonx.MustMarshalToString(m) }

type CreateConnectionRequest struct {
	Name                   string `json:"name"`
	ImageURL               string `json:"image_url"`
	TrackingURL            string `json:"tracking_url"`
	CreateFulfillmentURL   string `json:"create_fulfillment_url"`
	GetFulfillmentURL      string `json:"get_fulfillment_url"`
	CancelFulfillmentURL   string `json:"cancel_fulfillment_url"`
	GetShippingServicesURL string `json:"get_shipping_services_url"`
	SignInURL              string `json:"sign_in_url"`
	SignUpURL              string `json:"sign_up_url"`
}

func (m *CreateConnectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateConnectionRequest struct {
	ID                     dot.ID `json:"id"`
	Name                   string `json:"name"`
	ImageURL               string `json:"image_url"`
	TrackingURL            string `json:"tracking_url"`
	CreateFulfillmentURL   string `json:"create_fulfillment_url"`
	GetFulfillmentURL      string `json:"get_fulfillment_url"`
	CancelFulfillmentURL   string `json:"cancel_fulfillment_url"`
	GetShippingServicesURL string `json:"get_shipping_services_url"`
	SignInURL              string `json:"sign_in_url"`
	SignUpURL              string `json:"sign_up_url"`
}

func (m *UpdateConnectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetConnectionsResponse struct {
	Connections []*ShipmentConnection `json:"connections"`
}

func (m *GetConnectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateFulfillmentRequest struct {
	ShippingCode  string                  `json:"shipping_code"`
	ShippingState shippingstate.NullState `json:"shipping_state"`
}

func (m *UpdateFulfillmentRequest) String() string { return jsonx.MustMarshalToString(m) }
