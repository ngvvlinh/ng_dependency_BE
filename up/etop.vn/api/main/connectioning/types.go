package connectioning

import (
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:event:topic=event/connection

var (
	DefaultTopshipGHTKConnectionID   = dot.ID(1000804010396750738)
	DefaultTopshipGHNConnectionID    = dot.ID(1000803215822389663)
	DefaultTopshipVTPostConnectionID = dot.ID(1000804104889339180)
)

type Connection struct {
	ID                   dot.ID
	Name                 string
	Status               status3.Status
	PartnerID            dot.ID
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	DriverConfig         *ConnectionDriverConfig
	Driver               string
	ConnectionType       connection_type.ConnectionType
	ConnectionSubtype    connection_type.ConnectionSubtype
	ConnectionMethod     connection_type.ConnectionMethod
	ConnectionProvider   connection_type.ConnectionProvider
	EtopAffiliateAccount *EtopAffiliateAccount
	Code                 string
	ImageURL             string
	Services             []*ConnectionService
	WLPartnerID          dot.ID
}

type ConnectionDriverConfig struct {
	TrackingURL            string `json:"tracking_url"`
	CreateFulfillmentURL   string `json:"create_fulfillment_url"`
	GetFulfillmentURL      string `json:"get_fulfillment_url"`
	GetShippingServicesURL string `json:"get_shipping_services_url"`
	CancelFulfillmentURL   string `json:"cancel_fulfillment_url"`
	SignInURL              string `json:"sign_in_url"`
	SignUpURL              string `json:"sign_up_url"`
}

type EtopAffiliateAccount struct {
	UserID string
	Token  string
}

type ShopConnection struct {
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	Status         status3.Status
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	IsGlobal       bool
	ExternalData   *ShopConnectionExternalData
}

type ShopConnectionExternalData struct {
	UserID string
	Email  string
}

type ConnectionUpdatedEvent struct {
	meta.EventMeta
	ConnectionID dot.ID
}

type ShopConnectionUpdatedEvent struct {
	meta.EventMeta
	ShopID       dot.ID
	ConnectionID dot.ID
}

type ConnectionService struct {
	ServiceID string
	Name      string
}
