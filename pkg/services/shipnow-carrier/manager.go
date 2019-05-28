package shipnow_carrier

import (
	"context"

	"etop.vn/api/main/location"
	"etop.vn/api/main/shipnow/carrier"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

const MinShopBalance = -200000

var _ carrier.Manager = &ShipnowManager{}

type ShipnowManager struct {
	ahamove  ShipnowCarrier
	location location.Bus
	s        sqlstore.ShipnowStoreFactory
}

func NewManager(db cmsql.Database, locationBus location.Bus, ahamoveCarrier ShipnowCarrier) *ShipnowManager {
	return &ShipnowManager{
		ahamove:  ahamoveCarrier,
		location: locationBus,
		s:        sqlstore.NewShipnowStore(db),
	}
}

func (ShipnowManager) CreateExternalShipping(ctx context.Context, ffm *carrier.CreateExternalShipnowCommand) error {
	panic("implement me")
}

func (ShipnowManager) CancelExternalShipping(ctx context.Context, ffm *carrier.CancelExternalShipnowCommand) error {
	panic("implement me")
}

//
// func (ctrl *ShipnowManager) CreateExternalShipping(ctx context.Context, ffm *model.ShipnowFulfillment) error {
// 	// check balance of shop
// 	// if balance < MinShopBalance => can not create order
// 	// TODO: plus balance with current order's value
// 	{
// 		query := &etopmodel.GetBalanceShopCommand{
// 			ShopID: ffm.ShopID,
// 		}
// 		if err := bus.Dispatch(ctx, query); err != nil {
// 			return err
// 		}
// 		balance := query.Result.Amount
// 		if balance < MinShopBalance {
// 			return cm.Errorf(cm.FailedPrecondition, nil, "Bạn đang nợ cước số tiền: %v. Vui lòng liên hệ ETOP để xử lý.", balance)
// 		}
// 	}
// 	return ctrl.createSingleFulfillment(ctx, ffm)
// }
//
// func (ctrl *ShipnowManager) createSingleFulfillment(ctx context.Context, ffm *model.ShipnowFulfillment) (_err error) {
// 	carrier := ffm.Carrier
// 	shipnowCarrier := ctrl.GetShipnowCarrierDriver(carrier)
// 	if shipnowCarrier == nil {
// 		return cm.Errorf(cm.InvalidArgument, nil, "invalid carrier")
// 	}
//
// 	{
// 		updateFfm := &model.ShipnowFulfillment{
// 			ID:         ffm.ID,
// 			SyncStatus: etopmodel.S4SuperPos,
// 			SyncStates: &etopmodel.FulfillmentSyncStates{
// 				TrySyncAt:         time.Now(),
// 				NextShippingState: etopmodel.StateCreated,
// 			},
// 		}
// 		if err := ctrl.s.Update(shipnowconvert.Shipnow(updateFfm)); err != nil {
// 			return err
// 		}
// 	}
// 	defer func() {
// 		if _err == nil {
// 			return
// 		}
// 		updateFfm := &model.ShipnowFulfillment{
// 			ID:         ffm.ID,
// 			SyncStatus: etopmodel.S4Negative,
// 			SyncStates: &etopmodel.FulfillmentSyncStates{
// 				TrySyncAt:         time.Now(),
// 				Error:             etopmodel.ToError(_err),
// 				NextShippingState: etopmodel.StateCreated,
// 			},
// 		}
// 		if err := ctrl.s.Update(shipnowconvert.Shipnow(updateFfm)); err != nil {
// 			return
// 		}
// 	}()
//
// 	args := GetShippingServiceArgs{
// 		AccountID:      ffm.ShopID,
// 		DeliveryPoints: shipnowconvert.DeliveryPoints(ffm.DeliveryPoints),
// 	}
// 	services, err := shipnowCarrier.GetShippingServices(ctx, args)
// 	if err != nil {
// 		return err
// 	}
//
// 	service, err := checkShippingService(ffm, services)
// 	if err != nil {
// 		return err
// 	}
// 	ffmToUpdate, err := shipnowCarrier.CreateFulfillment(ctx, ffm, service)
// 	if err := ctrl.s.Update(shipnowconvert.Shipnow(ffmToUpdate)); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (ctrl *ShipnowManager) CancelExternalShipping(ctx context.Context, ffm *shipnow.ShipnowFulfillment) error {
// 	carrier := ffm.Carrier
// 	shipnowCarrier := ctrl.GetShipnowCarrierDriver(carrier)
// 	if shipnowCarrier == nil {
// 		return cm.Errorf(cm.InvalidArgument, nil, "invalid carrier")
// 	}
// 	if err := shipnowCarrier.CancelFulfillment(ctx, ffm); err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (ctrl *ShipnowManager) GetShipnowCarrierDriver(carrier model.Carrier) ShipnowCarrier {
// 	switch carrier {
// 	case model.ahamove:
// 		return ctrl.ahamove
// 	default:
// 		return nil
// 	}
// }
//
// func checkShippingService(ffm *model.ShipnowFulfillment, services []*etopmodel.AvailableShippingService) (service *etopmodel.AvailableShippingService, err error) {
// 	if ffm.ShippingServiceCode == "" {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần chọn gói dịch vụ giao hàng")
// 	}
// 	providerServiceID := ffm.ShippingServiceCode
// 	for _, s := range services {
// 		if s.ProviderServiceID == providerServiceID {
// 			service = s
// 		}
// 	}
// 	if service == nil {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil, "Gói dịch vụ đã chọn không hợp lệ")
// 	}
// 	if ffm.ShippingServiceFee != int32(service.ServiceFee) {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số tiền phí giao hàng không hợp lệ cho dịch vụ %v: Phí trên đơn hàng %v, phí từ dịch vụ giao hàng: %v", service.Name, ffm.ShippingServiceFee, service.ServiceFee)
// 	}
// 	return service, nil
// }
//
// func (ctrl *ShipnowManager) VerifyDistrictCode(addr *etopmodel.Address) (*location.District, *location.Province, error) {
// 	if addr == nil {
// 		return nil, nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
// 	}
// 	if addr.DistrictCode == "" {
// 		return nil, nil, cm.Error(cm.InvalidArgument, cm.F(
// 			`Địa chỉ %v, %v không thể được xác định bởi hệ thống.`,
// 			addr.District, addr.Province,
// 		), nil)
// 	}
//
// 	query := &location.GetLocationQuery{DistrictCode: addr.DistrictCode}
// 	if err := ctrl.location.Dispatch(context.TODO(), query); err != nil {
// 		return nil, nil, err
// 	}
//
// 	district := query.Result.District
// 	if district.Extra.GhnId == 0 {
// 		return nil, nil, cm.Errorf(cm.InvalidArgument, nil,
// 			"Địa chỉ %v, %v không thể được giao bởi dịch vụ vận chuyển.",
// 			addr.District, addr.Province,
// 		)
// 	}
// 	return district, query.Result.Province, nil
// }
//
// func (ctrl *ShipnowManager) VerifyWardCode(addr *etopmodel.Address) (*location.Ward, error) {
// 	if addr == nil {
// 		return nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
// 	}
// 	if addr.WardCode == "" {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil,
// 			`Thiếu thông tin phường xã (%v, %v).`,
// 			addr.District, addr.Province,
// 		)
// 	}
//
// 	query := &location.GetLocationQuery{WardCode: addr.WardCode}
// 	if err := ctrl.location.Dispatch(context.TODO(), query); err != nil {
// 		return nil, err
// 	}
// 	return query.Result.Ward, nil
// }
//
// func (ctrl *ShipnowManager) VerifyProvinceCode(addr *etopmodel.Address) (*location.Province, error) {
// 	if addr == nil {
// 		return nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
// 	}
// 	if addr.ProvinceCode == "" {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil,
// 			`Địa chỉ %v, %v không thể được xác định bởi hệ thống.`,
// 			addr.District, addr.Province,
// 		)
// 	}
//
// 	query := &location.GetLocationQuery{ProvinceCode: addr.ProvinceCode}
// 	if err := ctrl.location.Dispatch(context.TODO(), query); err != nil {
// 		return nil, err
// 	}
// 	return query.Result.Province, nil
// }
//
// func (ctrl *ShipnowManager) VerifyAddress(addr *etopmodel.Address, requireWard bool) (*location.Province, *location.District, *location.Ward, error) {
// 	if addr == nil {
// 		return nil, nil, nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
// 	}
// 	if addr.ProvinceCode == "" || addr.DistrictCode == "" {
// 		return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil,
// 			`Địa chỉ %v, %v không thể được xác định bởi hệ thống.`,
// 			addr.District, addr.Province,
// 		)
// 	}
// 	query := &location.GetLocationQuery{
// 		ProvinceCode: addr.ProvinceCode,
// 		DistrictCode: addr.DistrictCode,
// 	}
// 	if requireWard {
// 		if addr.WardCode == "" {
// 			return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil,
// 				`Cần cung cấp thông tin phường xã hợp lệ`)
// 		}
// 		query.WardCode = addr.WardCode
// 	}
// 	if err := ctrl.location.DispatchAll(context.TODO(), query); err != nil {
// 		return nil, nil, nil, err
// 	}
// 	if addr.Coordinates == nil || addr.Coordinates.Latitude == 0 || addr.Coordinates.Longitude == 0 {
// 		return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp Latitude và Longitude")
// 	}
// 	loc := query.Result
// 	return loc.Province, loc.District, loc.Ward, nil
// }
