package carrier

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/location"
	"etop.vn/api/main/ordering"
	"etop.vn/api/top/types/etc/connection_type"
	shippingstate "etop.vn/api/top/types/etc/shipping"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	addressmodel "etop.vn/backend/com/main/address/model"
	carriertypes "etop.vn/backend/com/main/shipping/carrier/types"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	shippingsharemodel "etop.vn/backend/com/main/shipping/sharemodel"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/syncgroup"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cipherx"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
	ghndriver "etop.vn/backend/pkg/integration/shipping/ghn/driver"
	ghtkclient "etop.vn/backend/pkg/integration/shipping/ghtk/client"
	ghtkdriver "etop.vn/backend/pkg/integration/shipping/ghtk/driver"
	vtpostdriver "etop.vn/backend/pkg/integration/shipping/vtpost/driver"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ll                 = l.New()
	GHNWebhookEndpoint string
)

const (
	MinShopBalance = -200000
	DefaultTTl     = 2 * 60 * 60
	SecretKey      = "connectionsecretkey"
	versionCaching = "v0.1"
)

type ShipmentManager struct {
	LocationQS     location.QueryBus
	ConnectionQS   connectioning.QueryBus
	connectionAggr connectioning.CommandBus
	Env            string
	driver         carriertypes.ShipmentCarrier
	redisStore     redis.Store
	cipherx        *cipherx.Cipherx
}

func NewShipmentManager(locationQS location.QueryBus, connectionQS connectioning.QueryBus, connectionAggr connectioning.CommandBus, redisS redis.Store) *ShipmentManager {
	_cipherx, _ := cipherx.NewCipherx(SecretKey)
	return &ShipmentManager{
		LocationQS:     locationQS,
		ConnectionQS:   connectionQS,
		connectionAggr: connectionAggr,
		Env:            cmenv.PartnerEnv(),
		redisStore:     redisS,
		cipherx:        _cipherx,
	}
}

func (m *ShipmentManager) SetDriver(driver carriertypes.ShipmentCarrier) {
	m.driver = driver
}

func (m *ShipmentManager) ResetDriver() {
	m.driver = nil
}

func (m *ShipmentManager) SetWebhookEndpoint(connectionProvider connection_type.ConnectionProvider, endpoint string) error {
	if connectionProvider == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "SetWebhookEndpoint: Missing connection provider")
	}
	if endpoint == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "SetWebhookEndpoint: Missing webhook endpoint (provider = %v)", connectionProvider)
	}
	switch connectionProvider {
	case connection_type.ConnectionProviderGHN:
		GHNWebhookEndpoint = endpoint
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "SetWebhookEndpoint: Do not support this provider (%v)", connectionProvider)
	}
	return nil
}

func (m *ShipmentManager) GetWebhookEndpoint(connectionProvider connection_type.ConnectionProvider) (string, error) {
	switch connectionProvider {
	case connection_type.ConnectionProviderGHN:
		return GHNWebhookEndpoint, nil
	default:
		return "", cm.Errorf(cm.InvalidArgument, nil, "GetWebhookEndpoint: Do not support this provider (%v)", connectionProvider)
	}
}

func (m *ShipmentManager) getShipmentDriver(ctx context.Context, connectionID dot.ID, shopID dot.ID) (carriertypes.ShipmentCarrier, error) {
	connection, err := m.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	etopAffiliateAccount := connection.EtopAffiliateAccount
	_shopID := shopID
	if connection.ConnectionMethod == connection_type.ConnectionMethodTopship {
		// ignore shopID
		_shopID = 0
	}
	shopConnection, err := m.getShopConnection(ctx, connectionID, _shopID)
	if err != nil {
		return nil, err
	}

	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid (check status or token)")
	}

	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderGHN:
		if shopConnection.ExternalData == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection ExternalData is missing (connection_id = %v, shop_id = %v)", connection.ID, shopConnection.ShopID)
		}
		userID := shopConnection.ExternalData.UserID
		clientID, err := strconv.Atoi(userID)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Connection ExternalData: UserID is invalid")
		}
		cfg := ghnclient.GHNAccountCfg{
			ClientID: clientID,
			Token:    shopConnection.Token,
		}
		if etopAffiliateAccount != nil {
			if affiliateID, err := strconv.Atoi(etopAffiliateAccount.UserID); err == nil {
				cfg.AffiliateID = affiliateID
			}
		}
		webhookEndpoint, err := m.GetWebhookEndpoint(connection_type.ConnectionProviderGHN)
		if err != nil {
			return nil, err
		}
		driver := ghndriver.New(m.Env, cfg, m.LocationQS, webhookEndpoint)
		return driver, nil

	case connection_type.ConnectionProviderGHTK:
		cfg := ghtkclient.GhtkAccount{
			Token: shopConnection.Token,
		}
		if etopAffiliateAccount != nil {
			cfg.AffiliateID = etopAffiliateAccount.UserID
			cfg.B2CToken = etopAffiliateAccount.Token
		}
		driver := ghtkdriver.New(m.Env, cfg, m.LocationQS)
		return driver, nil

	case connection_type.ConnectionProviderVTP:
		driver := vtpostdriver.New(m.Env, shopConnection.Token, m.LocationQS)
		return driver, nil

	case connection_type.ConnectionProviderPartner:
		return nil, cm.ErrTODO
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}

func (m *ShipmentManager) CreateFulfillments(ctx context.Context, order *ordering.Order, ffms []*shipmodel.Fulfillment) error {
	// check balance of shop
	// if balance < MinShopBalance => can not create order
	// TODO: raise event FulfillmentCreatingEvent after merge wallet (amount-available service)
	{
		query := &model.GetBalanceShopCommand{
			ShopID: order.ShopID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			return err
		}
		balance := query.Result.Amount
		if balance < MinShopBalance {
			return cm.Errorf(cm.FailedPrecondition, nil, "Bạn đang nợ cước số tiền: %v. Vui lòng liên hệ ETOP để xử lý.", balance)
		}
	}

	var err error
	g := syncgroup.New(len(ffms))
	for _, ffm := range ffms {
		ffm := ffm
		g.Go(func() error { return m.createSingleFulfillment(ctx, order, ffm) })
	}
	errs := g.Wait()
	if errs.IsAll() {
		err = errs[0]
	}
	return err
}

func (m *ShipmentManager) createSingleFulfillment(ctx context.Context, order *ordering.Order, ffm *shipmodel.Fulfillment) (_err error) {
	driver, err := m.getShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "invalid connection")
	}

	// UpdateInfo status to error
	defer func() {
		if _err == nil {
			return
		}
		updateFfm2 := &shipmodel.Fulfillment{
			ID:         ffm.ID,
			SyncStatus: status4.N,
			SyncStates: &shippingsharemodel.FulfillmentSyncStates{
				TrySyncAt: time.Now(),
				Error:     model.ToError(_err),

				NextShippingState: shippingstate.Created,
			},
		}
		cmd := &shipmodelx.UpdateFulfillmentCommand{Fulfillment: updateFfm2}

		// Keep the original error
		_ = bus.Dispatch(ctx, cmd)
	}()

	fromDistrict, _, err := m.VerifyDistrictCode(ffm.AddressFrom)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "FromDistrictCode: %v", err)
	}
	toDistrict, _, err := m.VerifyDistrictCode(ffm.AddressTo)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "ToDistrictCode: %v", err)
	}

	weight := ffm.GrossWeight
	if weight == 0 {
		weight = ffm.ChargeableWeight
	}

	args := &carriertypes.GetShippingServicesArgs{
		ArbitraryID:            ffm.ConnectionID,
		AccountID:              order.ShopID,
		FromDistrictCode:       fromDistrict.Code,
		ToDistrictCode:         toDistrict.Code,
		ChargeableWeight:       cm.CoalesceInt(weight, 100),
		Length:                 ffm.Length,
		Width:                  ffm.Width,
		Height:                 ffm.Height,
		IncludeInsurance:       ffm.IncludeInsurance,
		BasketValue:            ffm.BasketValue,
		CODAmount:              ffm.TotalCODAmount,
		IncludeTopshipServices: true,
	}

	allServices, err := driver.GetShippingServices(ctx, args)
	if err != nil {
		return err
	}

	// check if etop package
	var etopService, providerService *model.AvailableShippingService
	sType, isEtopService := etop_shipping_price.ParseEtopServiceCode(ffm.ProviderServiceID)
	if isEtopService {
		// ETOP serivce
		// => Get cheapest provider service
		etopService, err = GetEtopServiceFromSeviceCode(ffm.ProviderServiceID, ffm.ShippingServiceFee, allServices)
		if err != nil {
			return err
		}

		providerService = shipping_provider.GetCheapestService(allServices, sType)
		if providerService == nil {
			return cm.Error(cm.InvalidArgument, "Không có gói vận chuyển phù hợp.", nil)
		}
	} else {
		// Provider service
		// => Check price
		// => Get this service
		providerService, err = CheckShippingService(ffm, allServices)
		if err != nil {
			return err
		}
	}

	ffmToUpdate, err := driver.CreateFulfillment(ctx, ffm, args, providerService)
	if err != nil {
		return err
	}
	// update shipping service name
	ffmToUpdate.ShippingServiceName = providerService.Name

	if etopService != nil {
		err := ffmToUpdate.ApplyEtopPrice(etopService.ShippingFeeMain)
		if err != nil {
			return err
		}
		ffmToUpdate.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(ffmToUpdate.ProviderShippingFeeLines, ffmToUpdate.EtopPriceRule, dot.Int(ffmToUpdate.EtopAdjustedShippingFeeMain))
	}
	ffmToUpdate.ExternalAffiliateID = driver.GetAffiliateID()
	updateCmd := &shipmodelx.UpdateFulfillmentCommand{
		Fulfillment: ffmToUpdate,
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return cm.Trace(err)
	}
	return nil
}

func (m *ShipmentManager) VerifyDistrictCode(addr *addressmodel.Address) (*location.District, *location.Province, error) {
	if addr == nil {
		return nil, nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
	}
	if addr.DistrictCode == "" {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil,
			`Địa chỉ %v, %v không thể được xác định bởi hệ thống.`,
			addr.District, addr.Province,
		)
	}

	query := &location.GetLocationQuery{DistrictCode: addr.DistrictCode}
	if err := m.LocationQS.Dispatch(context.TODO(), query); err != nil {
		return nil, nil, err
	}
	district := query.Result.District
	return district, query.Result.Province, nil
}

func (m *ShipmentManager) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	driver, err := m.getShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "invalid connection")
	}
	return driver.CancelFulfillment(ctx, ffm)
}

func (m *ShipmentManager) GetShippingServices(ctx context.Context, accountID dot.ID, args *GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	shopConnections, err := m.GetAllShopConnections(ctx, accountID, args.ConnectionIDs)
	if err != nil {
		return nil, err
	}
	var res []*model.AvailableShippingService
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(shopConnections))
	for _, shopConn := range shopConnections {
		connID := shopConn.ConnectionID
		shopConn := shopConn
		go func() error {
			defer wg.Done()
			if shopConn.Status != status3.P || shopConn.Token == "" {
				return cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid (check status or token)")
			}
			conn, err := m.GetConnectionByID(ctx, connID)
			if err != nil {
				return err
			}
			if !m.validateConnection(ctx, conn) {
				return nil
			}

			var services []*model.AvailableShippingService
			driver, err := m.getShipmentDriver(ctx, connID, accountID)
			if err != nil {
				ll.Error("Driver shipment không hợp lệ", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
				return err
			}

			_args := args.ToShipmentServiceArgs() // clone the request to prevent race condition
			_args.ArbitraryID = connID
			_args.AccountID = accountID
			if conn.ConnectionMethod == connection_type.ConnectionMethodTopship {
				_args.IncludeTopshipServices = true
			}
			services, err = driver.GetShippingServices(ctx, _args)
			if err != nil {
				ll.Error("Get service error", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
				return err
			}

			// assign connection info to services
			for _, s := range services {
				s.ConnectionInfo = &model.ConnectionInfo{
					ID:       connID,
					Name:     conn.Name,
					ImageURL: conn.ImageURL,
				}
			}
			mutex.Lock()
			res = append(res, services...)
			mutex.Unlock()
			return nil
		}()
	}
	wg.Wait()

	if len(res) == 0 {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Không có gói giao hàng phù hợp")
	}
	res = shipping_provider.CompactServices(res)
	return res, nil
}

func (m *ShipmentManager) GetAllShopConnections(ctx context.Context, shopID dot.ID, connectionIDs []dot.ID) ([]*connectioning.ShopConnection, error) {
	// Get all shop_connection & global shop_connection
	query := &connectioning.ListShopConnectionsQuery{
		ShopID:        shopID,
		IncludeGlobal: true,
		ConnectionIDs: connectionIDs,
	}
	if err := m.ConnectionQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func (m *ShipmentManager) SignIn(ctx context.Context, args *ConnectionSignInArgs) (account *carriertypes.AccountResponse, _ error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}

	driver, err := m.getDriverByEtopAffiliateAccount(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}
	cmd := &carriertypes.SignInArgs{
		Email:    args.Email,
		Password: args.Password,
	}
	return driver.SignIn(ctx, cmd)
}

func (m *ShipmentManager) SignUp(ctx context.Context, args *ConnectionSignUpArgs) (newAccount *carriertypes.AccountResponse, _ error) {
	driver, err := m.getDriverByEtopAffiliateAccount(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}
	cmd := &carriertypes.SignUpArgs{
		Name:     args.Name,
		Email:    args.Email,
		Password: args.Password,
		Phone:    args.Phone,
		Province: args.Province,
		District: args.District,
		Address:  args.Address,
	}
	return driver.SignUp(ctx, cmd)
}

func (m *ShipmentManager) getDriverByEtopAffiliateAccount(ctx context.Context, connectionID dot.ID) (carriertypes.ShipmentCarrier, error) {
	conn, err := m.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}

	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Do not support this feature for this connection")
	}

	if conn.EtopAffiliateAccount == nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Missing EtopAffiliateAcount in Connection")
	}

	userID := conn.EtopAffiliateAccount.UserID
	token := conn.EtopAffiliateAccount.Token

	switch conn.ConnectionProvider {
	case connection_type.ConnectionProviderGHN:
		clientID, err := strconv.Atoi(userID)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "AffiliateAccount: UserID is invalid")
		}
		cfg := ghnclient.GHNAccountCfg{
			ClientID: clientID,
			Token:    token,
		}
		webhookEndpoint, err := m.GetWebhookEndpoint(connection_type.ConnectionProviderGHN)
		if err != nil {
			return nil, err
		}
		driver := ghndriver.New(m.Env, cfg, m.LocationQS, webhookEndpoint)
		return driver, nil
	case connection_type.ConnectionProviderGHTK:
		cfg := ghtkclient.GhtkAccount{
			AccountID: userID,
			Token:     token,
		}
		driver := ghtkdriver.New(m.Env, cfg, m.LocationQS)
		return driver, nil
	case connection_type.ConnectionProviderVTP:
		return nil, cm.Errorf(cm.Unimplemented, nil, "VTPost: không hỗ trợ affiliate account")
	case connection_type.ConnectionProviderPartner:
		return nil, cm.ErrTODO
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}

func (m *ShipmentManager) UpdateFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (updateFfm *shipmodel.Fulfillment, err error) {
	if ffm.ConnectionID == 0 && m.driver == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			"Không thể cập nhật fulfillment. Không tìm thấy driver hợp lệ (ffm_id = %v)", ffm.ID)
	}
	var driver carriertypes.ShipmentCarrier
	if ffm.ConnectionID != 0 {
		driver, err = m.getShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
		if err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "invalid connection (ffm_id = %v)", ffm.ID)
		}
	} else {
		driver = m.driver
	}

	updateFfm, err = driver.UpdateFulfillment(ctx, ffm)
	if err != nil {
		return nil, err
	}
	return
}

func (m *ShipmentManager) GetConnectionByID(ctx context.Context, connID dot.ID) (*connectioning.Connection, error) {
	connKey := getRedisConnectionKeyByID(connID)
	var connection connectioning.Connection
	err := m.loadRedis(connKey, &connection)
	if err != nil {
		query := &connectioning.GetConnectionByIDQuery{
			ID: connID,
		}
		if err := m.ConnectionQS.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).Wrap(cm.NotFound, "Connection not found").Throw()
		}
		connection = *query.Result
		connKeyCode := getRedisConnectionKeyByCode(connection.Code)
		m.setRedis(connKey, connection)
		m.setRedis(connKeyCode, connection)
	}
	return &connection, nil
}

func (m *ShipmentManager) GetConnectionByCode(ctx context.Context, connCode string) (*connectioning.Connection, error) {
	connKey := getRedisConnectionKeyByCode(connCode)
	var connection connectioning.Connection
	err := m.loadRedis(connKey, &connection)
	if err != nil {
		query := &connectioning.GetConnectionByCodeQuery{
			Code: connCode,
		}
		if err := m.ConnectionQS.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).Wrap(cm.NotFound, "Connection not found").Throw()
		}
		connection = *query.Result
		connKeyID := getRedisConnectionKeyByID(connection.ID)
		m.setRedis(connKey, connection)
		m.setRedis(connKeyID, connection)
	}
	return &connection, nil
}

func (m *ShipmentManager) getShopConnection(ctx context.Context, connID dot.ID, shopID dot.ID) (*connectioning.ShopConnection, error) {
	shopConnKey := getRedisShopConnectionKey(connID, shopID)
	var shopConnection connectioning.ShopConnection
	err := m.loadRedis(shopConnKey, &shopConnection)
	if err == nil {
		return &shopConnection, nil
	}
	query2 := &connectioning.GetShopConnectionByIDQuery{
		ConnectionID: connID,
		ShopID:       shopID,
	}
	if err := m.ConnectionQS.Dispatch(ctx, query2); err != nil {
		return nil, err
	}
	shopConnection = *query2.Result
	m.setRedis(shopConnKey, shopConnection)
	return &shopConnection, nil
}

func getRedisShopConnectionKey(connID dot.ID, shopID dot.ID) string {
	return fmt.Sprintf("shopConn:%v:%v%v", versionCaching, shopID.String(), connID.String())
}

func getRedisConnectionKeyByID(connID dot.ID) string {
	return fmt.Sprintf("conn:id:%v:%v", versionCaching, connID.String())
}

func getRedisConnectionKeyByCode(code string) string {
	return fmt.Sprintf("conn:code:%v:%v", versionCaching, code)
}

func (m *ShipmentManager) loadRedis(key string, v interface{}) error {
	if m.redisStore == nil {
		return cm.Errorf(cm.Internal, nil, "Redis service nil")
	}
	value, err := m.redisStore.GetString(key)
	if err != nil {
		return err
	}

	data, err := m.cipherx.Decrypt([]byte(value))
	if err != nil {
		ll.Error("Fail to decrypt from redis", l.Error(err))
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		ll.Error("Fail to unmarshal from redis", l.Error(err))
		return err
	}
	return nil
}

func (m *ShipmentManager) setRedis(key string, data interface{}) {
	if m.redisStore == nil {
		return
	}
	xData, err := json.Marshal(data)
	if err != nil {
		return
	}
	dataEncrypt, err := m.cipherx.Encrypt(xData)
	if err != nil {
		return
	}
	value := string(dataEncrypt)
	if err := m.redisStore.SetStringWithTTL(key, value, DefaultTTl); err != nil {
		ll.Error("Can not store to redis", l.Error(err))
	}
	return
}

// validateConnection
//
// Check if this connection is allowed in whitelabel partner
func (m *ShipmentManager) validateConnection(ctx context.Context, conn *connectioning.Connection) bool {
	wlPartner := wl.X(ctx)
	if !wlPartner.IsWhiteLabel() {
		return true
	}
	if conn.ConnectionMethod != connection_type.ConnectionMethodTopship {
		return true
	}
	topshipProvidersAllowed := wlPartner.Shipment.Topship
	for _, provider := range topshipProvidersAllowed {
		if provider == conn.ConnectionProvider {
			return true
		}
	}
	return false
}
