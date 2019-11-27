package carrier

import (
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/logic/etop_shipping_price"
	"etop.vn/backend/pkg/etop/model"
)

func GetEtopServiceFromSeviceCode(shippingServiceCode string, shippingServiceFee int, services []*model.AvailableShippingService) (etopService *model.AvailableShippingService, err error) {
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

func CheckShippingService(ffm *shipmodel.Fulfillment, services []*model.AvailableShippingService) (service *model.AvailableShippingService, _err error) {
	providerServiceID := cm.Coalesce(ffm.ProviderServiceID)
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
	return service, nil
}
