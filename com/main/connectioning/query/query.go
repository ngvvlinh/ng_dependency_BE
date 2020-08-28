package query

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
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
	mapServices         shippingservices.MapShipmentServices
}

func NewConnectionQuery(db com.MainDB, mapServices shippingservices.MapShipmentServices) *ConnectionQuery {
	return &ConnectionQuery{
		connectionStore:     sqlstore.NewConnectionStore(db),
		shopConnectionStore: sqlstore.NewShopConnectionStore(db),
		mapServices:         mapServices,
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
	// backward-compatible
	// set default subtype to shipment
	if args.ConnectionSubtype == 0 {
		args.ConnectionSubtype = connection_type.ConnectionSubtypeShipment
	}

	query := q.connectionStore(ctx).OptionalPartnerID(args.PartnerID).OptionalConnectionType(args.ConnectionType).OptionalConnectionSubtype(args.ConnectionSubtype).OptionalConnectionMethod(args.ConnectionMethod).OptionalConnectionProvider(args.ConnectionProvider)
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
		services := q.mapServices.ByCarrier(carrier)
		for _, s := range services {
			res = append(res, &connectioning.ConnectionService{
				ServiceID: s.ServiceID,
				Name:      s.Name,
			})
		}
	}
	return res, nil
}

func (q *ConnectionQuery) ListConnectionsByOriginConnectionID(ctx context.Context, originConnectionID dot.ID) ([]*connectioning.Connection, error) {
	return q.connectionStore(ctx).OriginConnectionID(originConnectionID).OptionalConnectionMethod(connection_type.ConnectionMethodBuiltin).ListConnections(status3.NullStatus{Valid: false})
}

func (q *ConnectionQuery) GetShopConnection(ctx context.Context, args *connectioning.GetShopConnectionArgs) (*connectioning.ShopConnection, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection_id")
	}
	query := q.shopConnectionStore(ctx).ConnectionID(args.ConnectionID)

	if args.IsGlobal {
		return query.IsGlobal(args.IsGlobal).GetShopConnection()
	}

	if args.ShopID == 0 && args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id or owner_id")
	}
	if args.ShopID != 0 && args.OwnerID != 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Only provide either shop_id or owner_id")
	}
	if args.ShopID != 0 {
		query = query.ShopID(args.ShopID)
	}
	if args.OwnerID != 0 {
		query = query.OwnerID(args.OwnerID)
	}
	return query.GetShopConnection()
}

func (q *ConnectionQuery) ListShopConnections(ctx context.Context, args *connectioning.ListShopConnectionsArgs) ([]*connectioning.ShopConnection, error) {
	// connection subtype;
	// - shipment: shop_id != null, owner_id = null
	// - shipnow: shop_id = null, owner_id != null
	if args.ShopID == 0 && args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID or OwnerID")
	}
	query := q.shopConnectionStore(ctx)
	if len(args.ConnectionIDs) > 0 {
		query = query.ConnectionIDs(args.ConnectionIDs...)
	}

	var res []*connectioning.ShopConnection
	if args.OwnerID != 0 {
		query1 := query.Clone()
		query1 = query1.IsGlobal(false)
		query1 = query1.OwnerID(args.OwnerID)
		res1, err := query1.ListShopConnections()
		if err != nil {
			return nil, err
		}
		res = append(res, res1...)
	}
	if args.ShopID != 0 {
		query2 := query.Clone()
		query2 = query2.IsGlobal(false)
		query2 = query2.ShopID(args.ShopID)
		res2, err := query2.ListShopConnections()
		if err != nil {
			return nil, err
		}
		res = append(res, res2...)
	}

	if args.IncludeGlobal {
		query3 := query.Clone()
		res3, err := query3.IsGlobal(true).ListShopConnections()
		if err != nil {
			return nil, err
		}
		res = append(res, res3...)
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
