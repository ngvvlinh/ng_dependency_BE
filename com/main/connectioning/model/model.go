package model

import (
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

// +sqlgen
type Connection struct {
	ID                   dot.ID
	Name                 string
	Status               status3.Status
	PartnerID            dot.ID
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
	DriverConfig         *connectioning.ConnectionDriverConfig `json:"driver_config"`
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
	OriginConnectionID   dot.ID
	Version              string // version api
}

// +sqlgen
type ShopConnection struct {
	OwnerID          dot.ID
	ShopID           dot.ID
	ConnectionID     dot.ID
	Token            string
	TokenExpiresAt   time.Time
	Status           status3.Status
	ConnectionStates *ConnectionStates
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
	DeletedAt        time.Time
	IsGlobal         bool
	ExternalData     *ShopConnectionExternalData
}

type ConnectionStates struct {
	Error *model.Error `json:"error"`
}

type EtopAffiliateAccount struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
	// shop_id used for GHN
	ShopID string `json:"shop_id"`
	// client_secret for NinjaVan
	SecretKey string `json:"secret_key,omitempty"`
}

type ShopConnectionExternalData struct {
	UserID string `json:"user_id"`
	// old: email
	// new: identifier include either email or phone
	Identifier string `json:"identifier"` // email or phone
	ShopID     string `json:"shop_id"`
}

type ConnectionService struct {
	ServiceID string `json:"service_id"`
	Name      string `json:"name"`
}
