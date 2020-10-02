package connectioning

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
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

	CreateBuiltinConnection(context.Context, *CreateBuiltinConnectionArgs) (*Connection, error)

	UpdateConnectionFromOrigin(ctx context.Context, ConnectionID dot.ID) error

	// -- Shop Connection -- //
	CreateShopConnection(context.Context, *CreateShopConnectionArgs) (*ShopConnection, error)

	CreateOrUpdateShopConnection(context.Context, *CreateShopConnectionArgs) (*ShopConnection, error)

	UpdateShopConnectionToken(context.Context, *UpdateShopConnectionExternalDataArgs) (*ShopConnection, error)

	ConfirmShopConnection(context.Context, *ShopConnectionQueryArgs) (updated int, err error)

	DeleteShopConnection(context.Context, *ShopConnectionQueryArgs) (deleted int, _ error)
}

type QueryService interface {
	// -- Connection -- //
	GetConnectionByID(ctx context.Context, ID dot.ID) (*Connection, error)

	GetConnectionByCode(ctx context.Context, code string) (*Connection, error)

	ListConnections(context.Context, *ListConnectionsArgs) ([]*Connection, error)

	ListConnectionServicesByID(ctx context.Context, ID dot.ID) ([]*ConnectionService, error)

	ListConnectionsByOriginConnectionID(ctx context.Context, OriginConnectionID dot.ID) ([]*Connection, error)

	// -- Shop Connection -- //
	GetShopConnection(context.Context, *GetShopConnectionArgs) (*ShopConnection, error)

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
	ImageURL           string
	Services           []*ConnectionService
	OriginConnectionID dot.ID
}

// +convert:update=Connection(PartnerID)
type UpdateConnectionArgs struct {
	ID              dot.ID
	PartnerID       dot.ID
	Name            string
	ImageURL        string
	Services        []*ConnectionService
	DriverConfig    *ConnectionDriverConfig
	IgnoreWLPartner bool
}

type DeleteConnectionArgs struct {
	ID        dot.ID
	PartnerID dot.ID
}

type CreateBuiltinConnectionArgs struct {
	ID           dot.ID
	Name         string
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
	OwnerID        dot.ID
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData
}

// +convert:update=ShopConnection(ShopID,ConnectionID)
type UpdateShopConnectionExternalDataArgs struct {
	OwnerID        dot.ID
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData
}

type ListConnectionsArgs struct {
	IDs                []dot.ID
	PartnerID          dot.ID
	Status             status3.NullStatus
	ConnectionType     connection_type.ConnectionType
	ConnectionSubtype  connection_type.ConnectionSubtype
	ConnectionMethod   connection_type.ConnectionMethod
	ConnectionProvider connection_type.ConnectionProvider
}

type ListShopConnectionsArgs struct {
	ShopID        dot.ID
	OwnerID       dot.ID
	IncludeGlobal bool
	ConnectionIDs []dot.ID
}

type GetShopConnectionArgs struct {
	ShopID       dot.ID
	OwnerID      dot.ID
	ConnectionID dot.ID
	// when IsGlobal = true => ignore ShopID & OwnerID
	IsGlobal bool
}

type ShopConnectionQueryArgs struct {
	OwnerID      dot.ID
	ShopID       dot.ID
	ConnectionID dot.ID
}
