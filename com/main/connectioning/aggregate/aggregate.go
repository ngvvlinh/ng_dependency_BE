package aggregate

import (
	"context"
	"fmt"

	"o.o/api/main/connectioning"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/connectioning/convert"
	"o.o/backend/com/main/connectioning/model"
	"o.o/backend/com/main/connectioning/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	etopmodel "o.o/backend/pkg/etop/model"
	etopsqlstore "o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ connectioning.Aggregate = &ConnectionAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type ConnectionAggregate struct {
	db                  *cmsql.Database
	txDB                cmsql.Transactioner
	connectionStore     sqlstore.ConnectionStoreFactory
	shopConnectionStore sqlstore.ShopConnectionStoreFactory
	eventBus            capi.EventBus
}

func NewConnectionAggregate(db com.MainDB, eventBus capi.EventBus) *ConnectionAggregate {
	return &ConnectionAggregate{
		db:                  db,
		txDB:                (*cmsql.Database)(db),
		eventBus:            eventBus,
		connectionStore:     sqlstore.NewConnectionStore(db),
		shopConnectionStore: sqlstore.NewShopConnectionStore(db),
	}
}

func ConnectionAggregateMessageBus(a *ConnectionAggregate) connectioning.CommandBus {
	b := bus.New()
	return connectioning.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *ConnectionAggregate) CreateConnection(ctx context.Context, args *connectioning.CreateConnectionArgs) (*connectioning.Connection, error) {
	args.Driver = fmt.Sprintf("%v/%v/%v/%v", args.ConnectionType, args.ConnectionSubtype, args.ConnectionMethod, args.ConnectionProvider)
	if args.ConnectionMethod == connection_type.ConnectionMethodDirect {
		if err := validateDirectConnection(args.DriverConfig); err != nil {
			return nil, err
		}
	}

	var conn connectioning.Connection
	if err := scheme.Convert(args, &conn); err != nil {
		return nil, err
	}
	conn.ID = cm.NewID()
	code, err := etopsqlstore.GenerateCodeWithoutTransaction(ctx, a.db, etopmodel.CodeTypeConnection, "")
	if err != nil {
		return nil, err
	}
	conn.Code = code
	// default status Z
	conn.Status = status3.Z

	return a.connectionStore(ctx).CreateConnection(&conn)
}

func validateDirectConnection(driverCfg *connectioning.ConnectionDriverConfig) error {
	if driverCfg == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp thông tin URL")
	}
	validateFields := []cmapi.Field{
		{
			Name:  "tracking_url",
			Value: driverCfg.TrackingURL,
		}, {
			Name:  "create_fulfillment_url",
			Value: driverCfg.CreateFulfillmentURL,
		}, {
			Name:  "get_fulfillment_url",
			Value: driverCfg.GetFulfillmentURL,
		}, {
			Name:  "cancel_fulfillment_url",
			Value: driverCfg.CancelFulfillmentURL,
		}, {
			Name:  "get_shipping_services_url",
			Value: driverCfg.GetShippingServicesURL,
		},
	}
	if err := cmapi.ValidateEmptyField(validateFields...); err != nil {
		return err
	}
	return nil
}

func (a *ConnectionAggregate) UpdateConnection(ctx context.Context, args *connectioning.UpdateConnectionArgs) (*connectioning.Connection, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	conn, err := a.connectionStore(ctx).ID(args.ID).OptionalPartnerID(args.PartnerID).GetConnection()
	if err != nil {
		return nil, err
	}
	if conn.ConnectionMethod == connection_type.ConnectionMethodDirect {
		if err := validateDirectConnection(args.DriverConfig); err != nil {
			return nil, err
		}
	}

	var res *connectioning.Connection
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var update connectioning.Connection
		if err := scheme.Convert(args, &update); err != nil {
			return err
		}
		query := a.connectionStore(ctx)
		if args.IgnoreWLPartner {
			query = query.IgnoreWLPartner()
		}
		res, err = a.connectionStore(ctx).UpdateConnection(&update)
		if err != nil {
			return err
		}
		return a.raiseConnectionUpdatedEvent(ctx, args.ID)
	})
	return res, err
}

func (a *ConnectionAggregate) ConfirmConnection(ctx context.Context, id dot.ID) (updated int, err error) {
	conn, err := a.connectionStore(ctx).ID(id).GetConnection()
	if err != nil {
		return 0, err
	}
	if conn.Status == status3.P {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "This Connection was confirmed")
	}
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.connectionStore(ctx).ConfirmConnection(id)
		if err != nil {
			return err
		}

		// enable/confirm all shop_connection
		if _, err := a.shopConnectionStore(ctx).ConnectionID(id).ConfirmShopConnection(); err != nil {
			return err
		}
		return a.raiseConnectionUpdatedEvent(ctx, id)
	})
	return
}

func (a *ConnectionAggregate) DeleteConnection(ctx context.Context, args *connectioning.DeleteConnectionArgs) (deleted int, err error) {
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		deleted, err = a.connectionStore(ctx).ID(args.ID).OptionalPartnerID(args.PartnerID).SoftDelete()
		if err != nil {
			return err
		}

		// disable all shop_connection
		if _, err := a.shopConnectionStore(ctx).ConnectionID(args.ID).DisableShopConnection(); err != nil {
			return err
		}

		return a.raiseConnectionUpdatedEvent(ctx, args.ID)
	})
	return
}

func (a *ConnectionAggregate) UpdateConnectionAffiliateAccount(ctx context.Context, args *connectioning.UpdateConnectionAffiliateAccountArgs) (updated int, err error) {
	conn, err := a.connectionStore(ctx).ID(args.ID).GetConnection()
	if err != nil {
		return 0, err
	}
	if conn.Status != status3.P {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Please confirm this connection first.")
	}
	if args.EtopAffiliateAccount == nil {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Affiliate account can't be empty.")
	}
	if args.EtopAffiliateAccount.Token == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing token affiliate account.")
	}
	if args.EtopAffiliateAccount.UserID == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing userID affiliate account.")
	}

	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var update connectioning.Connection
		if err := scheme.Convert(args, &update); err != nil {
			return err
		}
		if _, err := a.connectionStore(ctx).UpdateConnection(&update); err != nil {
			return err
		}
		return a.raiseConnectionUpdatedEvent(ctx, args.ID)
	})
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *ConnectionAggregate) DisableConnection(ctx context.Context, id dot.ID) (updated int, err error) {
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.connectionStore(ctx).DisableConnection(id)
		if err != nil {
			return err
		}

		// disable all shop_connection
		if _, err := a.shopConnectionStore(ctx).ConnectionID(id).DisableShopConnection(); err != nil {
			return err
		}

		return a.raiseConnectionUpdatedEvent(ctx, id)
	})
	return
}

func (a *ConnectionAggregate) UpdateConnectionFromOrigin(ctx context.Context, connectionID dot.ID) error {
	conn, err := a.connectionStore(ctx).ID(connectionID).GetConnection()
	if err != nil {
		return err
	}
	// chỉ quan tâm tới connection direct (NVC trực tiếp tích hợp)
	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil
	}

	childConns, err := a.connectionStore(ctx).OriginConnectionID(connectionID).IgnoreWLPartner().ListConnections(status3.NullStatus{})
	if err != nil {
		return err
	}

	// Lấy tất cả các connection có origin_connection_id = conn.ID
	// thay đổi thông tin của các conn đó theo thông tin của connection gốc
	for _, _conn := range childConns {
		update := &connectioning.UpdateConnectionArgs{
			ID:              _conn.ID,
			ImageURL:        conn.ImageURL,
			DriverConfig:    conn.DriverConfig,
			IgnoreWLPartner: true,
		}
		if _, err := a.UpdateConnection(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

func (a *ConnectionAggregate) raiseConnectionUpdatedEvent(ctx context.Context, connectionID dot.ID) error {
	event := &connectioning.ConnectionUpdatedEvent{
		EventMeta:    meta.NewEvent(),
		ConnectionID: connectionID,
	}
	return a.eventBus.Publish(ctx, event)
}

/*
	CreateTopshipConnection

	Xác nhận 1 direct connection thành builtin connection
	Cần cung cấp thêm thông tin tài khoản của topship đối với connection (token, externalData) đó để tạo shop_connection global tương ứng
*/
func (a *ConnectionAggregate) CreateBuiltinConnection(ctx context.Context, args *connectioning.CreateBuiltinConnectionArgs) (*connectioning.Connection, error) {
	conn, err := a.connectionStore(ctx).ID(args.ID).GetConnection()
	if err != nil {
		return nil, err
	}
	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid")
	}
	if args.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
	}
	if args.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing token")
	}
	if args.ExternalData == nil || args.ExternalData.UserID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing userID")
	}

	var result *connectioning.Connection
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		cmd := &connectioning.CreateConnectionArgs{
			Name:               args.Name,
			ConnectionType:     connection_type.Shipping,
			ConnectionSubtype:  connection_type.ConnectionSubtypeShipment,
			ConnectionMethod:   connection_type.ConnectionMethodBuiltin,
			ConnectionProvider: connection_type.ConnectionProviderPartner,
			DriverConfig:       conn.DriverConfig,
			PartnerID:          conn.PartnerID,
			ImageURL:           conn.ImageURL,
			Services:           conn.Services,
			OriginConnectionID: args.ID,
		}
		res, err := a.CreateConnection(ctx, cmd)
		if err != nil {
			return err
		}
		result = res

		// Tạo shop_connection global cho topship
		cmd2 := &model.ShopConnection{
			ConnectionID: res.ID,
			Token:        args.Token,
			Status:       status3.P,
			IsGlobal:     true,
			ExternalData: &model.ShopConnectionExternalData{
				UserID:     args.ExternalData.UserID,
				Identifier: args.ExternalData.Identifier,
				ShopID:     args.ExternalData.ShopID,
			},
		}

		if err := a.shopConnectionStore(ctx).CreateShopConnectionDB(cmd2); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ConnectionAggregate) CreateShopConnection(ctx context.Context, args *connectioning.CreateShopConnectionArgs) (*connectioning.ShopConnection, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	conn, err := a.connectionStore(ctx).ID(args.ConnectionID).GetConnection()
	if err != nil {
		return nil, err
	}

	var shopConn connectioning.ShopConnection
	if err := scheme.Convert(args, &shopConn); err != nil {
		return nil, err
	}
	if conn.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		shopConn.IsGlobal = true
	} else {
		if args.ShopID == 0 && args.OwnerID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID or OwnerID")
		}
	}

	// always set status = 1
	shopConn.Status = 1
	return a.shopConnectionStore(ctx).CreateShopConnection(&shopConn)
}

func (a *ConnectionAggregate) UpdateShopConnectionToken(ctx context.Context, args *connectioning.UpdateShopConnectionExternalDataArgs) (*connectioning.ShopConnection, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	if args.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Token")
	}
	cmd := &connectioning.CreateShopConnectionArgs{
		ShopID:       args.ShopID,
		OwnerID:      args.OwnerID,
		ConnectionID: args.ConnectionID,
		Token:        args.Token,
		ExternalData: args.ExternalData,
	}
	return a.CreateOrUpdateShopConnection(ctx, cmd)
}

func (a *ConnectionAggregate) CreateOrUpdateShopConnection(ctx context.Context, args *connectioning.CreateShopConnectionArgs) (res *connectioning.ShopConnection, err error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	if args.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Token")
	}
	conn, err := a.connectionStore(ctx).ID(args.ConnectionID).GetConnection()
	if err != nil {
		return nil, err
	}
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		if args.ShopID == 0 && args.OwnerID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID or OwnerID")
		}
	}

	shopConn, err := a.shopConnectionStore(ctx).OptionalShopID(args.ShopID).OptionalOwnerID(args.OwnerID).ConnectionID(args.ConnectionID).GetShopConnection()
	if err == nil {
		// Update
		err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
			update := &connectioning.UpdateShopConnectionExternalDataArgs{
				OwnerID:        shopConn.OwnerID,
				ShopID:         shopConn.ShopID,
				ConnectionID:   shopConn.ConnectionID,
				Token:          args.Token,
				TokenExpiresAt: args.TokenExpiresAt,
				ExternalData:   args.ExternalData,
			}
			res, err = a.shopConnectionStore(ctx).UpdateShopConnectionToken(update)
			if err != nil {
				return err
			}
			return a.raiseShopConnectionUpdatedEvent(ctx, shopConn.ShopID, shopConn.OwnerID, shopConn.ConnectionID)
		})
		return
	}

	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	// Create
	cmd := &connectioning.CreateShopConnectionArgs{
		OwnerID:        args.OwnerID,
		ShopID:         args.ShopID,
		ConnectionID:   args.ConnectionID,
		Token:          args.Token,
		TokenExpiresAt: args.TokenExpiresAt,
		ExternalData:   args.ExternalData,
	}
	return a.CreateShopConnection(ctx, cmd)

}

func (a *ConnectionAggregate) ConfirmShopConnection(ctx context.Context, args *connectioning.ShopConnectionQueryArgs) (updated int, err error) {
	shopID, ownerID, connectionID := args.ShopID, args.OwnerID, args.ConnectionID
	if shopID == 0 && ownerID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin shop ID hoặc owner ID")
	}
	shopConn, err := a.shopConnectionStore(ctx).ConnectionID(connectionID).
		OptionalShopID(shopID).
		OptionalOwnerID(ownerID).
		GetShopConnection()
	if err != nil {
		return 0, err
	}
	if shopConn.Status != status3.Z {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this Shop Connection")
	}
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.shopConnectionStore(ctx).ConnectionID(connectionID).
			OptionalShopID(shopID).
			OptionalOwnerID(ownerID).
			ConfirmShopConnection()
		if err != nil {
			return err
		}
		return a.raiseShopConnectionUpdatedEvent(ctx, shopID, ownerID, connectionID)
	})
	return
}

func (a *ConnectionAggregate) DeleteShopConnection(ctx context.Context, args *connectioning.ShopConnectionQueryArgs) (deleted int, err error) {
	shopID, ownerID, connectionID := args.ShopID, args.OwnerID, args.ConnectionID
	if shopID == 0 && ownerID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin shop ID hoặc owner ID")
	}
	conn, err := a.connectionStore(ctx).ID(connectionID).GetConnection()
	if err != nil {
		return 0, err
	}
	err = a.txDB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		switch conn.ConnectionSubtype {
		case connection_type.ConnectionSubtypeShipment:
			// Shop connection unique theo shop_id
			deleted, err = a.shopConnectionStore(ctx).ConnectionID(connectionID).
				ShopID(shopID).
				SoftDelete()
			if err != nil {
				return err
			}
		case connection_type.ConnectionSubtypeShipnow:
			// // Shop connection unique theo owner_id
			deleted, err = a.shopConnectionStore(ctx).ConnectionID(connectionID).
				OwnerID(ownerID).
				SoftDelete()
			if err != nil {
				return err
			}
		}

		if err = a.raiseShopConnectionUpdatedEvent(ctx, shopID, ownerID, connectionID); err != nil {
			return err
		}
		return nil
	})
	return
}

func (a *ConnectionAggregate) raiseShopConnectionUpdatedEvent(ctx context.Context, shopID, ownerID, connectionID dot.ID) error {
	event := &connectioning.ShopConnectionUpdatedEvent{
		EventMeta:    meta.NewEvent(),
		ShopID:       shopID,
		OwnerID:      ownerID,
		ConnectionID: connectionID,
	}
	return a.eventBus.Publish(ctx, event)
}
