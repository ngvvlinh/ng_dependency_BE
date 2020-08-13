package carrier

import (
	"context"
	"strings"
	"sync"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/main/ordering"
	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
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
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

const (
	PrefixMakeupPriceCode = "###"
)

type ShipmentManager struct {
	locationQS           location.QueryBus
	connectionQS         connectioning.QueryBus
	connectionAggr       connectioning.CommandBus
	env                  string
	shipmentServiceQS    shipmentservice.QueryBus
	shipmentPriceQS      shipmentprice.QueryBus
	shippingQS           shipping.QueryBus
	priceListPromotionQS pricelistpromotion.QueryBus
	ghnWebhookEndpoint   string
	ConnectionManager    *connectionmanager.ConnectionManager

	webhookEndpoints carriertypes.ConfigEndpoints
	carrierDriver    carriertypes.Driver

	eventBus capi.EventBus
}

func NewShipmentManager(
	eventBus capi.EventBus,
	locationQS location.QueryBus,
	connectionQS connectioning.QueryBus,
	connectionAggr connectioning.CommandBus,
	shipmentServiceQS shipmentservice.QueryBus,
	shipmentPriceQS shipmentprice.QueryBus,
	priceListPromotionQS pricelistpromotion.QueryBus,
	cfg carriertypes.Config,
	carrierDriver carriertypes.Driver,
	connectionManager *connectionmanager.ConnectionManager,
) (*ShipmentManager, error) {
	sm := &ShipmentManager{
		eventBus:             eventBus,
		locationQS:           locationQS,
		connectionQS:         connectionQS,
		connectionAggr:       connectionAggr,
		env:                  cmenv.PartnerEnv(),
		shipmentServiceQS:    shipmentServiceQS,
		shipmentPriceQS:      shipmentPriceQS,
		priceListPromotionQS: priceListPromotionQS,
		ConnectionManager:    connectionManager,
		webhookEndpoints:     cfg.Endpoints,
		carrierDriver:        carrierDriver,
	}
	return sm, nil
}

func (m *ShipmentManager) getShipmentDriver(ctx context.Context, connectionID dot.ID, shopID dot.ID) (carriertypes.ShipmentCarrier, error) {
	connection, err := m.ConnectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}
	_shopID := shopID
	if connection.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		// ignore shopID
		_shopID = 0
	}
	shopConnection, err := m.ConnectionManager.GetShopConnection(ctx, connectionID, _shopID)
	if err != nil {
		return nil, err
	}

	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid (check status or token)")
	}
	return m.carrierDriver.GetShipmentDriver(m.env, m.locationQS, connection, shopConnection, m.webhookEndpoints)
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
		EventMeta:    meta.NewEvent(),
		ShopID:       ffm.ShopID,
		FromAddress:  addressconvert.Convert_addressmodel_Address_orderingtypes_Address(ffm.AddressFrom, nil),
		ShippingFee:  ffm.ShippingServiceFee,
		ConnectionID: ffm.ConnectionID,
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
		IncludeInsurance: ffm.IncludeInsurance.Apply(false),
		InsuranceValue:   ffm.InsuranceValue,
		BasketValue:      ffm.BasketValue,
		CODAmount:        ffm.TotalCODAmount,
		Coupon:           ffm.Coupon,
	}
	allServices, err := m.GetShipmentServicesAndMakeupPrice(ctx, args, ffm.ConnectionID)
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

	_args := args.ToShipmentServiceArgs(ffm.ConnectionID, ffm.ShopID)
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
	if err := m.locationQS.Dispatch(context.TODO(), query); err != nil {
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
	if err := m.connectionQS.Dispatch(ctx, query); err != nil {
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
		Identifier: args.Identifier,
		Password:   args.Password,
		OTP:        args.OTP,
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
		Email:    args.Identifier,
		Password: args.Password,
		Phone:    args.Phone,
		Province: args.Province,
		District: args.District,
		Address:  args.Address,
	}
	return driver.SignUp(ctx, cmd)
}

func (m *ShipmentManager) getDriverByEtopAffiliateAccount(ctx context.Context, connectionID dot.ID) (carriertypes.ShipmentCarrier, error) {
	conn, err := m.ConnectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}

	// Chỉ có method=direct mới được login
	if conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Do not support this feature for this connection")
	}
	return m.carrierDriver.GetAffiliateShipmentDriver(m.env, m.locationQS, conn, m.webhookEndpoints)
}

func (m *ShipmentManager) RefreshFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) (updateFfm *shipmodel.Fulfillment, err error) {
	connectionID := shipping.GetConnectionID(ffm.ConnectionID, ffm.ShippingProvider)
	driver, err := m.getShipmentDriver(ctx, connectionID, ffm.ShopID)
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
	driver, err := m.getShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
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
	driver, err := m.getShipmentDriver(ctx, ffm.ConnectionID, ffm.ShopID)
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
	driver, err := m.getShipmentDriver(ctx, connID, accountID)
	if err != nil {
		// ll.Error("Driver shipment không hợp lệ", l.ID("shopID", accountID), l.ID("connectionID", connID), l.Error(err))
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
		// Nếu không có cấu hình giá (shipment_price) của mã dịch vụ (shipment_service_id) trong bảng giá (mã lỗi not_found)
		// Trả về kết quả của NVC
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
		wardCode = args.ToWardCode
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

func (m *ShipmentManager) makeupPriceByShipmentPrice(ctx context.Context, service *shippingsharemodel.AvailableShippingService, args *GetShippingServicesArgs) error {
	if service.ShipmentServiceInfo == nil || service.ShipmentServiceInfo.ID == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Thiếu shipment service.")
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
	driver, err := m.getShipmentDriver(ctx, connectionID, ffm.ShopID)
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
