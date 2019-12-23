package aggregate

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/connectioning/convert"
	"etop.vn/backend/com/main/connectioning/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	etopmodel "etop.vn/backend/pkg/etop/model"
	etopsqlstore "etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi/dot"
)

var _ connectioning.Aggregate = &ConnectionAggregate{}
var schemas = conversion.Build(convert.RegisterConversions)

type ConnectionAggregate struct {
	connectionStore     sqlstore.ConnectionStoreFactory
	shopConnectionStore sqlstore.ShopConnectionStoreFactory
}

func NewConnectionAggregate(db *cmsql.Database) *ConnectionAggregate {
	return &ConnectionAggregate{
		connectionStore:     sqlstore.NewConnectionStore(db),
		shopConnectionStore: sqlstore.NewShopConnectionStore(db),
	}
}

func (a *ConnectionAggregate) MessageBus() connectioning.CommandBus {
	b := bus.New()
	return connectioning.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *ConnectionAggregate) CreateConnection(ctx context.Context, args *connectioning.CreateConnectionArgs) (*connectioning.Connection, error) {
	if args.Driver == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Driver")
	}
	var conn connectioning.Connection
	if err := schemas.Convert(args, &conn); err != nil {
		return nil, err
	}
	conn.ID = cm.NewID()
	code, err := etopsqlstore.GenerateCodeWithoutTransaction(ctx, etopmodel.CodeTypeConnection, "")
	if err != nil {
		return nil, err
	}
	conn.Code = code
	return a.connectionStore(ctx).CreateConnection(&conn)
}

func (a *ConnectionAggregate) UpdateConnectionDriverConfig(ctx context.Context, args *connectioning.UpdateConnectionDriveConfig) (*connectioning.Connection, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	return a.connectionStore(ctx).UpdateConnectionDriverConfig(args)
}

func (a *ConnectionAggregate) ConfirmConnection(ctx context.Context, ID dot.ID) (updated int, err error) {
	conn, err := a.connectionStore(ctx).ID(ID).GetConnection()
	if err != nil {
		return 0, err
	}
	if conn.Status != status3.Z {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this Connection")
	}
	return a.connectionStore(ctx).ConfirmConnection(ID)
}

func (a *ConnectionAggregate) DeleteConnection(ctx context.Context, ID dot.ID) (deleted int, err error) {
	return a.connectionStore(ctx).ID(ID).SoftDelete()
}

func (a *ConnectionAggregate) CreateShopConnection(ctx context.Context, args *connectioning.CreateShopConnectionArgs) (*connectioning.ShopConnection, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	var shopConn connectioning.ShopConnection
	if err := schemas.Convert(args, &shopConn); err != nil {
		return nil, err
	}
	// always set status = 1
	shopConn.Status = 1
	return a.shopConnectionStore(ctx).CreateShopConnection(&shopConn)
}

func (a *ConnectionAggregate) UpdateShopConnectionToken(ctx context.Context, args *connectioning.UpdateShopConnectionExternalDataArgs) (*connectioning.ShopConnection, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	if args.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Token")
	}
	return a.shopConnectionStore(ctx).UpdateShopConnectionToken(args)
}

func (a *ConnectionAggregate) CreateOrUpdateShopConnection(ctx context.Context, args *connectioning.CreateShopConnectionArgs) (*connectioning.ShopConnection, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	if args.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Token")
	}

	conn, err := a.shopConnectionStore(ctx).ShopID(args.ShopID).ConnectionID(args.ConnectionID).GetShopConnection()
	if err == nil {
		// Update
		update := &connectioning.UpdateShopConnectionExternalDataArgs{
			ShopID:         conn.ShopID,
			ConnectionID:   conn.ConnectionID,
			Token:          args.Token,
			TokenExpiresAt: args.TokenExpiresAt,
			ExternalData:   args.ExternalData,
		}
		return a.shopConnectionStore(ctx).UpdateShopConnectionToken(update)
	}

	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	// Create
	cmd := &connectioning.CreateShopConnectionArgs{
		ShopID:         args.ShopID,
		ConnectionID:   args.ConnectionID,
		Token:          args.Token,
		TokenExpiresAt: args.TokenExpiresAt,
		ExternalData:   args.ExternalData,
	}
	return a.CreateShopConnection(ctx, cmd)

}

func (a *ConnectionAggregate) ConfirmShopConnection(ctx context.Context, shopID dot.ID, connectionID dot.ID) (updated int, err error) {
	conn, err := a.shopConnectionStore(ctx).ShopID(shopID).ConnectionID(connectionID).GetShopConnection()
	if err != nil {
		return 0, err
	}
	if conn.Status != status3.Z {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this Connection")
	}
	return a.shopConnectionStore(ctx).ConfirmShopConnection(shopID, connectionID)
}

func (a *ConnectionAggregate) DeleteShopConnection(ctx context.Context, shopID dot.ID, connectionID dot.ID) (deleted int, err error) {
	return a.shopConnectionStore(ctx).ShopID(shopID).ConnectionID(connectionID).SoftDelete()
}
