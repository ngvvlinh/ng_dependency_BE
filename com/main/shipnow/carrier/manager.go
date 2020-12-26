package carrier

import (
	"context"
	"sync"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/connectioning"
	connectiontypes "o.o/api/main/connectioning/types"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	shipnowtypes "o.o/api/main/shipnow/types"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	comshipnowcarriertypes "o.o/backend/com/main/shipnow/carrier/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ carrier.Manager = &ShipnowManager{}
var ll = l.New()

type ShipnowManager struct {
	env                string
	locationQS         location.QueryBus
	connectionQS       connectioning.QueryBus
	connectionAggr     connectioning.CommandBus
	redisStore         redis.Store
	connectionManager  *connectionmanager.ConnectionManager
	identityQS         identity.QueryBus
	shipnowQS          shipnow.QueryBus
	accountshipnowQS   accountshipnow.QueryBus
	accountshipnowAggr accountshipnow.CommandBus

	pathConfigs   comshipnowcarriertypes.Config
	carrierDriver comshipnowcarriertypes.Driver
}

func NewShipnowManager(locationQS location.QueryBus,
	connectionQS connectioning.QueryBus,
	connectionAggr connectioning.CommandBus,
	redisStore redis.Store,
	connectionManager *connectionmanager.ConnectionManager,
	identityQS identity.QueryBus,
	shipnowQS shipnow.QueryBus,
	accountshipnowQS accountshipnow.QueryBus,
	accountshipnowAggr accountshipnow.CommandBus,
	cfg comshipnowcarriertypes.Config,
	carrierDriver comshipnowcarriertypes.Driver,
) *ShipnowManager {
	return &ShipnowManager{
		env:                cmenv.PartnerEnv(),
		locationQS:         locationQS,
		connectionQS:       connectionQS,
		connectionAggr:     connectionAggr,
		redisStore:         redisStore,
		connectionManager:  connectionManager,
		identityQS:         identityQS,
		shipnowQS:          shipnowQS,
		accountshipnowQS:   accountshipnowQS,
		accountshipnowAggr: accountshipnowAggr,
		pathConfigs:        cfg,
		carrierDriver:      carrierDriver,
	}
}

type GetShipnowDriverArgs struct {
	ConnectionID dot.ID
	ShopID       dot.ID
	OwnerID      dot.ID
}

func (m *ShipnowManager) getShipnowDriver(ctx context.Context, args GetShipnowDriverArgs) (comshipnowcarriertypes.ShipnowCarrier, error) {
	connection, err := m.connectionManager.GetConnectionByID(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}
	if connection.Status != status3.P {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}

	if connection.ConnectionSubtype != connection_type.ConnectionSubtypeShipnow {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}

	getShopConnectionQuery := connectionmanager.GetShopConnectionArgs{
		ConnectionID: connection.ID,
	}
	if connection.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		// ignore ownerID
		getShopConnectionQuery.IsGlobal = true
	} else {
		_ownerID, err := m.GetOwnerID(ctx, GetOwnerIDArgs{
			OwnerID: args.OwnerID,
			ShopID:  args.ShopID,
		})
		if err != nil {
			return nil, err
		}
		if _ownerID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu owner_id")
		}
		getShopConnectionQuery.OwnerID = _ownerID
	}
	shopConnection, err := m.connectionManager.GetShopConnection(ctx, getShopConnectionQuery)
	if err != nil {
		return nil, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy shop connection")
	}
	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Shop connection không hợp lệ (vui lòng kiểm tra status hoặc token)")
	}
	return m.carrierDriver.GetShipnowDriver(
		m.env, m.locationQS,
		connection, shopConnection,
		m.identityQS,
		m.accountshipnowQS,
		m.pathConfigs,
	)
}

func (m *ShipnowManager) getShipnowDriverByEtopAffiliateAccount(ctx context.Context, connectionID dot.ID) (comshipnowcarriertypes.ShipnowCarrier, error) {
	conn, err := m.connectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	if conn.Status != status3.P {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}

	if conn.ConnectionSubtype != connection_type.ConnectionSubtypeShipnow {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
	return m.carrierDriver.GetAffiliateShipnowDriver(
		m.env, m.locationQS,
		conn,
		m.identityQS,
		m.accountshipnowQS,
	)
}

func (m *ShipnowManager) CreateExternalShipnow(ctx context.Context, args *carrier.CreateExternalShipnowCommand) (*carrier.ExternalShipnow, error) {
	// check balance of shop
	// if balance < MinShopBalance => can not create order

	query := &shipnow.GetShipnowFulfillmentQuery{
		ID:     args.ShipnowFulfillmentID,
		ShopID: args.ShopID,
	}
	if err := m.shipnowQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	shipnowFfm := query.Result.ShipnowFulfillment
	if err := m.ValidateShipnowAddress(shipnowFfm); err != nil {
		return nil, err
	}

	connID := cm.CoalesceID(shipnowFfm.ConnectionID, connectioning.DefaultTopShipAhamoveConnectionID)
	driver, err := m.getShipnowDriver(ctx, GetShipnowDriverArgs{
		ConnectionID: connID,
		ShopID:       args.ShopID,
		OwnerID:      args.OwnerID,
	})
	if err != nil {
		return nil, err
	}
	getShipnowServiceArgs := comshipnowcarriertypes.GetShipnowServiceArgs{
		ArbitraryID:    connID,
		ShopID:         args.ShopID,
		PickupAddress:  shipnowFfm.PickupAddress,
		DeliveryPoints: shipnowFfm.DeliveryPoints,
		Coupon:         args.Coupon,
	}
	services, err := driver.GetShipnowServices(ctx, getShipnowServiceArgs)
	if err != nil {
		return nil, err
	}

	service, err := m.CheckShippingService(shipnowFfm, driver, services)
	if err != nil {
		return nil, err
	}

	externalShipnow, err := driver.CreateExternalShipnow(ctx, args, service)
	if err != nil {
		return nil, err
	}
	externalShipnow.Service = service
	return externalShipnow, nil
}

func (m *ShipnowManager) CheckShippingService(ffm *shipnow.ShipnowFulfillment, driver comshipnowcarriertypes.ShipnowCarrier, services []*shipnowtypes.ShipnowService) (service *shipnowtypes.ShipnowService, err error) {
	if ffm.ShippingServiceCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần chọn gói dịch vụ giao hàng")
	}
	shippingServiceID, err := driver.ParseServiceCode(ffm.ShippingServiceCode)
	for _, s := range services {
		sID, _ := driver.ParseServiceCode(s.Code)
		if sID == shippingServiceID {
			service = s
			break
		}
	}
	if service == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Gói dịch vụ đã chọn không hợp lệ")
	}
	if ffm.ShippingServiceFee != service.Fee {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số tiền phí giao hàng không hợp lệ cho dịch vụ %v: Phí trên đơn hàng %v, phí từ dịch vụ giao hàng: %v", service.Name, ffm.ShippingServiceFee, service.Fee)
	}
	return service, nil
}

func (m *ShipnowManager) ValidateShipnowAddress(ffm *shipnow.ShipnowFulfillment) error {
	if _, _, _, err := m.VerifyAddress(ffm.PickupAddress, false); err != nil {
		return err
	}

	for _, point := range ffm.DeliveryPoints {
		if _, _, _, err := m.VerifyAddress(point.ShippingAddress, false); err != nil {
			return err
		}
	}
	return nil
}

func (m *ShipnowManager) VerifyAddress(addr *ordertypes.Address, requireWard bool) (*location.Province, *location.District, *location.Ward, error) {
	if addr == nil {
		return nil, nil, nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
	}
	query := &location.GetLocationQuery{
		ProvinceCode: addr.ProvinceCode,
		DistrictCode: addr.DistrictCode,
	}
	if requireWard {
		if addr.WardCode == "" {
			return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil,
				`Cần cung cấp thông tin phường xã hợp lệ`)
		}
		query.WardCode = addr.WardCode
	}
	if err := m.locationQS.DispatchAll(context.Background(), query); err != nil {
		return nil, nil, nil, err
	}
	if addr.Coordinates == nil || addr.Coordinates.Latitude == 0 || addr.Coordinates.Longitude == 0 {
		return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp Latitude và Longitude")
	}
	loc := query.Result
	return loc.Province, loc.District, loc.Ward, nil
}

func (m *ShipnowManager) CancelExternalShipping(ctx context.Context, args *carrier.CancelExternalShipnowCommand) error {
	// backward-compatible
	connID := cm.CoalesceID(args.ConnectionID, connectioning.DefaultTopShipAhamoveConnectionID)

	driver, err := m.getShipnowDriver(ctx, GetShipnowDriverArgs{
		ConnectionID: connID,
		ShopID:       args.ShopID,
		OwnerID:      args.OwnerID,
	})
	if err != nil {
		return err
	}
	if err := driver.CancelExternalShipnow(ctx, args); err != nil {
		return err
	}
	return nil
}

type GetOwnerIDArgs struct {
	OwnerID dot.ID
	ShopID  dot.ID
}

func (m *ShipnowManager) GetOwnerID(ctx context.Context, args GetOwnerIDArgs) (dot.ID, error) {
	if args.OwnerID != 0 {
		return args.OwnerID, nil
	}
	query := &identity.GetAccountByIDQuery{
		ID: args.ShopID,
	}
	if err := m.identityQS.Dispatch(ctx, query); err != nil {
		return 0, err
	}
	return query.Result.OwnerID, nil
}

func (m *ShipnowManager) GetExternalShipnowServices(ctx context.Context, args *carrier.GetExternalShipnowServicesCommand) ([]*shipnowtypes.ShipnowService, error) {
	var res []*shipnowtypes.ShipnowService

	shopConnections, err := m.GetAllShopConnections(ctx, args.ShopID, args.ConnectionIDs)
	if err != nil {
		return nil, err
	}

	getShipnowServicesArgs := comshipnowcarriertypes.GetShipnowServiceArgs{
		ShopID:         args.ShopID,
		PickupAddress:  args.PickupAddress,
		DeliveryPoints: args.DeliveryPoints,
		Coupon:         args.Coupon,
	}
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(shopConnections))
	for _, shopConn := range shopConnections {
		connID := shopConn.ConnectionID
		shopConn := shopConn
		go func() error {
			defer wg.Done()
			if shopConn.Status != status3.P || shopConn.Token == "" {
				return cm.Errorf(cm.FailedPrecondition, nil, "Connection shipnow does not valid (check status or token)")
			}
			getShipnowServicesArgs.ArbitraryID = connID
			services, err := m.getExternalShipnowServices(ctx, getShipnowServicesArgs, connID)
			if err != nil {
				ll.Error("GetShipnowService :: ", l.ID("connection_id", connID), l.Error(err))
				return err
			}
			mutex.Lock()
			res = append(res, services...)
			mutex.Unlock()
			return nil
		}()
	}
	wg.Wait()

	if len(res) == 0 {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Tuyến giao hàng không được hỗ trợ bởi đơn bị vận chuyển nào")
	}
	return res, nil
}

func (m *ShipnowManager) GetAllShopConnections(ctx context.Context, shopID dot.ID, connectionIDs []dot.ID) ([]*connectioning.ShopConnection, error) {
	// Get all shop_connection & global shop_connection
	if len(connectionIDs) == 0 {
		connQuery := &connectioning.ListConnectionsQuery{
			Status:            status3.P.Wrap(),
			ConnectionType:    connection_type.Shipping,
			ConnectionSubtype: connection_type.ConnectionSubtypeShipnow,
		}
		if err := m.connectionQS.Dispatch(ctx, connQuery); err != nil {
			return nil, err
		}
		for _, conn := range connQuery.Result {
			connectionIDs = append(connectionIDs, conn.ID)
		}
	}
	ownerID, err := m.GetOwnerID(ctx, GetOwnerIDArgs{
		ShopID: shopID,
	})
	if err != nil {
		return nil, err
	}
	query := &connectioning.ListShopConnectionsQuery{
		OwnerID:       ownerID,
		IncludeGlobal: true,
		ConnectionIDs: connectionIDs,
	}

	if err := m.connectionQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func (m *ShipnowManager) getExternalShipnowServices(ctx context.Context, args comshipnowcarriertypes.GetShipnowServiceArgs, connID dot.ID) ([]*shipnowtypes.ShipnowService, error) {
	conn, err := m.connectionManager.GetConnectionByID(ctx, connID)
	if err != nil {
		return nil, err
	}
	if err := validateConnection(conn); err != nil {
		return nil, err
	}

	driver, err := m.getShipnowDriver(ctx, GetShipnowDriverArgs{
		ConnectionID: conn.ID,
		ShopID:       args.ShopID,
		OwnerID:      args.OwnerID,
	})
	if err != nil {
		return nil, err
	}
	services, err := driver.GetShipnowServices(ctx, args)
	if err != nil {
		return nil, err
	}
	for _, s := range services {
		s.ConnectionInfo = &connectiontypes.ConnectionInfo{
			ID:       conn.ID,
			Name:     conn.Name,
			ImageURL: conn.ImageURL,
		}
	}
	return services, nil
}

func validateConnection(conn *connectioning.Connection) error {
	if conn.Status != status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Connection chưa được xác nhận")
	}
	if conn.ConnectionSubtype != connection_type.ConnectionSubtypeShipnow {
		return cm.Errorf(cm.FailedPrecondition, nil, "Loại Connection không hợp lệ. (Yêu cầu: %v, thực tế: %v", connection_type.ConnectionSubtypeShipnow.String(), conn.ConnectionSubtype.String())
	}
	return nil
}

func (m *ShipnowManager) RegisterExternalAccount(ctx context.Context, cmd *carrier.RegisterExternalAccountCommand) (*carrier.RegisterExternalAccountResult, error) {
	if cmd.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin số điện thoại")
	}
	if cmd.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin tên tài khoản")
	}

	connectionID := carrier.GetConnectionID(cmd.ConnectionID, cmd.Carrier, connection_type.ConnectionMethodDirect)
	driver, err := m.getShipnowDriverByEtopAffiliateAccount(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	args := &comshipnowcarriertypes.RegisterExternalAccountArgs{
		Phone:   cmd.Phone,
		Name:    cmd.Name,
		Address: cmd.Address,
	}
	return driver.RegisterExternalAccount(ctx, args)
}

func (m *ShipnowManager) GetExternalAccount(ctx context.Context, cmd *carrier.GetExternalAccountCommand) (*carrier.ExternalAccount, error) {
	connectionID := carrier.GetConnectionID(cmd.ConnectionID, cmd.Carrier, connection_type.ConnectionMethodDirect)
	driver, err := m.getShipnowDriver(ctx, GetShipnowDriverArgs{
		ConnectionID: connectionID,
		ShopID:       cmd.ShopID,
		OwnerID:      cmd.OwnerID,
	})
	if err != nil {
		return nil, err
	}
	args := &comshipnowcarriertypes.GetExternalAccountArgs{}
	return driver.GetExternalAccount(ctx, args)
}

func (m *ShipnowManager) VerifyExternalAccount(ctx context.Context, cmd *carrier.VerifyExternalAccountCommand) (*carrier.VerifyExternalAccountResult, error) {
	connectionID := carrier.GetConnectionID(cmd.ConnectionID, cmd.Carrier, connection_type.ConnectionMethodDirect)
	driver, err := m.getShipnowDriver(ctx, GetShipnowDriverArgs{
		ConnectionID: connectionID,
		ShopID:       cmd.ShopID,
		OwnerID:      cmd.OwnerID,
	})
	if err != nil {
		return nil, err
	}
	args := &comshipnowcarriertypes.VerifyExternalAccountArgs{
		OwnerID: cmd.OwnerID,
	}
	return driver.VerifyExternalAccount(ctx, args)
}

func (m *ShipnowManager) RefreshToken(ctx context.Context, args *carrier.RefreshTokenArgs) error {
	if args.ConnectionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Shipnow refresh token error: missing connection ID")
	}
	ownerID, err := m.GetOwnerID(ctx, GetOwnerIDArgs{
		OwnerID: args.OwnerID,
		ShopID:  args.ShopID,
	})
	if err != nil {
		return err
	}
	queryOwner := &identity.GetUserByIDQuery{
		UserID: ownerID,
	}
	if err := m.identityQS.Dispatch(ctx, queryOwner); err != nil {
		return err
	}
	user := queryOwner.Result

	registerCmd := &carrier.RegisterExternalAccountCommand{
		Phone:        user.Phone,
		Name:         user.FullName,
		ConnectionID: args.ConnectionID,
		OwnerID:      ownerID,
	}
	xAccount, err := m.RegisterExternalAccount(ctx, registerCmd)
	if err != nil {
		return err
	}

	// update shop_connection
	update := &connectioning.UpdateShopConnectionCommand{
		OwnerID:      ownerID,
		ConnectionID: args.ConnectionID,
		Token:        xAccount.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			Identifier: user.Phone,
		},
	}
	if err := m.connectionAggr.Dispatch(ctx, update); err != nil {
		return err
	}
	return nil
}
