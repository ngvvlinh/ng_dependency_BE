package query

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipping_provider"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/connectioning/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	shippingservices "o.o/backend/pkg/integration/shipping/services"
	"o.o/capi/dot"
)

var _ connectioning.QueryService = &ConnectionQuery{}

type ConnectionQuery struct {
	connectionStore     sqlstore.ConnectionStoreFactory
	shopConnectionStore sqlstore.ShopConnectionStoreFactory
}

func NewConnectionQuery(db com.MainDB) *ConnectionQuery {
	return &ConnectionQuery{
		connectionStore:     sqlstore.NewConnectionStore(db),
		shopConnectionStore: sqlstore.NewShopConnectionStore(db),
	}
}

func ConnectionQueryMessageBus(q *ConnectionQuery) connectioning.QueryBus {
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
	query := q.connectionStore(ctx).OptionalPartnerID(args.PartnerID).OptionalConnectionType(args.ConnectionType).OptionalConnectionMethod(args.ConnectionMethod).OptionalConnectionProvider(args.ConnectionProvider)
	return query.ListConnections(args.Status)
}

func (q *ConnectionQuery) ListConnectionServicesByID(ctx context.Context, id dot.ID) ([]*connectioning.ConnectionService, error) {
	conn, err := q.connectionStore(ctx).ID(id).GetConnection()
	if err != nil {
		return nil, err
	}
	if conn.Services != nil && len(conn.Services) > 0 {
		return conn.Services, nil
	}

	// Get default service in case of shipment
	var res = []*connectioning.ConnectionService{}
	if conn.ConnectionType == connection_type.Shipping {
		carrier, ok := shipping_provider.ParseShippingProvider(conn.ConnectionProvider.Name())
		if !ok {
			return res, nil
		}
		services := shippingservices.GetServicesByCarrier(carrier)
		for _, s := range services {
			res = append(res, &connectioning.ConnectionService{
				ServiceID: s.ServiceID,
				Name:      s.Name,
			})
		}
	}
	return res, nil
}

func (q *ConnectionQuery) GetShopConnectionByID(ctx context.Context, ShopID dot.ID, ConnectionID dot.ID) (*connectioning.ShopConnection, error) {
	return q.shopConnectionStore(ctx).OptionalShopID(ShopID).ConnectionID(ConnectionID).GetShopConnection()
}

func (q *ConnectionQuery) ListShopConnections(ctx context.Context, args *connectioning.ListShopConnectionsArgs) ([]*connectioning.ShopConnection, error) {
	query := q.shopConnectionStore(ctx)
	if args.ShopID == 0 && !args.IncludeGlobal {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ListShopConnections failed. Invalid ShopID")
	}
	if len(args.ConnectionIDs) > 0 {
		query = query.ConnectionIDs(args.ConnectionIDs...)
	}

	var res []*connectioning.ShopConnection
	if args.ShopID != 0 {
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
