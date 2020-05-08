package partnercarrier

import (
	"time"

	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

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
