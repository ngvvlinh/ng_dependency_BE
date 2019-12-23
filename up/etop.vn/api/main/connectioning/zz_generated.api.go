// +build !generator

// Code generated by generator api. DO NOT EDIT.

package connectioning

import (
	context "context"
	time "time"

	meta "etop.vn/api/meta"
	connection_type "etop.vn/api/top/types/etc/connection_type"
	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type ConfirmConnectionCommand struct {
	ID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleConfirmConnection(ctx context.Context, msg *ConfirmConnectionCommand) (err error) {
	msg.Result, err = h.inner.ConfirmConnection(msg.GetArgs(ctx))
	return err
}

type ConfirmShopConnectionCommand struct {
	ShopID       dot.ID
	ConnectionID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleConfirmShopConnection(ctx context.Context, msg *ConfirmShopConnectionCommand) (err error) {
	msg.Result, err = h.inner.ConfirmShopConnection(msg.GetArgs(ctx))
	return err
}

type CreateConnectionCommand struct {
	Name               string
	PartnerID          dot.ID
	Driver             string
	DriverConfig       *ConnectionDriverConfig
	ConnectionType     connection_type.ConnectionType
	ConnectionSubtype  connection_type.ConnectionSubtype
	ConnectionMethod   connection_type.ConnectionMethod
	ConnectionProvider connection_type.ConnectionProvider

	Result *Connection `json:"-"`
}

func (h AggregateHandler) HandleCreateConnection(ctx context.Context, msg *CreateConnectionCommand) (err error) {
	msg.Result, err = h.inner.CreateConnection(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateShopConnectionCommand struct {
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData

	Result *ShopConnection `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateShopConnection(ctx context.Context, msg *CreateOrUpdateShopConnectionCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateShopConnection(msg.GetArgs(ctx))
	return err
}

type CreateShopConnectionCommand struct {
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData

	Result *ShopConnection `json:"-"`
}

func (h AggregateHandler) HandleCreateShopConnection(ctx context.Context, msg *CreateShopConnectionCommand) (err error) {
	msg.Result, err = h.inner.CreateShopConnection(msg.GetArgs(ctx))
	return err
}

type DeleteConnectionCommand struct {
	ID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteConnection(ctx context.Context, msg *DeleteConnectionCommand) (err error) {
	msg.Result, err = h.inner.DeleteConnection(msg.GetArgs(ctx))
	return err
}

type DeleteShopConnectionCommand struct {
	ShopID       dot.ID
	ConnectionID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteShopConnection(ctx context.Context, msg *DeleteShopConnectionCommand) (err error) {
	msg.Result, err = h.inner.DeleteShopConnection(msg.GetArgs(ctx))
	return err
}

type UpdateConnectionDriverConfigCommand struct {
	ConnectionID dot.ID
	DriverConfig *ConnectionDriverConfig

	Result *Connection `json:"-"`
}

func (h AggregateHandler) HandleUpdateConnectionDriverConfig(ctx context.Context, msg *UpdateConnectionDriverConfigCommand) (err error) {
	msg.Result, err = h.inner.UpdateConnectionDriverConfig(msg.GetArgs(ctx))
	return err
}

type UpdateShopConnectionTokenCommand struct {
	ShopID         dot.ID
	ConnectionID   dot.ID
	Token          string
	TokenExpiresAt time.Time
	ExternalData   *ShopConnectionExternalData

	Result *ShopConnection `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopConnectionToken(ctx context.Context, msg *UpdateShopConnectionTokenCommand) (err error) {
	msg.Result, err = h.inner.UpdateShopConnectionToken(msg.GetArgs(ctx))
	return err
}

type GetConnectionByCodeQuery struct {
	Code string

	Result *Connection `json:"-"`
}

func (h QueryServiceHandler) HandleGetConnectionByCode(ctx context.Context, msg *GetConnectionByCodeQuery) (err error) {
	msg.Result, err = h.inner.GetConnectionByCode(msg.GetArgs(ctx))
	return err
}

type GetConnectionByIDQuery struct {
	ID dot.ID

	Result *Connection `json:"-"`
}

func (h QueryServiceHandler) HandleGetConnectionByID(ctx context.Context, msg *GetConnectionByIDQuery) (err error) {
	msg.Result, err = h.inner.GetConnectionByID(msg.GetArgs(ctx))
	return err
}

type GetShopConnectionByIDQuery struct {
	ShopID       dot.ID
	ConnectionID dot.ID

	Result *ShopConnection `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopConnectionByID(ctx context.Context, msg *GetShopConnectionByIDQuery) (err error) {
	msg.Result, err = h.inner.GetShopConnectionByID(msg.GetArgs(ctx))
	return err
}

type ListConnectionsQuery struct {
	ConnectionType     connection_type.ConnectionType
	ConnectionMethod   connection_type.ConnectionMethod
	ConnectionProvider connection_type.ConnectionProvider

	Result []*Connection `json:"-"`
}

func (h QueryServiceHandler) HandleListConnections(ctx context.Context, msg *ListConnectionsQuery) (err error) {
	msg.Result, err = h.inner.ListConnections(msg.GetArgs(ctx))
	return err
}

type ListGlobalShopConnectionsQuery struct {
	Result []*ShopConnection `json:"-"`
}

func (h QueryServiceHandler) HandleListGlobalShopConnections(ctx context.Context, msg *ListGlobalShopConnectionsQuery) (err error) {
	msg.Result, err = h.inner.ListGlobalShopConnections(msg.GetArgs(ctx))
	return err
}

type ListShopConnectionsQuery struct {
	ShopID        dot.ID
	IncludeGlobal bool
	ConnectionIDs []dot.ID

	Result []*ShopConnection `json:"-"`
}

func (h QueryServiceHandler) HandleListShopConnections(ctx context.Context, msg *ListShopConnectionsQuery) (err error) {
	msg.Result, err = h.inner.ListShopConnections(msg.GetArgs(ctx))
	return err
}

type ListShopConnectionsByConnectionIDQuery struct {
	ConnectionID dot.ID

	Result []*ShopConnection `json:"-"`
}

func (h QueryServiceHandler) HandleListShopConnectionsByConnectionID(ctx context.Context, msg *ListShopConnectionsByConnectionIDQuery) (err error) {
	msg.Result, err = h.inner.ListShopConnectionsByConnectionID(msg.GetArgs(ctx))
	return err
}

type ListShopConnectionsByShopIDQuery struct {
	ShopID dot.ID

	Result []*ShopConnection `json:"-"`
}

func (h QueryServiceHandler) HandleListShopConnectionsByShopID(ctx context.Context, msg *ListShopConnectionsByShopIDQuery) (err error) {
	msg.Result, err = h.inner.ListShopConnectionsByShopID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *ConfirmConnectionCommand) command()            {}
func (q *ConfirmShopConnectionCommand) command()        {}
func (q *CreateConnectionCommand) command()             {}
func (q *CreateOrUpdateShopConnectionCommand) command() {}
func (q *CreateShopConnectionCommand) command()         {}
func (q *DeleteConnectionCommand) command()             {}
func (q *DeleteShopConnectionCommand) command()         {}
func (q *UpdateConnectionDriverConfigCommand) command() {}
func (q *UpdateShopConnectionTokenCommand) command()    {}

func (q *GetConnectionByCodeQuery) query()               {}
func (q *GetConnectionByIDQuery) query()                 {}
func (q *GetShopConnectionByIDQuery) query()             {}
func (q *ListConnectionsQuery) query()                   {}
func (q *ListGlobalShopConnectionsQuery) query()         {}
func (q *ListShopConnectionsQuery) query()               {}
func (q *ListShopConnectionsByConnectionIDQuery) query() {}
func (q *ListShopConnectionsByShopIDQuery) query()       {}

// implement conversion

func (q *ConfirmConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *ConfirmShopConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID, ConnectionID dot.ID) {
	return ctx,
		q.ShopID,
		q.ConnectionID
}

func (q *CreateConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateConnectionArgs) {
	return ctx,
		&CreateConnectionArgs{
			Name:               q.Name,
			PartnerID:          q.PartnerID,
			Driver:             q.Driver,
			DriverConfig:       q.DriverConfig,
			ConnectionType:     q.ConnectionType,
			ConnectionSubtype:  q.ConnectionSubtype,
			ConnectionMethod:   q.ConnectionMethod,
			ConnectionProvider: q.ConnectionProvider,
		}
}

func (q *CreateConnectionCommand) SetCreateConnectionArgs(args *CreateConnectionArgs) {
	q.Name = args.Name
	q.PartnerID = args.PartnerID
	q.Driver = args.Driver
	q.DriverConfig = args.DriverConfig
	q.ConnectionType = args.ConnectionType
	q.ConnectionSubtype = args.ConnectionSubtype
	q.ConnectionMethod = args.ConnectionMethod
	q.ConnectionProvider = args.ConnectionProvider
}

func (q *CreateOrUpdateShopConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShopConnectionArgs) {
	return ctx,
		&CreateShopConnectionArgs{
			ShopID:         q.ShopID,
			ConnectionID:   q.ConnectionID,
			Token:          q.Token,
			TokenExpiresAt: q.TokenExpiresAt,
			ExternalData:   q.ExternalData,
		}
}

func (q *CreateOrUpdateShopConnectionCommand) SetCreateShopConnectionArgs(args *CreateShopConnectionArgs) {
	q.ShopID = args.ShopID
	q.ConnectionID = args.ConnectionID
	q.Token = args.Token
	q.TokenExpiresAt = args.TokenExpiresAt
	q.ExternalData = args.ExternalData
}

func (q *CreateShopConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShopConnectionArgs) {
	return ctx,
		&CreateShopConnectionArgs{
			ShopID:         q.ShopID,
			ConnectionID:   q.ConnectionID,
			Token:          q.Token,
			TokenExpiresAt: q.TokenExpiresAt,
			ExternalData:   q.ExternalData,
		}
}

func (q *CreateShopConnectionCommand) SetCreateShopConnectionArgs(args *CreateShopConnectionArgs) {
	q.ShopID = args.ShopID
	q.ConnectionID = args.ConnectionID
	q.Token = args.Token
	q.TokenExpiresAt = args.TokenExpiresAt
	q.ExternalData = args.ExternalData
}

func (q *DeleteConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *DeleteShopConnectionCommand) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID, ConnectionID dot.ID) {
	return ctx,
		q.ShopID,
		q.ConnectionID
}

func (q *UpdateConnectionDriverConfigCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateConnectionDriveConfig) {
	return ctx,
		&UpdateConnectionDriveConfig{
			ConnectionID: q.ConnectionID,
			DriverConfig: q.DriverConfig,
		}
}

func (q *UpdateConnectionDriverConfigCommand) SetUpdateConnectionDriveConfig(args *UpdateConnectionDriveConfig) {
	q.ConnectionID = args.ConnectionID
	q.DriverConfig = args.DriverConfig
}

func (q *UpdateShopConnectionTokenCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShopConnectionExternalDataArgs) {
	return ctx,
		&UpdateShopConnectionExternalDataArgs{
			ShopID:         q.ShopID,
			ConnectionID:   q.ConnectionID,
			Token:          q.Token,
			TokenExpiresAt: q.TokenExpiresAt,
			ExternalData:   q.ExternalData,
		}
}

func (q *UpdateShopConnectionTokenCommand) SetUpdateShopConnectionExternalDataArgs(args *UpdateShopConnectionExternalDataArgs) {
	q.ShopID = args.ShopID
	q.ConnectionID = args.ConnectionID
	q.Token = args.Token
	q.TokenExpiresAt = args.TokenExpiresAt
	q.ExternalData = args.ExternalData
}

func (q *GetConnectionByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, code string) {
	return ctx,
		q.Code
}

func (q *GetConnectionByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetShopConnectionByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID, ConnectionID dot.ID) {
	return ctx,
		q.ShopID,
		q.ConnectionID
}

func (q *ListConnectionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListConnectionsArgs) {
	return ctx,
		&ListConnectionsArgs{
			ConnectionType:     q.ConnectionType,
			ConnectionMethod:   q.ConnectionMethod,
			ConnectionProvider: q.ConnectionProvider,
		}
}

func (q *ListConnectionsQuery) SetListConnectionsArgs(args *ListConnectionsArgs) {
	q.ConnectionType = args.ConnectionType
	q.ConnectionMethod = args.ConnectionMethod
	q.ConnectionProvider = args.ConnectionProvider
}

func (q *ListGlobalShopConnectionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *ListGlobalShopConnectionsQuery) SetEmpty(args *meta.Empty) {
}

func (q *ListShopConnectionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListShopConnectionsArgs) {
	return ctx,
		&ListShopConnectionsArgs{
			ShopID:        q.ShopID,
			IncludeGlobal: q.IncludeGlobal,
			ConnectionIDs: q.ConnectionIDs,
		}
}

func (q *ListShopConnectionsQuery) SetListShopConnectionsArgs(args *ListShopConnectionsArgs) {
	q.ShopID = args.ShopID
	q.IncludeGlobal = args.IncludeGlobal
	q.ConnectionIDs = args.ConnectionIDs
}

func (q *ListShopConnectionsByConnectionIDQuery) GetArgs(ctx context.Context) (_ context.Context, ConnectionID dot.ID) {
	return ctx,
		q.ConnectionID
}

func (q *ListShopConnectionsByShopIDQuery) GetArgs(ctx context.Context) (_ context.Context, ShopID dot.ID) {
	return ctx,
		q.ShopID
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleConfirmConnection)
	b.AddHandler(h.HandleConfirmShopConnection)
	b.AddHandler(h.HandleCreateConnection)
	b.AddHandler(h.HandleCreateOrUpdateShopConnection)
	b.AddHandler(h.HandleCreateShopConnection)
	b.AddHandler(h.HandleDeleteConnection)
	b.AddHandler(h.HandleDeleteShopConnection)
	b.AddHandler(h.HandleUpdateConnectionDriverConfig)
	b.AddHandler(h.HandleUpdateShopConnectionToken)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetConnectionByCode)
	b.AddHandler(h.HandleGetConnectionByID)
	b.AddHandler(h.HandleGetShopConnectionByID)
	b.AddHandler(h.HandleListConnections)
	b.AddHandler(h.HandleListGlobalShopConnections)
	b.AddHandler(h.HandleListShopConnections)
	b.AddHandler(h.HandleListShopConnectionsByConnectionID)
	b.AddHandler(h.HandleListShopConnectionsByShopID)
	return QueryBus{b}
}
