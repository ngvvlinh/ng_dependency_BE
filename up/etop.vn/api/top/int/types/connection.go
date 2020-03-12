package types

import (
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type LoginShopConnectionRequest struct {
	ConnectionID dot.ID `json:"connection_id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

func (m *LoginShopConnectionRequest) Reset() {
	*m = LoginShopConnectionRequest{}
}
func (m *LoginShopConnectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterShopConnectionRequest struct {
	ConnectionID dot.ID `json:"connection_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	Province     string `json:"province"`
	District     string `json:"district"`
	Address      string `json:"address"`
}

func (m *RegisterShopConnectionRequest) Reset() {
	*m = RegisterShopConnectionRequest{}
}
func (m *RegisterShopConnectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type Connection struct {
	ID                 dot.ID                             `json:"id"`
	Name               string                             `json:"name"`
	Status             status3.Status                     `json:"status"`
	CreatedAt          dot.Time                           `json:"created_at"`
	UpdatedAt          dot.Time                           `json:"updated_at"`
	ConnectionType     connection_type.ConnectionType     `json:"connection_type"`
	ConnectionSubtype  connection_type.ConnectionSubtype  `json:"connection_subtype"`
	ConnectionMethod   connection_type.ConnectionMethod   `json:"connection_method"`
	ConnectionProvider connection_type.ConnectionProvider `json:"connection_provider"`
}

func (m *Connection) Reset() {
	*m = Connection{}
}
func (m *Connection) String() string { return jsonx.MustMarshalToString(m) }

type ShopConnection struct {
	ShopID       dot.ID                      `json:"shop_id"`
	ConnectionID dot.ID                      `json:"connection_id"`
	Status       status3.Status              `json:"status"`
	CreatedAt    dot.Time                    `json:"created_at"`
	UpdatedAt    dot.Time                    `json:"updated_at"`
	DeletedAt    dot.Time                    `json:"deleted_at"`
	IsGlobal     bool                        `json:"is_global"`
	ExternalData *ShopConnectionExternalData `json:"external_data"`
}

func (m *ShopConnection) Reset() {
	*m = ShopConnection{}
}
func (m *ShopConnection) String() string { return jsonx.MustMarshalToString(m) }

type ShopConnectionExternalData struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
}

func (m *ShopConnectionExternalData) Reset() {
	*m = ShopConnectionExternalData{}
}
func (m *ShopConnectionExternalData) String() string { return jsonx.MustMarshalToString(m) }

type GetConnectionsResponse struct {
	Connections []*Connection `json:"connections"`
}

func (m *GetConnectionsResponse) Reset() {
	*m = GetConnectionsResponse{}
}
func (m *GetConnectionsResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type ConnectionService struct {
	ServiceID string `json:"service_id"`
	Name      string `json:"name"`
}

func (m *ConnectionService) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetConnectionServicesResponse struct {
	ConnectionService []*ConnectionService `json:"connection_service"`
}

func (m *GetConnectionServicesResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type GetShopConnectionsResponse struct {
	ShopConnections []*ShopConnection `json:"shop_connections"`
}

func (m *GetShopConnectionsResponse) Reset() {
	*m = GetShopConnectionsResponse{}
}
func (m *GetShopConnectionsResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type DeleteShopConnectionRequest struct {
	ConnectionID dot.ID `json:"connection_id"`
}

func (m *DeleteShopConnectionRequest) Reset() {
	*m = DeleteShopConnectionRequest{}
}
func (m *DeleteShopConnectionRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type UpdateShopConnectionRequest struct {
	ConnectionID dot.ID                      `json:"connection_id"`
	Token        string                      `json:"token"`
	ExternalData *ShopConnectionExternalData `json:"external_data"`
}

func (m *UpdateShopConnectionRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreateTopshipConnectionRequest struct {
	ConnectionID dot.ID                      `json:"connection_id"`
	Token        string                      `json:"token"`
	ExternalData *ShopConnectionExternalData `json:"external_data"`
}

func (m *CreateTopshipConnectionRequest) String() string {
	return jsonx.MustMarshalToString(m)
}
