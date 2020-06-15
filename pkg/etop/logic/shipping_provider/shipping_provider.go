package shipping_provider

import (
	"context"
	"strings"
	"time"

	"o.o/api/main/location"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	addressconvert "o.o/backend/com/main/address/convert"
	addressmodel "o.o/backend/com/main/address/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/syncgroup"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/logic/etop_shipping_price"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
)

type codeCarrier struct {
	code   shipping_provider.ShippingProvider
	driver CarrierDriver
}

type CarrierManager struct {
	carriers []codeCarrier
	location location.QueryBus
	eventBus capi.EventBus
}

func NewCtrl(eventBus capi.EventBus, locationBus location.QueryBus, carriers []CarrierDriver) *CarrierManager {
	p := &CarrierManager{
		eventBus: eventBus,
		location: locationBus,
	}
	for _, c := range carriers {
		p.carriers = append(p.carriers, codeCarrier{code: c.Code(), driver: c})
	}
	return p
}

func (ctrl *CarrierManager) CreateExternalShipping(ctx context.Context, order *ordermodel.Order, ffms []*shipmodel.Fulfillment) error {
	return ctrl.createFulfillments(ctx, order, ffms)
}

func (ctrl *CarrierManager) createFulfillments(ctx context.Context, order *ordermodel.Order, ffms []*shipmodel.Fulfillment) error {
	var err error
	g := syncgroup.New(len(ffms))
	for _, ffm := range ffms {
		ffm := ffm // capture closure
		g.Go(func() error { return ctrl.createSingleFulfillment(ctx, order, ffm) })
	}
	errs := g.Wait()
	if errs.IsAll() {
		err = errs[0]
	}
	return err
}

func (ctrl *CarrierManager) GetShippingProviderDriver(provider shipping_provider.ShippingProvider) CarrierDriver {
	for _, c := range ctrl.carriers {
		if c.code == provider {
			return c.driver
		}
	}
	return nil
}

func (ctrl *CarrierManager) createSingleFulfillment(ctx context.Context, order *ordermodel.Order, ffm *shipmodel.Fulfillment) (_err error) {
	// TODO: handle case when ffm.shipping_provider is different with order.shipping_provider
	provider := order.ShopShipping.ShippingProvider

	if err := checkBlockCarrier(provider); err != nil {
		return err
	}

	shippingProvider := ctrl.GetShippingProviderDriver(provider)
	if shippingProvider == nil {
		return cm.Errorf(cm.Internal, nil, "invalid carrier")
	}

	// raise event to check balance
	event := &shippingcore.SingleFulfillmentCreatingEvent{
		EventMeta:    meta.NewEvent(),
		ShopID:       ffm.ShopID,
		FromAddress:  addressconvert.Convert_addressmodel_Address_orderingtypes_Address(ffm.AddressFrom, nil),
		ShippingFee:  order.ShopShipping.ExternalShippingFee,
		ConnectionID: ffm.ConnectionID,
	}
	if err := ctrl.eventBus.Publish(ctx, event); err != nil {
		return err
	}

	{
		// UpdateInfo status to pending
		updateFfm := &shipmodel.Fulfillment{
			ID:         ffm.ID,
			SyncStatus: status4.S,
			SyncStates: &shippingsharemodel.FulfillmentSyncStates{
				TrySyncAt:         time.Now(),
				NextShippingState: shipping.Created,
			},
		}
		cmd := &shipmodelx.UpdateFulfillmentCommand{Fulfillment: updateFfm}
		if err := bus.Dispatch(ctx, cmd); err != nil {
			return err
		}
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

				NextShippingState: shipping.Created,
			},
		}
		cmd := &shipmodelx.UpdateFulfillmentCommand{Fulfillment: updateFfm2}

		// Keep the original error
		_ = bus.Dispatch(ctx, cmd)
	}()

	shopShipping := order.ShopShipping

	fromDistrict, _, err := ctrl.VerifyDistrictCode(ffm.AddressFrom)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "FromDistrictCode: %v", err)
	}
	toDistrict, _, err := ctrl.VerifyDistrictCode(ffm.AddressTo)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "ToDistrictCode: %v", err)
	}

	weight, length, width, height := order.TotalWeight, 10, 10, 10
	if weight == 0 {
		weight = 100
	}

	// note := GetShippingProviderNote(order, ffm)
	noteCode := order.GhnNoteCode
	if noteCode == 0 {
		// harcode
		noteCode = ghn_note_code.CHOXEMHANGKHONGTHU
	}

	args := GetShippingServicesArgs{
		AccountID:        order.ShopID,
		FromDistrictCode: fromDistrict.Code,
		ToDistrictCode:   toDistrict.Code,
		ChargeableWeight: weight,
		Length:           length,
		Width:            width,
		Height:           height,
		IncludeInsurance: shopShipping.IncludeInsurance, // TODO: fix it
		BasketValue:      order.BasketValue,
		CODAmount:        ffm.TotalCODAmount,
	}

	allServices, err := shippingProvider.GetAllShippingServices(ctx, args)
	if err != nil {
		return err
	}

	// check if etop package
	var etopService, providerService *shippingsharemodel.AvailableShippingService
	sType, isEtopService := etop_shipping_price.ParseEtopServiceCode(shopShipping.ProviderServiceID)
	if isEtopService {
		// ETOP serivce
		// => Get cheapest provider service
		etopService, err = GetEtopServiceFromShopShipping(order.ShopShipping, allServices)
		if err != nil {
			return err
		}

		providerService = GetCheapestService(allServices, sType)
		if providerService == nil {
			return cm.Error(cm.InvalidArgument, "Không có gói vận chuyển phù hợp.", nil)
		}
	} else {
		// Provider service
		// => Check price
		// => Get this service
		providerService, err = CheckShippingService(order, allServices)
		if err != nil {
			return err
		}
	}

	ffmToUpdate, err := shippingProvider.CreateFulfillment(ctx, order, ffm, args, providerService)
	if err != nil {
		return err
	}

	if etopService != nil {
		ffmToUpdate.ApplyEtopPrice(etopService.ShippingFeeMain)
		ffmToUpdate.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(ffmToUpdate.ProviderShippingFeeLines, ffmToUpdate.EtopPriceRule, dot.Int(ffmToUpdate.EtopAdjustedShippingFeeMain))
	}

	updateCmd := &shipmodelx.UpdateFulfillmentCommand{
		Fulfillment: ffmToUpdate,
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return cm.Trace(err)
	}
	return nil
}

func GetShippingProviderNote(order *ordermodel.Order, ffm *shipmodel.Fulfillment) string {
	noteB := strings.Builder{}
	if note := ffm.AddressFrom.Notes.GetFullNote(); note != "" {
		noteB.WriteString("Lấy hàng: ")
		noteB.WriteString(note)
		noteB.WriteString("\n")
	}
	if note := ffm.AddressTo.Notes.GetFullNote(); note != "" || order.ShippingNote != "" {
		noteB.WriteString("Giao hàng: ")
		if order.ShippingNote != "" {
			noteB.WriteString(order.ShippingNote)
			noteB.WriteString(". \n")
		}
		noteB.WriteString(note)
		noteB.WriteString("\n")
	}
	noteB.WriteString("Giao hàng không thành công hoặc giao một phần, xin gọi lại cho shop. KHÔNG ĐƯỢC TỰ Ý HOÀN HÀNG khi chưa thông báo cho shop.")
	return noteB.String()
}

func CheckShippingService(order *ordermodel.Order, services []*shippingsharemodel.AvailableShippingService) (service *shippingsharemodel.AvailableShippingService, _err error) {
	if order.ShopShipping != nil {
		providerServiceID := cm.Coalesce(order.ShopShipping.ProviderServiceID, order.ShopShipping.ExternalServiceID)
		if providerServiceID == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần chọn gói dịch vụ giao hàng")
		}
		for _, s := range services {
			if s.ProviderServiceID == providerServiceID {
				service = s
			}
		}
		if service == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Gói dịch vụ giao hàng đã chọn không hợp lệ")
		}
		if order.ShopShipping.ExternalShippingFee != service.ServiceFee {
			return nil, cm.Errorf(cm.InvalidArgument, nil,
				"Số tiền phí giao hàng không hợp lệ cho dịch vụ %v: Phí trên đơn hàng %v, phí từ dịch vụ giao hàng: %v",
				service.Name, order.ShopShipping.ExternalShippingFee, service.ServiceFee)
		}
		return service, nil
	}
	return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần chọn gói dịch vụ giao hàng")
}

func (ctrl *CarrierManager) VerifyDistrictCode(addr *addressmodel.Address) (*location.District, *location.Province, error) {
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
	if err := ctrl.location.Dispatch(context.TODO(), query); err != nil {
		return nil, nil, err
	}

	district := query.Result.District
	if district.Extra.GhnId == 0 {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil,
			"Địa chỉ %v, %v không thể được giao bởi dịch vụ vận chuyển.",
			addr.District, addr.Province,
		)
	}
	return district, query.Result.Province, nil
}

func (ctrl *CarrierManager) VerifyWardCode(addr *addressmodel.Address) (*location.Ward, error) {
	if addr == nil {
		return nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
	}
	if addr.WardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Thiếu thông tin phường xã (%v, %v).`,
			addr.District, addr.Province,
		)
	}

	query := &location.GetLocationQuery{WardCode: addr.WardCode}
	if err := ctrl.location.Dispatch(context.TODO(), query); err != nil {
		return nil, err
	}
	return query.Result.Ward, nil
}

func (ctrl *CarrierManager) VerifyProvinceCode(addr *addressmodel.Address) (*location.Province, error) {
	if addr == nil {
		return nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
	}
	if addr.ProvinceCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			`Địa chỉ %v, %v không thể được xác định bởi hệ thống.`,
			addr.District, addr.Province,
		)
	}

	query := &location.GetLocationQuery{ProvinceCode: addr.ProvinceCode}
	if err := ctrl.location.Dispatch(context.TODO(), query); err != nil {
		return nil, err
	}
	return query.Result.Province, nil
}

func (ctrl *CarrierManager) VerifyAddress(addr *addressmodel.Address, requireWard bool) (*location.Province, *location.District, *location.Ward, error) {
	if addr == nil {
		return nil, nil, nil, cm.Errorf(cm.Internal, nil, "Địa chỉ không tồn tại")
	}
	if addr.ProvinceCode == "" || addr.DistrictCode == "" {
		return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil,
			`Địa chỉ %v, %v không thể được xác định bởi hệ thống.`,
			addr.District, addr.Province,
		)
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
	loc := query.Result
	return loc.Province, loc.District, loc.Ward, nil
}

func (ctrl *CarrierManager) ParseServiceCode(carrier shipping_provider.ShippingProvider, code string) (serviceName string, ok bool) {
	c := ctrl.GetShippingProviderDriver(carrier)
	if c == nil {
		return "", false
	}
	return c.ParseServiceCode(code)
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func CheckShippingFeeWithinDelta(providerShippingFee int, shippingFee int) bool {
	return abs(providerShippingFee-shippingFee) < 10
}

func GetEtopServiceFromShopShipping(shopShipping *ordermodel.OrderShipping, services []*shippingsharemodel.AvailableShippingService) (etopService *shippingsharemodel.AvailableShippingService, err error) {
	if shopShipping == nil || shopShipping.ProviderServiceID == "" {
		return nil, cm.Error(cm.InvalidArgument, "ShopShipping is invalid", nil)
	}

	sType, isEtopService := etop_shipping_price.ParseEtopServiceCode(shopShipping.ProviderServiceID)
	if !isEtopService {
		return nil, cm.Error(cm.InvalidArgument, "ProviderServiceID is invalid", nil)
	}
	for _, service := range services {
		if service.Name == sType && service.ServiceFee == shopShipping.ExternalShippingFee && service.Source == model.TypeShippingSourceEtop {
			etopService = service
			return etopService, nil
		}
	}
	return nil, cm.Error(cm.NotFound, "Không có gói vận chuyển phù hợp", nil)
}

// GetCheapestService chooses cheapest provider service (except etop service) using ServiceType (Nhanh | Chuan)
func GetCheapestService(services []*shippingsharemodel.AvailableShippingService, sType string) *shippingsharemodel.AvailableShippingService {
	var service *shippingsharemodel.AvailableShippingService
	for _, s := range services {
		if s.Source == model.TypeShippingSourceEtop || s.Name != sType {
			continue
		}
		if service == nil {
			service = s
			continue
		}
		if service.ServiceFee > s.ServiceFee {
			service = s
		}
	}
	return service
}
