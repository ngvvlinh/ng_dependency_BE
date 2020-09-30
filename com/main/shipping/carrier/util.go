package carrier

import (
	"fmt"

	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

func CheckShippingService(ffm *shipmodel.Fulfillment, services []*shippingsharemodel.AvailableShippingService) (service *shippingsharemodel.AvailableShippingService, _err error) {
	providerServiceID := ffm.ProviderServiceID
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
	if ffm.ShippingServiceFee != service.ServiceFee {
		return nil, cm.Errorf(cm.InvalidArgument, nil,
			"Số tiền phí giao hàng không hợp lệ cho dịch vụ %v: Phí trên đơn hàng %v, phí từ dịch vụ giao hàng: %v",
			service.Name, ffm.ShippingServiceFee, service.ServiceFee)
	}

	if service.ShipmentServiceInfo != nil {
		if !service.ShipmentServiceInfo.IsAvailable {
			return nil, cm.Errorf(cm.InvalidArgument, nil, service.ShipmentServiceInfo.ErrorMessage)
		}
	}

	return service, nil
}

// filterShipmentServicesByCode
//
// Filter theo `ed_code` trong shipment_service.
// Là mã do admin TopShip định nghĩa cho từng gói dịch vụ.
// Yêu cầu chỉ trả về 1 gói duy nhất theo mã TopShip định nghĩa, với giá rẻ nhất từ NVC
func filterShipmentServicesByEdCode(services []*shippingsharemodel.AvailableShippingService) []*shippingsharemodel.AvailableShippingService {
	type UniqueService struct {
		service      *shippingsharemodel.AvailableShippingService
		carrierPrice int
	}
	var mapServices = make(map[string]*UniqueService)
	res := []*shippingsharemodel.AvailableShippingService{}
	for _, s := range services {
		if s.ShipmentServiceInfo == nil || s.ShipmentServiceInfo.Code == "" {
			res = append(res, s)
			continue
		}
		code := s.ShipmentServiceInfo.Code
		if mapServices[code] == nil {
			_service := &UniqueService{service: s}
			if s.ShipmentPriceInfo != nil {
				_service.carrierPrice = s.ShipmentPriceInfo.OriginFee
			} else {
				_service.carrierPrice = s.ServiceFee
			}
			mapServices[code] = _service
			continue
		}

		if !s.ShipmentServiceInfo.IsAvailable {
			continue
		}
		carrierPrice := getCarrierPriceFromService(s)
		currentService := mapServices[code]
		if !currentService.service.ShipmentServiceInfo.IsAvailable ||
			currentService.carrierPrice > carrierPrice {
			mapServices[code] = &UniqueService{
				service:      s,
				carrierPrice: carrierPrice,
			}
		}
	}

	if len(mapServices) == 0 {
		return res
	}

	for _, s := range mapServices {
		res = append(res, s.service)
	}
	return res
}

func getCarrierPriceFromService(s *shippingsharemodel.AvailableShippingService) int {
	if s.ShipmentPriceInfo == nil {
		return s.ServiceFee
	}
	return s.ShipmentPriceInfo.OriginFee
}

// CompactServices Loại bỏ các service không sử dụng
// Trường hợp:
// - Có gói TopShip: chỉ sử dụng gói TopShip
// - Mỗi NVC phải có 2 dịch vụ: Nhanh và Chuẩn, ưu tiên gói TopShip
// - Không có gói TopShip: Sử dụng gói của NVC như bình thường
func CompactServices(services []*shippingsharemodel.AvailableShippingService) []*shippingsharemodel.AvailableShippingService {
	var res []*shippingsharemodel.AvailableShippingService
	carrierServicesIndex := make(map[string][]*shippingsharemodel.AvailableShippingService)
	for _, s := range services {
		connectionID := dot.ID(0)
		if s.ConnectionInfo != nil {
			connectionID = s.ConnectionInfo.ID
		}
		key := fmt.Sprintf("%v_%v_%v", s.Provider.String(), s.Name, connectionID)
		carrierServicesIndex[key] = append(carrierServicesIndex[key], s)
	}
	for _, carrierServices := range carrierServicesIndex {
		var ss []*shippingsharemodel.AvailableShippingService
		for _, s := range carrierServices {
			if s.Source == model.TypeShippingSourceEtop {
				ss = append(ss, s)
			}
		}
		if len(ss) > 0 {
			res = append(res, ss...)
		} else {
			res = append(res, carrierServices...)
		}
	}
	return res
}
