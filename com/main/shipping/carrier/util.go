package carrier

import (
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/logic/etop_shipping_price"
	"o.o/backend/pkg/etop/model"
)

func GetEtopServiceFromSeviceCode(shippingServiceCode string, shippingServiceFee int, services []*shippingsharemodel.AvailableShippingService) (etopService *shippingsharemodel.AvailableShippingService, err error) {
	if shippingServiceCode == "" {
		return nil, cm.Error(cm.InvalidArgument, "ShopShipping is invalid", nil)
	}

	sType, isEtopService := etop_shipping_price.ParseEtopServiceCode(shippingServiceCode)
	if !isEtopService {
		return nil, cm.Error(cm.InvalidArgument, "ShippingServiceCode is invalid", nil)
	}
	for _, service := range services {
		if service.Name == sType && service.ServiceFee == shippingServiceFee && service.Source == model.TypeShippingSourceEtop {
			etopService = service
			return etopService, nil
		}
	}
	return nil, cm.Error(cm.NotFound, "Không có gói vận chuyển phù hợp", nil)
}

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
