package carrier

import (
	"context"
	"strings"
	"sync"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/main/shippingcode"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/filter_type"
	"o.o/api/top/types/etc/location_type"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	addressconvert "o.o/backend/com/main/address/convert"
	addressmodel "o.o/backend/com/main/address/model"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	locationutil "o.o/backend/com/main/location/util"
	shipmentpriceconvert "o.o/backend/com/main/shipmentpricing/shipmentprice/convert"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/syncgroup"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

const (
	PrefixMakeupPriceCode = "###"

	ThirtyMinutes = 30 * time.Minute
)

type ShipmentManager struct {
	locationQS           location.QueryBus
	identityQS           identity.QueryBus
	connectionQS         connectioning.QueryBus
	connectionAggr       connectioning.CommandBus
	env                  string
	shippingcodeQS       shippingcode.QueryBus
	shipmentServiceQS    shipmentservice.QueryBus
	shipmentPriceQS      shipmentprice.QueryBus
	shippingQS           shipping.QueryBus
	priceListPromotionQS pricelistpromotion.QueryBus
	ConnectionManager    *connectionmanager.ConnectionManager

	carrierDriver carriertypes.Driver
	eventBus      capi.EventBus
	OrderStore    sqlstore.OrderStoreInterface
}

func NewShipmentManager(
	eventBus capi.EventBus,
	locationQS location.QueryBus,
	identityQS identity.QueryBus,
	connectionQS connectioning.QueryBus,
	connectionAggr connectioning.CommandBus,
	shippingcodeQS shippingcode.QueryBus,
	shipmentServiceQS shipmentservice.QueryBus,
	shipmentPriceQS shipmentprice.QueryBus,
	priceListPromotionQS pricelistpromotion.QueryBus,
	carrierDriver carriertypes.Driver,
	connectionManager *connectionmanager.ConnectionManager,
	OrderStore sqlstore.OrderStoreInterface,
) (*ShipmentManager, error) {
	sm := &ShipmentManager{
		eventBus:             eventBus,
		locationQS:           locationQS,
		identityQS:           identityQS,
		connectionQS:         connectionQS,
		connectionAggr:       connectionAggr,
		env:                  cmenv.PartnerEnv(),
		shippingcodeQS:       shippingcodeQS,
		shipmentServiceQS:    shipmentServiceQS,
		shipmentPriceQS:      shipmentPriceQS,
		priceListPromotionQS: priceListPromotionQS,
		ConnectionManager:    connectionManager,
		carrierDriver:        carrierDriver,
		OrderStore:           OrderStore,
	}
	return sm, nil
}

func (m *ShipmentManager) GetShipmentDriver(ctx context.Context, connectionID dot.ID, shopID dot.ID) (carriertypes.ShipmentCarrier, error) {
	connection, err := m.ConnectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	getShopConnectionQuery := connectionmanager.GetShopConnectionArgs{
		ConnectionID: connectionID,
		ShopID:       shopID,
	}
	if connection.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		// ignore shopID
		getShopConnectionQuery.ShopID = 0
		getShopConnectionQuery.IsGlobal = true
	}

	shopConnection, err := m.ConnectionManager.GetShopConnection(ctx, getShopConnectionQuery)
	if err != nil {
		return nil, err
	}

	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid (check status or token)")
	}

	// check if token is expired then generate new token
	if err := m.generateToken(ctx, connection, shopConnection); err != nil {
		return nil, err
	}

	shipmentDriver, err := m.carrierDriver.GetShipmentDriver(m.env, m.locationQS, m.identityQS, connection, shopConnection, m.shippingcodeQS)
	if err != nil {
		return nil, err
	}

	return shipmentDriver, nil
}

func (m *ShipmentManager) generateToken(ctx context.Context, connection *connectioning.Connection, shopConnection *connectioning.ShopConnection) error {
	expiresAt := shopConnection.TokenExpiresAt
	if expiresAt.IsZero() {
		return nil
	}
	now := time.Now()
	// 30p tr?????c khi h???t h???n
	expiresAt.Add(-ThirtyMinutes)
	if expiresAt.After(now) {
		return nil
	}

	// get driver
	_driver, err := m.GetDriverByEtopAffiliateAccount(ctx, connection.ID)
	if err != nil {
		return err
	}

	// re-generate token
	generateTokenResp, err := _driver.GenerateToken(ctx)
	if err != nil {
		return err
	}

	// update shopConnection
	updateShopConnectionCmd := connectioning.CreateOrUpdateShopConnectionCommand{
		ShopID:         shopConnection.ShopID,
		ConnectionID:   shopConnection.ConnectionID,
		Token:          generateTokenResp.AccessToken,
		TokenExpiresAt: generateTokenResp.ExpiresAt,
		ExternalData:   shopConnection.ExternalData,
	}
	if err := m.connectionAggr.Dispatch(ctx, &updateShopConnectionCmd); err != nil {
		return err
	}
	*shopConnection = *updateShopConnectionCmd.Result
	return nil
}

func (m *ShipmentManager) CreateFulfillments(ctx context.Context, ffms []*shipmodel.Fulfillment) error {
	var err error
	g := syncgroup.New(len(ffms))
	for _, ffm := range ffms {
		ffm := ffm
		g.Go(func() error { return m.createSingleFulfillment(ctx, ffm) })
	}
	errs := g.Wait()
	if errs.IsAll() {
		err = errs[0]
	}
	return err
}

func (m *ShipmentManager) createSingleFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (_err error) {
	connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
	driver, err := m.GetShipmentDriver(ctx, connectionID, ffm.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "invalid connection")
	}

	// raise event to check balance
	event := &shipping.SingleFulfillmentCreatingEvent{
		EventMeta:    meta.NewEvent(),
		ShopID:       ffm.ShopID,
		FromAddress:  addressconvert.Convert_addressmodel_Address_orderingtypes_Address(ffm.AddressFrom, nil),
		ShippingFee:  ffm.ShippingServiceFee,
		ConnectionID: connectionID,
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
		_ = m.OrderStore.UpdateFulfillment(ctx, cmd)
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
		IncludeInsurance: ffm.IncludeInsurance.Apply(false),
		InsuranceValue:   ffm.InsuranceValue,
		BasketValue:      ffm.BasketValue,
		CODAmount:        ffm.TotalCODAmount,
		Coupon:           ffm.Coupon,
	}
	allServices, err := m.GetShipmentServicesAndMakeupPrice(ctx, args, connectionID)
	if err != nil {
		return err
	}

	isMakeupPrice := false
	makeupPriceMain := 0
	var providerService *shippingsharemodel.AvailableShippingService
	providerService, err = CheckShippingService(ffm, allServices)
	if err != nil {
		return err
	}

	// Check service with makeup fee
	providerServiceID := providerService.ProviderServiceID
	if strings.HasPrefix(providerServiceID, PrefixMakeupPriceCode) {
		isMakeupPrice = true
		makeupPriceMain = providerService.ShippingFeeMain
		providerServiceID = providerServiceID[len(PrefixMakeupPriceCode):]
	}
	providerService.ProviderServiceID = providerServiceID

	_args := args.ToShipmentServiceArgs(connectionID, ffm.ShopID)
	ffm.ShippingNote = carriertypes.GetShippingCarrierNote(ffm)

	ffmToUpdate, err := driver.CreateFulfillment(ctx, ffm, _args, providerService)
	if err != nil {
		return err
	}

	// update shipping service name
	ffmToUpdate.ShippingServiceName = providerService.Name

	if providerService.ShipmentPriceInfo != nil {
		ffmToUpdate.ShipmentPriceInfo = providerService.ShipmentPriceInfo
	}
	if isMakeupPrice {
		ffmToUpdate.ApplyEtopPrice(makeupPriceMain)
		ffmToUpdate.ShippingFeeShopLines = providerService.ShippingFeeLines
	} else {
		ffmToUpdate.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(ffmToUpdate.ProviderShippingFeeLines, false, dot.Int(0))
	}
	ffmToUpdate.ExternalAffiliateID = driver.GetAffiliateID()
	ffmToUpdate.ChargeableWeight = weight
	updateCmd := &shipmodelx.UpdateFulfillmentCommand{
		Fulfillment: ffmToUpdate,
	}
	if err := m.OrderStore.UpdateFulfillment(ctx, updateCmd); err != nil {
		return cm.Trace(err)
	}
	return nil
}

func (m *ShipmentManager) VerifyDistrictCode(addr *addressmodel.Address) (*location.District, *location.Province, error) {
	if addr == nil {
		return nil, nil, cm.Errorf(cm.Internal, nil, "?????a ch??? kh??ng t???n t???i")
	}
	if addr.DistrictCode == "" {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil,
			`?????a ch??? %v, %v kh??ng th??? ???????c x??c ?????nh b???i h??? th???ng.`,
			addr.District, addr.Province,
		)
	}

	query := &location.GetLocationQuery{DistrictCode: addr.DistrictCode}
	if err := m.locationQS.Dispatch(context.TODO(), query); err != nil {
		return nil, nil, err
	}
	district := query.Result.District
	return district, query.Result.Province, nil
}

func (m *ShipmentManager) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	driver, err := m.GetShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "invalid connection")
	}
	return driver.CancelFulfillment(ctx, ffm)
}

func (m *ShipmentManager) GetShippingServices(ctx context.Context, args *GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	accountID := args.AccountID
	shopConnections, err := m.GetAllShopConnections(ctx, accountID, args.ConnectionIDs)
	if err != nil {
		return nil, err
	}
	var res []*shippingsharemodel.AvailableShippingService
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
				ll.Error("GetShipmentServicesAndMakeupPrice :: ", l.ID("connection_id", connID), l.Error(err))
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
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Kh??ng c?? g??i giao h??ng ph?? h???p")
	}
	res = CompactServices(res)
	return res, nil
}

func (m *ShipmentManager) GetAllShopConnections(ctx context.Context, shopID dot.ID, connectionIDs []dot.ID) ([]*connectioning.ShopConnection, error) {
	if len(connectionIDs) == 0 {
		// Ch??? l???y nh???ng connection type shipping, subtype shipment
		queryConn := &connectioning.ListConnectionsQuery{
			IDs:               connectionIDs,
			Status:            status3.P.Wrap(),
			ConnectionType:    connection_type.Shipping,
			ConnectionSubtype: connection_type.ConnectionSubtypeShipment,
		}
		if err := m.connectionQS.Dispatch(ctx, queryConn); err != nil {
			return nil, err
		}
		conns := queryConn.Result
		if len(conns) == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Kh??ng c?? nh?? v???n chuy???n h???p l???")
		}
		for _, conn := range conns {
			connectionIDs = append(connectionIDs, conn.ID)
		}
	}

	// Get all shop_connection & global shop_connection
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

func (m *ShipmentManager) SignIn(ctx context.Context, args *ConnectionSignInArgs) (*carriertypes.AccountResponse, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}

	driver, err := m.GetDriverByEtopAffiliateAccount(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}
	cmd := &carriertypes.SignInArgs{
		Identifier: args.Identifier,
		Password:   args.Password,
		OTP:        args.OTP,
	}
	return driver.SignIn(ctx, cmd)
}

func (m *ShipmentManager) SignUp(ctx context.Context, args *ConnectionSignUpArgs) (newAccount *carriertypes.AccountResponse, _ error) {
	driver, err := m.GetDriverByEtopAffiliateAccount(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}
	cmd := &carriertypes.SignUpArgs{
		Name:     args.Name,
		Email:    args.Identifier,
		Password: args.Password,
		Phone:    args.Phone,
		Province: args.Province,
		District: args.District,
		Address:  args.Address,
	}
	return driver.SignUp(ctx, cmd)
}

func (m *ShipmentManager) GetDriverByEtopAffiliateAccount(ctx context.Context, connectionID dot.ID) (carriertypes.ShipmentCarrier, error) {
	conn, err := m.ConnectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}

	return m.carrierDriver.GetAffiliateShipmentDriver(m.env, m.locationQS, m.identityQS, conn, m.shippingcodeQS)
}

func (m *ShipmentManager) RefreshFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (updateFfm *shipmodel.Fulfillment, err error) {
	connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
	driver, err := m.GetShipmentDriver(ctx, connectionID, ffm.ShopID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "invalid connection (ffm_id = %v)", ffm.ID)
	}

	updateFfm, err = driver.RefreshFulfillment(ctx, ffm)
	if err != nil {
		return nil, err
	}
	return
}

func (m *ShipmentManager) UpdateFulfillmentInfo(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	driver, err := m.GetShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "invalid connection")
	}

	if _, _, err := m.VerifyDistrictCode(ffm.AddressFrom); err != nil {
		return cm.Errorf(cm.Internal, err, "FromDistrictCode: %v", err)
	}
	if _, _, err := m.VerifyDistrictCode(ffm.AddressTo); err != nil {
		return cm.Errorf(cm.Internal, err, "ToDistrictCode: %v", err)
	}

	return driver.UpdateFulfillmentInfo(ctx, ffm)
}

func (m *ShipmentManager) UpdateFulfillmentCOD(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	driver, err := m.GetShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "invalid connection")
	}
	return driver.UpdateFulfillmentCOD(ctx, ffm)
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
	if len(topshipProvidersAllowed) == 0 {
		return true
	}
	for _, provider := range topshipProvidersAllowed {
		if provider == conn.ConnectionProvider {
			return true
		}
	}
	return false
}

func (m *ShipmentManager) GetShipmentServicesAndMakeupPrice(ctx context.Context, args *GetShippingServicesArgs, connID dot.ID) ([]*shippingsharemodel.AvailableShippingService, error) {
	accountID := args.AccountID
	conn, err := m.ConnectionManager.GetConnectionByID(ctx, connID)
	if err != nil {
		return nil, err
	}
	if !m.validateConnection(ctx, conn) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Connection does not valid")
	}

	var services []*shippingsharemodel.AvailableShippingService
	driver, err := m.GetShipmentDriver(ctx, connID, accountID)
	if err != nil {
		// ll.Error("Driver shipment kh??ng h???p l???", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
		return nil, err
	}

	_args := args.ToShipmentServiceArgs(connID, accountID)
	services, err = driver.GetShippingServices(ctx, _args)
	if err != nil {
		// ll.Error("Get service error", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
		return nil, err
	}

	var res []*shippingsharemodel.AvailableShippingService
	// assign connection info to services
	for _, s := range services {
		s.ConnectionInfo = &shippingsharemodel.ConnectionInfo{
			ID:       connID,
			Name:     conn.Name,
			ImageURL: conn.ImageURL,
		}
		if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
			// kh??ng ??p d???ng b???ng gi??
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
		// N???u kh??ng c?? c???u h??nh gi?? (shipment_price) c???a m?? d???ch v??? (shipment_service_id) trong b???ng gi?? (m?? l???i not_found)
		// Tr??? v??? k???t qu??? c???a NVC
		if err = m.makeupPriceByShipmentPrice(ctx, s, args); err != nil && cm.ErrorCode(err) != cm.NotFound {
			ll.Error("MakeupPriceByShipmentPrice failed", l.String("serviceID", serviceID), l.ID("connectionID", connID), l.Error(err))
			continue
		}
		res = append(res, s)
	}
	return filterShipmentServicesByEdCode(res), nil
}

func (m *ShipmentManager) mapWithShipmentService(ctx context.Context, args *GetShippingServicesArgs, serviceID string, connID dot.ID, service *shippingsharemodel.AvailableShippingService) error {
	sService, err := m.getShipmentService(ctx, args, serviceID, connID, false)
	if err != nil {
		return err
	}

	service.ShipmentServiceInfo = &shippingsharemodel.ShipmentServiceInfo{
		ID:          sService.ID,
		Code:        sService.EdCode,
		Name:        sService.Name,
		IsAvailable: true,
	}
	service.Name = sService.Name

	// Ki???m tra c??c khu v???c blacklist
	// N???u n???m trong khu v???c blacklist th?? v???n tr??? v??? g??i d???ch v???, k??m theo th??ng tin l???i ????? client hi???n th???
	// Khi t???o ????n v???i g??i n??y c???n ki???m tra `IsAvailable` hay kh??ng v?? tr??? v??? l???i n???u c??
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
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "G??i d???ch v??? ???? t???t (shipment_service_id = %v).", sService.ID)
	}

	// Ki???m tra c??c khu v???c kh??? d???ng c???a g??i
	if err := m.checkShipmentServiceAvailableLocations(ctx, args, sService.AvailableLocations); err != nil {
		ll.Error("checkShipmentServiceAvailableLocation failed", l.String("serviceID", serviceID), l.ID("shipment_service_id", sService.ID), l.ID("connectionID", connID), l.Error(err))
		return nil, err
	}

	// Ki???m tra kh???i l?????ng kh??? d???ng
	if sService.OtherCondition != nil {
		weight := args.ChargeableWeight
		minWeight := sService.OtherCondition.MinWeight
		maxWeight := sService.OtherCondition.MaxWeight
		if weight < minWeight || (weight > maxWeight && maxWeight != -1) {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Kh???i l?????ng n???m ngo??i m???c kh??? d???ng c???a g??i")
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
		return cm.Errorf(cm.Internal, nil, "filter_type kh??ng h???p l???").WithMetap("availableLocation", al)
	}
	switch al.ShippingLocationType {
	case location_type.ShippingLocationPick:
		provinceCode = args.FromProvinceCode

	case location_type.ShippingLocationDeliver:
		provinceCode = args.ToProvinceCode
	default:
		return cm.Errorf(cm.Internal, nil, "shipping_location_type kh??ng h???p l???").WithMetap("availableLocation", al)
	}

	shippingLocationLabel := al.ShippingLocationType.GetLabelRefName()
	if len(al.RegionTypes) > 0 {
		regionType := locationutil.GetRegion(provinceCode, "")
		isContain := location_type.RegionTypeContains(al.RegionTypes, regionType)
		if isInclude && !isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v n???m ngo??i mi???n quy ?????nh", shippingLocationLabel)
		}
		if !isInclude && isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v n???m trong mi???n b??? lo???i tr???", shippingLocationLabel)
		}
	}

	if len(al.CustomRegionIDs) > 0 {
		query := &location.ListCustomRegionsByCodeQuery{
			ProvinceCode: provinceCode,
		}
		if err := m.locationQS.Dispatch(ctx, query); err != nil {
			return err
		}
		isContain := false
		for _, customRegion := range query.Result {
			isContain = cm.IDsContain(al.CustomRegionIDs, customRegion.ID)
			if isContain {
				break
			}
		}
		if isInclude && !isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v n???m ngo??i v??ng quy ?????nh", shippingLocationLabel)
		}
		if !isInclude && isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v n???m trong v??ng b??? lo???i tr???", shippingLocationLabel)
		}
	}

	if len(al.ProvinceCodes) > 0 {
		isContain := cm.StringsContain(al.ProvinceCodes, provinceCode)
		if isInclude && !isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v n???m ngo??i t???nh qu??? ?????nh", shippingLocationLabel)
		}
		if !isInclude && isContain {
			return cm.Errorf(cm.FailedPrecondition, nil, "%v n???m trong v??ng b??? lo???i tr???", shippingLocationLabel)
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
		wardCode = args.ToWardCode
	default:
		return cm.Errorf(cm.Internal, nil, "shipping_location_type kh??ng h???p l???").WithMetap("blacklist", bl)
	}

	shippingLocationLabel := bl.ShippingLocationType.GetLabelRefName()
	if cm.StringsContain(bl.ProvinceCodes, provinceCode) ||
		cm.StringsContain(bl.DistrictCodes, districtCode) ||
		cm.StringsContain(bl.WardCodes, wardCode) {
		return cm.Errorf(cm.FailedPrecondition, nil, "%v kh??ng kh??? d???ng. %v.", shippingLocationLabel, bl.Reason)
	}
	return nil
}

func (m *ShipmentManager) makeupPriceByShipmentPrice(ctx context.Context, service *shippingsharemodel.AvailableShippingService, args *GetShippingServicesArgs) error {
	if service.ShipmentServiceInfo == nil || service.ShipmentServiceInfo.ID == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Thi???u shipment service.")
	}
	originFee := service.ServiceFee
	addFeeTypes := []shipping_fee_type.ShippingFeeType{}
	if args.IncludeInsurance {
		addFeeTypes = append(addFeeTypes, shipping_fee_type.Insurance)
	}

	query := &shipmentprice.CalculateShippingFeesQuery{
		AccountID:           args.AccountID,
		ShipmentPriceListID: args.ShipmentPriceListID,
		FromDistrictCode:    args.FromDistrictCode,
		ToDistrictCode:      args.ToDistrictCode,
		ShipmentServiceID:   service.ShipmentServiceInfo.ID,
		ConnectionID:        service.ConnectionInfo.ID,
		Weight:              args.ChargeableWeight,
		BasketValue:         args.BasketValue,
		CODAmount:           args.CODAmount,
		AdditionalFeeTypes:  addFeeTypes,
	}

	// get pricelist promotion
	queryPromotion := &GetPriceListPromotionArgs{
		ShopID:           args.AccountID,
		FromProvinceCode: args.FromProvinceCode,
		ConnectionID:     service.ConnectionInfo.ID,
	}
	if promotionPriceListID, err := m.getPromotionPriceListID(ctx, queryPromotion); err == nil {
		query.PromotionPriceListID = promotionPriceListID
	}

	if err := m.shipmentPriceQS.Dispatch(ctx, query); err != nil {
		return err
	}

	calcShippingFeesRes := query.Result
	service.ProviderServiceID = PrefixMakeupPriceCode + service.ProviderServiceID
	service.ServiceFee = calcShippingFeesRes.TotalFee
	feeLines := shipmentpriceconvert.Convert_shipmentprice_ShippingFees_To_shippingsharemodel_ShippingFeeLines(calcShippingFeesRes.FeeLines)
	service.ShippingFeeLines = feeLines
	service.ShippingFeeMain = shippingsharemodel.GetMainFee(feeLines)
	service.ShipmentPriceInfo = &shippingsharemodel.ShipmentPriceInfo{
		ShipmentPriceID:     calcShippingFeesRes.ShipmentPriceID,
		ShipmentPriceListID: calcShippingFeesRes.ShipmentPriceListID,
		OriginFee:           originFee,
		MakeupFee:           calcShippingFeesRes.TotalFee,
	}
	return nil
}

func (m *ShipmentManager) getPromotionPriceListID(ctx context.Context, args *GetPriceListPromotionArgs) (priceListID dot.ID, _ error) {
	query := &pricelistpromotion.GetValidPriceListPromotionQuery{
		ShopID:           args.ShopID,
		FromProvinceCode: args.FromProvinceCode,
		ConnectionID:     args.ConnectionID,
	}
	if err := m.priceListPromotionQS.Dispatch(ctx, query); err != nil {
		return 0, err
	}

	return query.Result.PriceListID, nil
}

type CalcMakeupShippingFeesByFfmArgs struct {
	Fulfillment        *shipping.Fulfillment
	Weight             int
	State              shippingstate.State
	AdditionalFeeTypes []shipping_fee_type.ShippingFeeType
}

type CalcMakeupShippingFeesByFfmResponse struct {
	ShipmentPriceID     dot.ID
	ShipmentPriceListID dot.ID
	ShippingFeeLines    []*shippingtypes.ShippingFeeLine
}

func (m *ShipmentManager) CalcMakeupShippingFeesByFfm(ctx context.Context, args *CalcMakeupShippingFeesByFfmArgs) (*CalcMakeupShippingFeesByFfmResponse, error) {
	ffm := args.Fulfillment
	connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
	driver, err := m.GetShipmentDriver(ctx, connectionID, ffm.ShopID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "invalid connection (ffm_id = %v)", ffm.ID)
	}
	serviceID, err := driver.ParseServiceID(ffm.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	serviceArgs := &GetShippingServicesArgs{
		FromDistrictCode: ffm.AddressFrom.DistrictCode,
		FromProvinceCode: ffm.AddressFrom.ProvinceCode,
		FromWardCode:     ffm.AddressFrom.WardCode,
		ToDistrictCode:   ffm.AddressTo.DistrictCode,
		ToProvinceCode:   ffm.AddressTo.ProvinceCode,
		ToWardCode:       ffm.AddressTo.WardCode,
		ChargeableWeight: cm.CoalesceInt(args.Weight, ffm.TotalWeight),
		BasketValue:      ffm.BasketValue,
		CODAmount:        ffm.TotalCODAmount,
	}
	shipmentService, err := m.getShipmentService(ctx, serviceArgs, serviceID, connectionID, true)
	if err != nil {
		return nil, err
	}
	addFeeTypes := args.AdditionalFeeTypes
	if serviceArgs.IncludeInsurance && !shipping_fee_type.Contain(addFeeTypes, shipping_fee_type.Insurance) {
		addFeeTypes = append(addFeeTypes, shipping_fee_type.Insurance)
	}
	if shipping.IsStateReturn(args.State) && !shipping_fee_type.Contain(addFeeTypes, shipping_fee_type.Return) {
		addFeeTypes = append(addFeeTypes, shipping_fee_type.Return)
	}
	shipmentPriceListID := dot.ID(0)
	if ffm.ShipmentPriceInfo != nil {
		shipmentPriceListID = ffm.ShipmentPriceInfo.ShipmentPriceListID
	}
	query := &shipmentprice.CalculateShippingFeesQuery{
		AccountID:           ffm.ShopID,
		ShipmentPriceListID: shipmentPriceListID,
		FromDistrictCode:    ffm.AddressFrom.DistrictCode,
		FromProvinceCode:    ffm.AddressFrom.ProvinceCode,
		ToDistrictCode:      ffm.AddressTo.DistrictCode,
		ToProvinceCode:      ffm.AddressTo.ProvinceCode,
		ShipmentServiceID:   shipmentService.ID,
		ConnectionID:        connectionID,
		Weight:              cm.CoalesceInt(args.Weight, ffm.TotalWeight),
		BasketValue:         ffm.BasketValue,
		CODAmount:           ffm.TotalCODAmount,
		AdditionalFeeTypes:  addFeeTypes,
	}
	if err := m.shipmentPriceQS.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := &CalcMakeupShippingFeesByFfmResponse{
		ShipmentPriceID:     query.Result.ShipmentPriceID,
		ShipmentPriceListID: query.Result.ShipmentPriceListID,
		ShippingFeeLines:    shipmentpriceconvert.Convert_shipmentprice_ShippingFees_To_shipping_ShippingFeeLines(query.Result.FeeLines),
	}
	return res, nil
}
