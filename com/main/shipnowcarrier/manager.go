package shipnowcarrier

import (
	"context"

	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	com "o.o/backend/com/main"
	shipnowstore "o.o/backend/com/main/shipnow/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
)

const MinShopBalance = -200000

var _ carrier.Manager = &ShipnowManager{}

type Carrier struct {
	Code carriertypes.ShipnowCarrier
	ShipnowCarrier
	ShipnowCarrierAccount
}

type ShipnowManager struct {
	drivers      []*Carrier
	location     location.QueryBus
	shipnowQuery shipnow.QueryBus
	store        shipnowstore.ShipnowStoreFactory

	MoneyTxStore sqlstore.MoneyTxStoreInterface
}

func NewManager(
	db com.MainDB,
	locationBus location.QueryBus,
	shipnowQuery shipnow.QueryBus,
	carriers []*Carrier,
	MoneyTxStore sqlstore.MoneyTxStoreInterface,
) *ShipnowManager {
	return &ShipnowManager{
		drivers:      carriers,
		location:     locationBus,
		shipnowQuery: shipnowQuery,
		store:        shipnowstore.NewShipnowStore(db),
		MoneyTxStore: MoneyTxStore,
	}
}

func (ctrl *ShipnowManager) CreateExternalShipnow(ctx context.Context, cmd *carrier.CreateExternalShipnowCommand) (*carrier.ExternalShipnow, error) {
	// check balance of shop
	// if balance < MinShopBalance => can not create order
	// TODO: plus balance with current order'store value
	// TODO: move to credit aggregate
	{
		query := &model.GetBalanceShopCommand{
			ShopID: cmd.ShopID,
		}
		if err := ctrl.MoneyTxStore.CalcBalanceShop(ctx, query); err != nil {
			return nil, err
		}
		balance := query.Result.Amount
		if balance < MinShopBalance {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Bạn đang nợ cước số tiền: %v. Vui lòng liên hệ ETOP để xử lý.", balance)
		}
	}

	return ctrl.createSingleFulfillment(ctx, cmd)
}

func (ctrl *ShipnowManager) CancelExternalShipping(ctx context.Context, cmd *carrier.CancelExternalShipnowCommand) error {
	shipnowCarrier, err := ctrl.GetShipnowCarrierDriver(cmd.Carrier)
	if err != nil {
		return nil
	}
	if err := shipnowCarrier.CancelExternalShipnow(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (ctrl *ShipnowManager) createSingleFulfillment(ctx context.Context, cmd *carrier.CreateExternalShipnowCommand) (externalShipnow *carrier.ExternalShipnow, _err error) {
	query := &shipnow.GetShipnowFulfillmentQuery{
		ID:     cmd.ShipnowFulfillmentID,
		ShopID: cmd.ShopID,
	}
	if err := ctrl.shipnowQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	ffm := query.Result.ShipnowFulfillment
	if err := ctrl.ValidateShipnowAddress(ffm); err != nil {
		return nil, err
	}

	shipnowCarrier, err := ctrl.GetShipnowCarrierDriver(ffm.Carrier)
	if err != nil {
		return nil, nil
	}

	args := GetShipnowServiceArgs{
		ShopID:         cmd.ShopID,
		PickupAddress:  ffm.PickupAddress,
		DeliveryPoints: ffm.DeliveryPoints,
	}
	services, err := shipnowCarrier.GetShipnowServices(ctx, args)
	if err != nil {
		return nil, err
	}

	service, err := ctrl.CheckShippingService(ffm, services)
	if err != nil {
		return nil, err
	}

	externalShipnow, err = shipnowCarrier.CreateExternalShipnow(ctx, cmd, service)
	if err != nil {
		return nil, err
	}
	externalShipnow.Service = service
	return externalShipnow, nil
}

func (ctrl *ShipnowManager) GetShipnowCarrierDriver(c carriertypes.ShipnowCarrier) (*Carrier, error) {
	for _, d := range ctrl.drivers {
		if d.Code == c {
			return d, nil
		}
	}
	return nil, cm.Errorf(cm.InvalidArgument, nil, "Đơn vị vận chuyển không hợp lệ")
}

func (ctrl *ShipnowManager) CheckShippingService(ffm *shipnow.ShipnowFulfillment, services []*shipnowtypes.ShipnowService) (service *shipnowtypes.ShipnowService, err error) {
	if ffm.ShippingServiceCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần chọn gói dịch vụ giao hàng")
	}
	shipnowCarrier, err := ctrl.GetShipnowCarrierDriver(ffm.Carrier)
	if err != nil {
		return nil, err
	}
	shippingServiceID, err := shipnowCarrier.ParseServiceCode(ffm.ShippingServiceCode)
	for _, s := range services {
		sID, _ := shipnowCarrier.ParseServiceCode(s.Code)
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

func (ctrl *ShipnowManager) VerifyAddress(addr *ordertypes.Address, requireWard bool) (*location.Province, *location.District, *location.Ward, error) {
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
	if err := ctrl.location.DispatchAll(context.TODO(), query); err != nil {
		return nil, nil, nil, err
	}
	if addr.Coordinates == nil || addr.Coordinates.Latitude == 0 || addr.Coordinates.Longitude == 0 {
		return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp Latitude và Longitude")
	}
	loc := query.Result
	return loc.Province, loc.District, loc.Ward, nil
}

func (ctrl *ShipnowManager) ValidateShipnowAddress(ffm *shipnow.ShipnowFulfillment) error {
	if _, _, _, err := ctrl.VerifyAddress(ffm.PickupAddress, false); err != nil {
		return err
	}

	for _, point := range ffm.DeliveryPoints {
		if _, _, _, err := ctrl.VerifyAddress(point.ShippingAddress, false); err != nil {
			return err
		}
	}
	return nil
}

func (ctrl *ShipnowManager) GetExternalShipnowServices(ctx context.Context, cmd *carrier.GetExternalShipnowServicesCommand) ([]*shipnowtypes.ShipnowService, error) {
	args := GetShipnowServiceArgs{
		ShopID:         cmd.ShopID,
		PickupAddress:  cmd.PickupAddress,
		DeliveryPoints: cmd.DeliveryPoints,
	}
	var res []*shipnowtypes.ShipnowService

	n := len(ctrl.drivers)
	ch := make(chan []*shipnowtypes.ShipnowService, n)
	for _, driver := range ctrl.drivers {
		driver := driver // closure
		go func() {
			defer cm.RecoverAndLog()
			var services []*shipnowtypes.ShipnowService
			var err error
			defer func() { sendServices(ch, services, err) }()
			services, err = driver.GetShipnowServices(ctx, args)
		}()
	}

	for i := 0; i < n; i++ {
		res = append(res, <-ch...)
	}
	if len(res) == 0 {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Tuyến giao hàng không được hỗ trợ bởi đơn bị vận chuyển nào")
	}
	return res, nil
}

func sendServices(ch chan<- []*shipnowtypes.ShipnowService, services []*shipnowtypes.ShipnowService, err error) {
	if err == nil {
		ch <- services
	} else {
		ch <- nil
	}
}

func (ctrl *ShipnowManager) RegisterExternalAccount(ctx context.Context, cmd *carrier.RegisterExternalAccountCommand) (*carrier.RegisterExternalAccountResult, error) {
	shipnowCarrier, err := ctrl.GetShipnowCarrierDriver(cmd.Carrier)
	if err != nil {
		return nil, err
	}
	args := &RegisterExternalAccountArgs{
		Phone:   cmd.Phone,
		Name:    cmd.Name,
		Address: cmd.Address,
	}
	return shipnowCarrier.RegisterExternalAccount(ctx, args)
}

func (ctrl *ShipnowManager) GetExternalAccount(ctx context.Context, cmd *carrier.GetExternalAccountCommand) (*carrier.ExternalAccount, error) {
	shipnowCarrier, err := ctrl.GetShipnowCarrierDriver(cmd.Carrier)
	if err != nil {
		return nil, err
	}
	args := &GetExternalAccountArgs{
		OwnerID: cmd.OwnerID,
	}
	return shipnowCarrier.GetExternalAccount(ctx, args)
}

func (ctrl *ShipnowManager) VerifyExternalAccount(ctx context.Context, cmd *carrier.VerifyExternalAccountCommand) (*carrier.VerifyExternalAccountResult, error) {
	shipnowCarrier, err := ctrl.GetShipnowCarrierDriver(cmd.Carrier)
	if err != nil {
		return nil, err
	}
	args := &VerifyExternalAccountArgs{
		OwnerID: cmd.OwnerID,
	}
	return shipnowCarrier.VerifyExternalAccount(ctx, args)
}
