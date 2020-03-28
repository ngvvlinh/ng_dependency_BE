package model

import (
	"time"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
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
}

// +sqlgen
type ShopConnection struct {
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
}

type ShopConnectionExternalData struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
