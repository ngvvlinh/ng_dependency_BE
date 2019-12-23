package query

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/connectioning/sqlstore"
	cm "etop.vn/backend/pkg/common"
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

func (q *ConnectionQuery) GetConnectionByCode(ctx context.Context, code string) (*connectioning.Connection, error) {
	return q.connectionStore(ctx).Code(code).GetConnection()
}

func (q *ConnectionQuery) ListConnections(ctx context.Context, args *connectioning.ListConnectionsArgs) ([]*connectioning.Connection, error) {
	query := q.connectionStore(ctx).ConnectionTypeOptional(args.ConnectionType).ConnectionMethodOptional(args.ConnectionMethod).ConnectionProviderOptional(args.ConnectionProvider)
	return query.ListConnections()
}

func (q *ConnectionQuery) GetShopConnectionByID(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).OptionalShopID(ShopID).ConnectionID(ConnectionID).GetShopConnection()
}

func (q *ConnectionQuery) ListShopConnections(ctx context.Context, args *connectioning.ListShopConnectionsArgs) ([]*connectioning.ShopConnection, error) {
	query := q.shopConnectionStore(ctx)
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ListShopConnections failed. Invalid ShopID")
	}
	if len(args.ConnectionIDs) > 0 {
		query = query.ConnectionIDs(args.ConnectionIDs...)
	}

	var res []*connectioning.ShopConnection
	{
		query1 := query.Clone()
		res1, err := query1.ShopID(args.ShopID).ListShopConnections()
		if err != nil {
			return nil, err
		}
		res = append(res, res1...)
	}
	if args.IncludeGlobal {
		query2 := query.Clone()
		res2, err := query2.IsGlobal(true).ListShopConnections()
		if err != nil {
			return nil, err
		}
		res = append(res, res2...)
	}
	return res, nil
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
