package connectioning

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	// -- Connection -- //
	CreateConnection(context.Context, *CreateConnectionArgs) (*Connection, error)

	UpdateConnectionDriverConfig(context.Context, *UpdateConnectionDriveConfig) (*Connection, error)

	ConfirmConnection(ctx context.Context, ID dot.ID) (updated int, err error)

	DeleteConnection(ctx context.Context, ID dot.ID) (deleted int, _ error)

	// -- Shop Connection -- //
	CreateShopConnection(context.Context, *CreateShopConnectionArgs) (*ShopConnection, error)

	CreateOrUpdateShopConnection(context.Context, *CreateShopConnectionArgs) (*ShopConnection, error)

	UpdateShopConnectionToken(context.Context, *UpdateShopConnectionExternalDataArgs) (*ShopConnection, error)

	ConfirmShopConnection(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (updated int, err error)

	DeleteShopConnection(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (deleted int, _ error)
}

type QueryService interface {
	// -- Connection -- //
	GetConnectionByID(ctx context.Context, ID dot.ID) (*Connection, error)

	ListConnections(context.Context, *ListConnectionsArgs) ([]*Connection, error)

	// -- Shop Connection -- //
	GetShopConnectionByID(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (*ShopConnection, error)

	ListShopConnections(context.Context, *meta.Empty) ([]*ShopConnection, error)

	ListGlobalShopConnections(context.Context, *meta.Empty) ([]*ShopConnection, error)

	ListShopConnectionsByShopID(ctx context.Context, ShopID dot.ID) ([]*ShopConnection, error)

	ListShopConnectionsByConnectionID(ctx context.Context, ConnectionID dot.ID) ([]*ShopConnection, error)
}

// +convert:create=Connection
type CreateConnectionArgs struct {
	Name               string
	PartnerID          dot.ID
	Driver             string
	DriverConfig       *ConnectionDriverConfig
	ConnectionType     connection_type.ConnectionType
	ConnectionSubtype  connection_type.ConnectionSubtype
	ConnectionMethod   connection_type.ConnectionMethod
	ConnectionProvider connection_type.ConnectionProvider
}

type UpdateConnectionDriveConfig struct {
	ConnectionID dot.ID
	DriverConfig *ConnectionDriverConfig
}

// +convert:create=ShopConnection
type CreateShopConnectionArgs struct {
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData
}

// +convert:update=ShopConnection(ShopID,ConnectionID)
type UpdateShopConnectionExternalDataArgs struct {
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData
}

type ListConnectionsArgs struct {
	ConnectionType     connection_type.ConnectionType
	ConnectionMethod   connection_type.ConnectionMethod
	ConnectionProvider connection_type.ConnectionProvider
}
