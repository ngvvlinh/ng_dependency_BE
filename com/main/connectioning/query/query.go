package query

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/connectioning/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ connectioning.QueryService = &ConnectionQuery{}

type ConnectionQuery struct {
	connectionStore     sqlstore.ConnectionStoreFactory
	shopConnectionStore sqlstore.ShopConnectionStoreFactory
}

func NewConnectionQuery(db *cmsql.Database) *ConnectionQuery {
	return &ConnectionQuery{
		connectionStore:     sqlstore.NewConnectionStore(db),
		shopConnectionStore: sqlstore.NewShopConnectionStore(db),
	}
}

func (q *ConnectionQuery) MessageBus() connectioning.QueryBus {
	b := bus.New()
	return connectioning.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *ConnectionQuery) GetConnectionByID(ctx context.Context, ID dot.ID) (*connectioning.Connection, error) {
	return q.connectionStore(ctx).ID(ID).GetConnection()
}

func (q *ConnectionQuery) ListConnections(ctx context.Context, args *connectioning.ListConnectionsArgs) ([]*connectioning.Connection, error) {
	query := q.connectionStore(ctx).ConnectionTypeOptional(args.ConnectionType).ConnectionMethodOptional(args.ConnectionMethod).ConnectionProviderOptional(args.ConnectionProvider)
	return query.ListConnections()
}

func (q *ConnectionQuery) GetShopConnectionByID(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).OptionalShopID(ShopID).ConnectionID(ConnectionID).GetShopConnection()
}

func (q *ConnectionQuery) ListShopConnections(ctx context.Context, _ *meta.Empty) ([]*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).ListShopConnections()
}

func (q *ConnectionQuery) ListGlobalShopConnections(ctx context.Context, _ *meta.Empty) ([]*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).IsGlobal(true).ListShopConnections()
}

func (q *ConnectionQuery) ListShopConnectionsByShopID(ctx context.Context, ShopID dot.ID) ([]*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).ShopID(ShopID).ListShopConnections()
}

func (q *ConnectionQuery) ListShopConnectionsByConnectionID(ctx context.Context, ConnectionID dot.ID) ([]*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).ConnectionID(ConnectionID).ListShopConnections()
}
