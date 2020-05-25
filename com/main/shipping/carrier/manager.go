package carrier

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/filter_type"
	"o.o/api/top/types/etc/location_type"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	addressconvert "o.o/backend/com/main/address/convert"
	addressmodel "o.o/backend/com/main/address/model"
	locationutil "o.o/backend/com/main/location/util"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/syncgroup"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cipherx"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/logic/etop_shipping_price"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/model"
	directclient "o.o/backend/pkg/integration/shipping/direct/client"
	directdriver "o.o/backend/pkg/integration/shipping/direct/driver"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	ghndriver "o.o/backend/pkg/integration/shipping/ghn/driver"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	ghtkdriver "o.o/backend/pkg/integration/shipping/ghtk/driver"
	vtpostdriver "o.o/backend/pkg/integration/shipping/vtpost/driver"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll                 = l.New()
	GHNWebhookEndpoint string
)

const (
	MinShopBalance        = -200000
	DefaultTTl            = 2 * 60 * 60
	SecretKey             = "connectionsecretkey"
	VersionCaching        = "0.1.8"
	PrefixMakeupPriceCode = "###"
)

type ShipmentManager struct {
	LocationQS             location.QueryBus
	ConnectionQS           connectioning.QueryBus
	connectionAggr         connectioning.CommandBus
	Env                    string
	redisStore             redis.Store
	cipherx                *cipherx.Cipherx
	shipmentServiceQS      shipmentservice.QueryBus
	shipmentPriceQS        shipmentprice.QueryBus
	shippingQS             shipping.QueryBus
	FlagApplyShipmentPrice bool
	eventBus               capi.EventBus
}

type Config struct {
	Endpoints []ConfigEndpoint
}

type ConfigEndpoint struct {
	Provider connection_type.ConnectionProvider
	Endpoint string
}

type FlagApplyShipmentPrice bool

func NewShipmentManager(
	eventBus capi.EventBus,
	locationQS location.QueryBus,
	connectionQS connectioning.QueryBus,
	connectionAggr connectioning.CommandBus,
	redisS redis.Store,
	shipmentServiceQS shipmentservice.QueryBus,
	shipmentPriceQS shipmentprice.QueryBus,
	flagApplyShipmentPrice FlagApplyShipmentPrice,
	cfg Config,
) (*ShipmentManager, error) {
	_cipherx, _ := cipherx.NewCipherx(SecretKey)
	sm := &ShipmentManager{
		eventBus:               eventBus,
		LocationQS:             locationQS,
		ConnectionQS:           connectionQS,
		connectionAggr:         connectionAggr,
		Env:                    cmenv.PartnerEnv(),
		redisStore:             redisS,
		cipherx:                _cipherx,
		shipmentServiceQS:      shipmentServiceQS,
		shipmentPriceQS:        shipmentPriceQS,
		FlagApplyShipmentPrice: bool(flagApplyShipmentPrice),
	}
	for _, endpoint := range cfg.Endpoints {
		err := sm.setWebhookEndpoint(endpoint.Provider, endpoint.Endpoint)
		if err != nil {
			return nil, err
		}
	}
	return sm, nil
}

func (m *ShipmentManager) setWebhookEndpoint(connectionProvider connection_type.ConnectionProvider, endpoint string) error {
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
	if connection.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
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
		cfg := directclient.PartnerAccountCfg{
			Token:      shopConnection.Token,
			Connection: connection,
		}
		if etopAffiliateAccount != nil {
			cfg.AffiliateID = etopAffiliateAccount.UserID
		}
		return directdriver.New(m.LocationQS, cfg)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}

func (m *ShipmentManager) CreateFulfillments(ctx context.Context, order *ordering.Order, ffms []*shipmodel.Fulfillment) error {
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

	// raise event to check balance
	event := &shipping.SingleFulfillmentCreatingEvent{
		EventMeta:   meta.NewEvent(),
		ShopID:      ffm.ShopID,
		FromAddress: addressconvert.Convert_addressmodel_Address_orderingtypes_Address(ffm.AddressFrom, nil),
		ShippingFee: ffm.ShippingServiceFee,
	}
	if err := m.eventBus.Publish(ctx, event); err != nil {
		return err
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

	fromDistrict, fromProvince, err := m.VerifyDistrictCode(ffm.AddressFrom)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "FromDistrictCode: %v", err)
	}
	toDistrict, toProvince, err := m.VerifyDistrictCode(ffm.AddressTo)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "ToDistrictCode: %v", err)
	}

	weight, err := ValidateFfmWeight(ffm.GrossWeight, ffm.Length, ffm.Width, ffm.Height, ffm.ChargeableWeight)
	if err != nil {
		return err
	}

	args := &GetShippingServicesArgs{
		AccountID:        ffm.ShopID,
		FromDistrictCode: fromDistrict.Code,
		FromProvinceCode: fromProvince.Code,
		FromWardCode:     ffm.AddressFrom.WardCode,
		ToDistrictCode:   toDistrict.Code,
		ToProvinceCode:   toProvince.Code,
		ToWardCode:       ffm.AddressTo.WardCode,
		ChargeableWeight: cm.CoalesceInt(weight, 100),
		Length:           ffm.Length,
		Width:            ffm.Width,
		Height:           ffm.Height,
		IncludeInsurance: ffm.IncludeInsurance,
		BasketValue:      ffm.BasketValue,
		CODAmount:        ffm.TotalCODAmount,
	}
	allServices, err := m.GetShipmentServicesAndMakeupPrice(ctx, args, ffm.ConnectionID)
	if err != nil {
		return err
	}

	isMakeupPrice := false
	makeupPrice := 0
	var providerService *model.AvailableShippingService
	if m.FlagApplyShipmentPrice {
		// áp dụng bảng giá TopShip/setting giá
		providerService, err = CheckShippingService(ffm, allServices)
		if err != nil {
			return err
		}

		// Check service with markup fee
		providerServiceID := providerService.ProviderServiceID
		if strings.HasPrefix(providerServiceID, PrefixMakeupPriceCode) {
			isMakeupPrice = true
			makeupPrice = providerService.ShippingFeeMain
			providerServiceID = providerServiceID[len(PrefixMakeupPriceCode):]
		}
		providerService.ProviderServiceID = providerServiceID
	} else {
		// backward-compatible
		// Kiểm tra gói ETOP
		var etopService *model.AvailableShippingService
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
			isMakeupPrice = true
			makeupPrice = etopService.ShippingFeeMain
		} else {
			// Provider service
			// => Check price
			// => Get this service
			providerService, err = CheckShippingService(ffm, allServices)
			if err != nil {
				return err
			}
		}
	}

	_args := args.ToShipmentServiceArgs(ffm.ConnectionID, ffm.ShopID)
	ffmToUpdate, err := driver.CreateFulfillment(ctx, ffm, _args, providerService)
	if err != nil {
		return err
	}
	// update shipping service name
	ffmToUpdate.ShippingServiceName = providerService.Name

	if isMakeupPrice {
		err := ffmToUpdate.ApplyEtopPrice(makeupPrice)
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

func (m *ShipmentManager) GetShippingServices(ctx context.Context, args *GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	accountID := args.AccountID
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
				return cm.Errorf(cm.FailedPrecondition, nil, "Connection does not valid (check status or token)")
			}
			services, err := m.GetShipmentServicesAndMakeupPrice(ctx, args, connID)
			if err != nil {
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

	// Chỉ có method=direct mới được login
	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Do not support this feature for this connection")
	}

	var userID, token string
	if conn.EtopAffiliateAccount != nil {
		userID = conn.EtopAffiliateAccount.UserID
		token = conn.EtopAffiliateAccount.Token
	}

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
	case connection_type.ConnectionProviderPartner:
		cfg := directclient.PartnerAccountCfg{
			Connection: conn,
			Token:      token,
		}
		return directdriver.New(m.LocationQS, cfg)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hỗ trợ affiliate account (connType = %v, connName = %v)", conn.ConnectionProvider, conn.Name)
	}
}

func (m *ShipmentManager) UpdateFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (updateFfm *shipmodel.Fulfillment, err error) {
	connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
	driver, err := m.getShipmentDriver(ctx, connectionID, ffm.ShopID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "invalid connection (ffm_id = %v)", ffm.ID)
	}

	updateFfm, err = driver.UpdateFulfillment(ctx, ffm)
	if err != nil {
		return nil, err
	}
	return
}

func (m *ShipmentManager) GetConnectionByID(ctx context.Context, connID dot.ID) (*connectioning.Connection, error) {
	connKey := GetRedisConnectionKeyByID(connID)
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
		connKeyID := GetRedisConnectionKeyByID(connection.ID)
		m.setRedis(connKey, connection)
		m.setRedis(connKeyID, connection)
	}
	return &connection, nil
}

func (m *ShipmentManager) getShopConnection(ctx context.Context, connID dot.ID, shopID dot.ID) (*connectioning.ShopConnection, error) {
	shopConnKey := GetRedisShopConnectionKey(connID, shopID)
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

func GetRedisShopConnectionKey(connID dot.ID, shopID dot.ID) string {
	return fmt.Sprintf("shopConn:%v:%v%v", VersionCaching, shopID.String(), connID.String())
}

func GetRedisConnectionKeyByID(connID dot.ID) string {
	return fmt.Sprintf("conn:id:%v:%v", VersionCaching, connID.String())
}

func getRedisConnectionKeyByCode(code string) string {
	return fmt.Sprintf("conn:code:%v:%v", VersionCaching, code)
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
	if conn.Status != status3.P {
		return false
	}
	wlPartner := wl.X(ctx)
	if !wlPartner.IsWhiteLabel() {
		return true
	}
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
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

func (m *ShipmentManager) GetShipmentServicesAndMakeupPrice(ctx context.Context, args *GetShippingServicesArgs, connID dot.ID) ([]*model.AvailableShippingService, error) {
	accountID := args.AccountID
	conn, err := m.GetConnectionByID(ctx, connID)
	if err != nil {
		return nil, err
	}
	if !m.validateConnection(ctx, conn) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Connection does not valid")
	}

	var services []*model.AvailableShippingService
	driver, err := m.getShipmentDriver(ctx, connID, accountID)
	if err != nil {
		// ll.Error("Driver shipment không hợp lệ", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
		return nil, err
	}

	_args := args.ToShipmentServiceArgs(connID, accountID)
	if !m.FlagApplyShipmentPrice {
		if conn.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
			_args.IncludeTopshipServices = true
		}
	}

	services, err = driver.GetShippingServices(ctx, _args)
	if err != nil {
		// ll.Error("Get service error", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
		return nil, err
	}

	var res []*model.AvailableShippingService
	// assign connection info to services
	for _, s := range services {
		s.ConnectionInfo = &model.ConnectionInfo{
			ID:       connID,
			Name:     conn.Name,
			ImageURL: conn.ImageURL,
		}
		if !m.FlagApplyShipmentPrice {
			// không áp dụng bảng giá
			res = append(res, s)
			continue
		}
		serviceID, err := driver.ParseServiceID(s.ProviderServiceID)
		if err != nil {
			ll.Error("Parse service ID failed", l.String("serviceID", serviceID), l.ID("connectionID", connID), l.Error(err))
		}

		if err := m.mapWithShipmentService(ctx, args, serviceID, connID, s); err != nil {
			continue
		}

		// Makeup price & change provider_service_id
		if err = m.makeupPriceByShipmentPrice(ctx, s, args); err != nil {
			// ll.Error("MakeupPriceByShipmentPrice failed", l.String("serviceID", serviceID), l.ID("connectionID", connID), l.Error(err))
			// continue
		}
		res = append(res, s)
	}
	return res, nil
}

func (m *ShipmentManager) mapWithShipmentService(ctx context.Context, args *GetShippingServicesArgs, serviceID string, connID dot.ID, service *model.AvailableShippingService) error {
	sService, err := m.getShipmentService(ctx, args, serviceID, connID, false)
	if err != nil {
		return err
	}

	service.ShipmentServiceInfo = &model.ShipmentServiceInfo{
		ID:          sService.ID,
		Code:        sService.EdCode,
		Name:        sService.Name,
		IsAvailable: true,
	}
	service.Name = sService.Name

	// Kiểm tra các khu vực blacklist
	// Nếu nằm trong khu vực blacklist thì vẫn trả về gói dịch vụ, kèm theo thông tin lỗi để client hiển thị
	// Khi tạo đơn với gói này cần kiểm tra `IsAvailable` hay không và trả về lỗi nếu có
	if err := checkShipmentServiceBlacklists(args, sService.BlacklistLocations); err != nil {
		service.ShipmentServiceInfo.IsAvailable = false
		service.ShipmentServiceInfo.ErrorMessage = err.Error()
	}

	return nil
}

func (m *ShipmentManager) getShipmentService(ctx context.Context, args *GetShippingServicesArgs, serviceID string, connID dot.ID, checkBlackListLocation bool) (*shipmentservice.ShipmentService, error) {
	query := &shipmentservice.GetShipmentServiceByServiceIDQuery{
		ServiceID: serviceID,
		ConnID:    connID,
	}
	if err := m.shipmentServiceQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	sService := query.Result
	if sService.Status != status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Gói dịch vụ đã tắt (shipment_service_id = %v).", sService.ID)
	}

	// Kiểm tra các khu vực khả dụng của gói
	if err := m.checkShipmentServiceAvailableLocations(ctx, args, sService.AvailableLocations); err != nil {
		ll.Error("checkShipmentServiceAvailableLocation failed", l.String("serviceID", serviceID), l.ID("shipment_service_id", sService.ID), l.ID("connectionID", connID), l.Error(err))
		return nil, err
	}

	// Kiểm tra khối lượng khả dụng
	if sService.OtherCondition != nil {
		weight := args.ChargeableWeight
		minWeight := sService.OtherCondition.MinWeight
		maxWeight := sService.OtherCondition.MaxWeight
		if weight < minWeight || (weight > maxWeight && maxWeight != -1) {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Khối lượng nằm ngoài mức khả dụng của gói")
		}
	}

	if checkBlackListLocation {
		if err := checkShipmentServiceBlacklists(args, sService.BlacklistLocations); err != nil {
			return nil, err
		}
	}
	return sService, nil
}

func (m *ShipmentManager) checkShipmentServiceAvailableLocations(ctx context.Context, args *GetShippingServicesArgs, als []*shipmentservice.AvailableLocation) error {
	for _, al := range als {
		if err := m.checkShipmentServiceAvailableLocation(ctx, args, al); err != nil {
			return err
		}
	}
	return nil
}

func (m *ShipmentManager) checkShipmentServiceAvailableLocation(ctx context.Context, args *GetShippingServicesArgs, al *shipmentservice.AvailableLocation) error {
	if al == nil {
		return nil
	}
	var isInclude bool
	var provinceCode string
	switch al.FilterType {
	case filter_type.Include:
		isInclude = true
	case filter_type.Exclude:
		isInclude = false
	default:
		return cm.Errorf(cm.Internal, nil, "filter_type không hợp lệ").WithMetap("availableLocation", al)
	}
	switch al.ShippingLocationType {
	case location_type.ShippingLocationPick:
		provinceCode = args.FromProvinceCode

	case location_type.ShippingLocationDeliver:
		provinceCode = args.ToProvinceCode
	default:
		return cm.Errorf(cm.Internal, nil, "shipping_location_type không hợp lệ").WithMetap("availableLocation", al)
	}

	shippingLocationLabel := al.ShippingLocationType.GetLabelRefName()
	if len(al.RegionTypes) > 0 {
		regionType := locationutil.GetRegion(provinceCode, "")
		isContain := location_type.RegionTypeContains(al.RegionTypes, regionType)
		if isInclude && !isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v nằm ngoài miền quy định", shippingLocationLabel)
		}
		if !isInclude && isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v nằm trong miền bị loại trừ", shippingLocationLabel)
		}
	}

	if len(al.CustomRegionIDs) > 0 {
		var customRegionID dot.ID
		query := &location.GetCustomRegionByCodeQuery{
			ProvinceCode: provinceCode,
		}
		if err := m.LocationQS.Dispatch(ctx, query); err != nil {
			return err
		}
		customRegionID = query.Result.ID
		isContain := cm.IDsContain(al.CustomRegionIDs, customRegionID)
		if isInclude && !isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v nằm ngoài vùng quy định", shippingLocationLabel)
		}
		if !isInclude && isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v nằm trong vùng bị loại trừ", shippingLocationLabel)
		}
	}

	if len(al.ProvinceCodes) > 0 {
		isContain := cm.StringsContain(al.ProvinceCodes, provinceCode)
		if isInclude && !isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v nằm ngoài tỉnh quỷ định", shippingLocationLabel)
		}
		if !isInclude && isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v nằm trong vùng bị loại trừ", shippingLocationLabel)
		}
	}
	return nil
}

func checkShipmentServiceBlacklists(args *GetShippingServicesArgs, bls []*shipmentservice.BlacklistLocation) error {
	for _, bl := range bls {
		if err := checkShipmentServiceBlacklist(args, bl); err != nil {
			return err
		}
	}
	return nil
}

func checkShipmentServiceBlacklist(args *GetShippingServicesArgs, bl *shipmentservice.BlacklistLocation) error {
	if bl == nil {
		return nil
	}
	var provinceCode, districtCode, wardCode string
	switch bl.ShippingLocationType {
	case location_type.ShippingLocationPick:
		provinceCode = args.FromProvinceCode
		districtCode = args.FromDistrictCode
		wardCode = args.FromWardCode

	case location_type.ShippingLocationDeliver:
		provinceCode = args.ToProvinceCode
		districtCode = args.ToDistrictCode
		wardCode = args.FromWardCode
	default:
		return cm.Errorf(cm.Internal, nil, "shipping_location_type không hợp lệ").WithMetap("blacklist", bl)
	}

	shippingLocationLabel := bl.ShippingLocationType.GetLabelRefName()
	if cm.StringsContain(bl.ProvinceCodes, provinceCode) ||
		cm.StringsContain(bl.DistrictCodes, districtCode) ||
		cm.StringsContain(bl.WardCodes, wardCode) {
		return cm.Errorf(cm.FailedPrecondition, nil, "%v không khả dụng. %v.", shippingLocationLabel, bl.Reason)
	}
	return nil
}

func (m *ShipmentManager) makeupPriceByShipmentPrice(ctx context.Context, service *model.AvailableShippingService, args *GetShippingServicesArgs) error {
	// Cập nhật giá Makeup vào phí chính của gói
	// TH gói dịch vụ ko tính được phí chính
	// Tạm thời để an toàn không áp dụng bảng giá vào đây
	if service.ShippingFeeMain == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Gói dịch vụ chưa tính được phí chính. Không đủ điều kiện để áp dụng bảng giá vào đây")
	}
	if service.ShipmentServiceInfo == nil || service.ShipmentServiceInfo.ID == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Thiếu shipment service.")
	}
	originMainFee := service.ShippingFeeMain
	additionFee := service.ServiceFee - originMainFee
	originFee := service.ServiceFee

	query := &shipmentprice.CalculatePriceQuery{
		AccountID:           args.AccountID,
		ShipmentPriceListID: args.ShipmentPriceListID,
		FromDistrictCode:    args.FromDistrictCode,
		ToDistrictCode:      args.ToDistrictCode,
		ShipmentServiceID:   service.ShipmentServiceInfo.ID,
		Weight:              args.ChargeableWeight,
	}
	err := m.shipmentPriceQS.Dispatch(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NoError:
	// continue
	case cm.NotFound:
		return nil
	default:
		return cm.Errorf(cm.Internal, err, "")
	}
	calcPriceResp := query.Result
	service.ProviderServiceID = PrefixMakeupPriceCode + service.ProviderServiceID
	service.ServiceFee = calcPriceResp.Price + additionFee
	service.ShippingFeeMain = calcPriceResp.Price
	service.ShipmentPriceInfo = &model.ShipmentPriceInfo{
		ID:            calcPriceResp.ShipmentPriceID,
		OriginFee:     originFee,
		OriginMainFee: originMainFee,
		MakeupMainFee: calcPriceResp.Price,
	}
	return nil
}

func (m *ShipmentManager) CalcMakeupShipmentPrice(ctx context.Context, ffm *shipmodel.Fulfillment, weight int) (int, error) {
	connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
	driver, err := m.getShipmentDriver(ctx, connectionID, ffm.ShopID)
	if err != nil {
		return 0, cm.Errorf(cm.InvalidArgument, err, "invalid connection (ffm_id = %v)", ffm.ID)
	}
	serviceID, err := driver.ParseServiceID(ffm.ProviderServiceID)
	if err != nil {
		return 0, err
	}

	args := &GetShippingServicesArgs{
		FromDistrictCode: ffm.AddressFrom.DistrictCode,
		FromProvinceCode: ffm.AddressFrom.ProvinceCode,
		FromWardCode:     ffm.AddressFrom.WardCode,
		ToDistrictCode:   ffm.AddressTo.DistrictCode,
		ToProvinceCode:   ffm.AddressTo.ProvinceCode,
		ToWardCode:       ffm.AddressTo.WardCode,
		ChargeableWeight: weight,
		BasketValue:      ffm.BasketValue,
		CODAmount:        ffm.TotalCODAmount,
	}

	shipmentService, err := m.getShipmentService(ctx, args, serviceID, connectionID, true)
	if err != nil {
		return 0, err
	}

	query := &shipmentprice.CalculatePriceQuery{
		ShipmentPriceListID: args.ShipmentPriceListID,
		FromDistrictCode:    args.FromDistrictCode,
		ToDistrictCode:      args.ToDistrictCode,
		ShipmentServiceID:   shipmentService.ID,
		Weight:              args.ChargeableWeight,
	}
	if err := m.shipmentPriceQS.Dispatch(ctx, query); err != nil {
		return 0, err
	}
	return query.Result.Price, nil
}
