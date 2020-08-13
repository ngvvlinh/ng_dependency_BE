package carrier

import (
	"context"
	"sync"

	"o.o/api/main/connectioning"
	connectiontypes "o.o/api/main/connectioning/types"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	shipnowcarriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	"o.o/backend/com/main/shipnowcarrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	ahamoveclient "o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ carrier.Manager = &ShipnowManager{}
var ll = l.New()

type ShipnowManager struct {
	env               string
	locationQS        location.QueryBus
	connectionQS      connectioning.QueryBus
	redisStore        redis.Store
	connectionManager *connectionmanager.ConnectionManager
	identityQS        identity.QueryBus
	shipnowQS         shipnow.QueryBus
}

func NewShipnowManager(locationQS location.QueryBus,
	connectionQS connectioning.QueryBus,
	redisStore redis.Store,
	connectionManager *connectionmanager.ConnectionManager,
	identityQS identity.QueryBus,
	shipnowQS shipnow.QueryBus,
) *ShipnowManager {
	return &ShipnowManager{
		env:               cmenv.PartnerEnv(),
		locationQS:        locationQS,
		connectionQS:      connectionQS,
		redisStore:        redisStore,
		connectionManager: connectionManager,
		identityQS:        identityQS,
		shipnowQS:         shipnowQS,
	}
}

func (m *ShipnowManager) getShipnowDriver(ctx context.Context, connectionID dot.ID, shopID dot.ID) (*shipnowcarrier.Carrier, error) {
	connection, err := m.connectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	if connection.Status != status3.P {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}

	if connection.ConnectionSubtype != connection_type.ConnectionSubtypeShipnow {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
	_shopID := shopID
	if connection.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		// ignore shopID
		_shopID = 0
	}
	shopConnection, err := m.connectionManager.GetShopConnection(ctx, connectionID, _shopID)
	if err != nil {
		return nil, cm.Errorf(cm.ErrorCode(err), err, "Không tìm thấy shop connection")
	}
	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Shop connection không hợp lệ (vui lòng kiểm tra status hoặc token)")
	}
	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderAhamove:
		cfg := ahamoveclient.Config{
			Env:   m.env,
			Name:  "",
			Token: shopConnection.Token,
		}
		client := ahamoveclient.New(cfg)
		driver := ahamove.New(client, m.locationQS, m.identityQS)
		return &shipnowcarrier.Carrier{
			Code:                  shipnowcarriertypes.Ahamove,
			ShipnowCarrier:        driver,
			ShipnowCarrierAccount: nil,
		}, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
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
	driver, err := m.getShipnowDriver(ctx, connID, shipnowFfm.ShopID)
	if err != nil {
		return nil, nil
	}
	getShipnowServiceArgs := shipnowcarrier.GetShipnowServiceArgs{
		ShopID:         args.ShopID,
		PickupAddress:  shipnowFfm.PickupAddress,
		DeliveryPoints: shipnowFfm.DeliveryPoints,
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

func (m *ShipnowManager) CheckShippingService(ffm *shipnow.ShipnowFulfillment, driver *shipnowcarrier.Carrier, services []*shipnowtypes.ShipnowService) (service *shipnowtypes.ShipnowService, err error) {
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

	driver, err := m.getShipnowDriver(ctx, connID, args.ShopID)
	if err != nil {
		return err
	}
	if err := driver.CancelExternalShipnow(ctx, args); err != nil {
		return err
	}
	return nil
}

func (m *ShipnowManager) GetExternalShipnowServices(ctx context.Context, args *carrier.GetExternalShipnowServicesCommand) ([]*shipnowtypes.ShipnowService, error) {
	var res []*shipnowtypes.ShipnowService

	shopConnections, err := m.GetAllShopConnections(ctx, args.ShopID, args.ConnectionIDs)
	if err != nil {
		return nil, err
	}

	getShipnowServicesArgs := shipnowcarrier.GetShipnowServiceArgs{
		ShopID:         args.ShopID,
		PickupAddress:  args.PickupAddress,
		DeliveryPoints: args.DeliveryPoints,
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
			services, err := m.getExternalShipnowService(ctx, getShipnowServicesArgs, connID)
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
	query := &connectioning.ListShopConnectionsQuery{
		ShopID:        shopID,
		IncludeGlobal: true,
		ConnectionIDs: connectionIDs,
	}

	if err := m.connectionQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func (m *ShipnowManager) getExternalShipnowService(ctx context.Context, args shipnowcarrier.GetShipnowServiceArgs, connID dot.ID) ([]*shipnowtypes.ShipnowService, error) {
	conn, err := m.connectionManager.GetConnectionByID(ctx, connID)
	if err != nil {
		return nil, err
	}
	if err := validateConnection(conn); err != nil {
		return nil, err
	}

	driver, err := m.getShipnowDriver(ctx, connID, args.ShopID)
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
	panic("implement me")
}

func (m *ShipnowManager) GetExternalAccount(ctx context.Context, cmd *carrier.GetExternalAccountCommand) (*carrier.ExternalAccount, error) {
	panic("implement me")
}

func (m *ShipnowManager) VerifyExternalAccount(ctx context.Context, cmd *carrier.VerifyExternalAccountCommand) (*carrier.VerifyExternalAccountResult, error) {
	panic("implement me")
}
