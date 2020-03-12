package connectioning

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	// -- Connection -- //
	CreateConnection(context.Context, *CreateConnectionArgs) (*Connection, error)

	UpdateConnection(context.Context, *UpdateConnectionArgs) (*Connection, error)

	UpdateConnectionAffiliateAccount(context.Context, *UpdateConnectionAffiliateAccountArgs) (updated int, err error)

	ConfirmConnection(ctx context.Context, ID dot.ID) (updated int, err error)

	DisableConnection(ctx context.Context, ID dot.ID) (updated int, err error)

	DeleteConnection(context.Context, *DeleteConnectionArgs) (deleted int, _ error)

	CreateTopshipConnection(context.Context, *CreateTopshipConnectionArgs) (*Connection, error)

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

	GetConnectionByCode(ctx context.Context, code string) (*Connection, error)

	ListConnections(context.Context, *ListConnectionsArgs) ([]*Connection, error)

	ListConnectionServicesByID(ctx context.Context, ID dot.ID) ([]*ConnectionService, error)

	// -- Shop Connection -- //
	GetShopConnectionByID(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (*ShopConnection, error)

	ListShopConnections(context.Context, *ListShopConnectionsArgs) ([]*ShopConnection, error)

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

// +convert:update=Connection(PartnerID)
type UpdateConnectionArgs struct {
	ID           dot.ID
	PartnerID    dot.ID
	Name         string
	ImageURL     string
	DriverConfig *ConnectionDriverConfig
}

type DeleteConnectionArgs struct {
	ID        dot.ID
	PartnerID dot.ID
}

type CreateTopshipConnectionArgs struct {
	ID           dot.ID
	Token        string
	ExternalData *ShopConnectionExternalData
}

// +convert:update=Connection
type UpdateConnectionAffiliateAccountArgs struct {
	ID                   dot.ID
	EtopAffiliateAccount *EtopAffiliateAccount
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
	PartnerID          dot.ID
	Status             status3.NullStatus
	ConnectionType     connection_type.ConnectionType
	ConnectionMethod   connection_type.ConnectionMethod
	ConnectionProvider connection_type.ConnectionProvider
}

type ListShopConnectionsArgs struct {
	ShopID        dot.ID
	IncludeGlobal bool
	ConnectionIDs []dot.ID
}
