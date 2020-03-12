package aggregate

import (
	"context"
	"fmt"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/connectioning/convert"
	"etop.vn/backend/com/main/connectioning/model"
	"etop.vn/backend/com/main/connectioning/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	etopmodel "etop.vn/backend/pkg/etop/model"
	etopsqlstore "etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ connectioning.Aggregate = &ConnectionAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type ConnectionAggregate struct {
	db                  cmsql.Transactioner
	connectionStore     sqlstore.ConnectionStoreFactory
	shopConnectionStore sqlstore.ShopConnectionStoreFactory
	eventBus            capi.EventBus
}

func NewConnectionAggregate(db *cmsql.Database, eventBus capi.EventBus) *ConnectionAggregate {
	return &ConnectionAggregate{
		db:                  db,
		eventBus:            eventBus,
		connectionStore:     sqlstore.NewConnectionStore(db),
		shopConnectionStore: sqlstore.NewShopConnectionStore(db),
	}
}

func (a *ConnectionAggregate) MessageBus() connectioning.CommandBus {
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
	code, err := etopsqlstore.GenerateCodeWithoutTransaction(ctx, etopmodel.CodeTypeConnection, "")
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
	var update connectioning.Connection
	if err := scheme.Convert(args, &update); err != nil {
		return nil, err
	}
	res, err := a.connectionStore(ctx).UpdateConnection(&update)
	if err != nil {
		return nil, err
	}

	// raise event & ignore error
	event := &connectioning.ConnectionUpdatedEvent{
		EventMeta:    meta.NewEvent(),
		ConnectionID: args.ID,
	}
	_ = a.eventBus.Publish(ctx, event)
	return res, nil
}

func (a *ConnectionAggregate) ConfirmConnection(ctx context.Context, ID dot.ID) (updated int, err error) {
	conn, err := a.connectionStore(ctx).ID(ID).GetConnection()
	if err != nil {
		return 0, err
	}
	if conn.Status != status3.Z {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Can not confirm this Connection")
	}
	res, err := a.connectionStore(ctx).ConfirmConnection(ID)
	if err != nil {
		return 0, err
	}
	// raise event & ignore error
	event := &connectioning.ConnectionUpdatedEvent{
		EventMeta:    meta.NewEvent(),
		ConnectionID: ID,
	}
	_ = a.eventBus.Publish(ctx, event)
	return res, nil
}

func (a *ConnectionAggregate) DeleteConnection(ctx context.Context, args *connectioning.DeleteConnectionArgs) (deleted int, err error) {
	return a.connectionStore(ctx).ID(args.ID).OptionalPartnerID(args.PartnerID).SoftDelete()
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
	var update connectioning.Connection
	if err := scheme.Convert(args, &update); err != nil {
		return 0, err
	}
	if _, err := a.connectionStore(ctx).UpdateConnection(&update); err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *ConnectionAggregate) DisableConnection(ctx context.Context, id dot.ID) (updated int, err error) {
	res, err := a.connectionStore(ctx).DisableConnection(id)
	if err != nil {
		return 0, err
	}
	// raise event & ignore error
	event := &connectioning.ConnectionUpdatedEvent{
		EventMeta:    meta.NewEvent(),
		ConnectionID: id,
	}
	_ = a.eventBus.Publish(ctx, event)
	return res, nil
}

/*
	CreateTopshipConnection

	Xác nhận 1 direct connection thành topship connection
	Cần cung cấp thêm thông tin tài khoản của topship đối với connection (token, externalData) đó để tạo shop_connection global tương ứng
*/
func (a *ConnectionAggregate) CreateTopshipConnection(ctx context.Context, args *connectioning.CreateTopshipConnectionArgs) (*connectioning.Connection, error) {
	conn, err := a.connectionStore(ctx).ID(args.ID).GetConnection()
	if err != nil {
		return nil, err
	}
	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid")
	}
	if args.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing token")
	}
	if args.ExternalData == nil || args.ExternalData.UserID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing userID")
	}

	var result *connectioning.Connection
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		cmd := &connectioning.CreateConnectionArgs{
			Name:               "Topship - " + conn.Name,
			ConnectionType:     connection_type.Shipping,
			ConnectionSubtype:  connection_type.ConnectionSubtypeShipment,
			ConnectionMethod:   connection_type.ConnectionMethodTopship,
			ConnectionProvider: connection_type.ConnectionProviderPartner,
			DriverConfig:       conn.DriverConfig,
			PartnerID:          conn.PartnerID,
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
				UserID: args.ExternalData.UserID,
				Email:  args.ExternalData.Email,
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
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ShopID")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	var shopConn connectioning.ShopConnection
	if err := scheme.Convert(args, &shopConn); err != nil {
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
	cmd := &connectioning.CreateShopConnectionArgs{
		ShopID:       args.ShopID,
		ConnectionID: args.ConnectionID,
		Token:        args.Token,
		ExternalData: args.ExternalData,
	}
	return a.CreateOrUpdateShopConnection(ctx, cmd)
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
