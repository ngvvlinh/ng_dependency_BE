package driver

import (
	"etop.vn/api/top/types/etc/shipping_provider"
	shippingsharemodel "etop.vn/backend/com/main/shipping/sharemodel"
	etopmodel "etop.vn/backend/pkg/etop/model"
	directclient "etop.vn/backend/pkg/integration/shipping/direct/client"
)

func toShippingFeeLine(line *directclient.ShippingFeeLine) *shippingsharemodel.ShippingFeeLine {
	if line == nil {
		return nil
	}
	return &shippingsharemodel.ShippingFeeLine{
		ShippingFeeType: line.ShippingFeeType,
		Cost:            line.Cost.Int(),
	}
}

func toShippingFeeLines(lines []*directclient.ShippingFeeLine) (res []*shippingsharemodel.ShippingFeeLine) {
	for _, line := range lines {
		res = append(res, toShippingFeeLine(line))
	}
	return
}

func toAvailableShippingService(s *directclient.ShippingService) *etopmodel.AvailableShippingService {
	if s == nil {
		return nil
	}
	return &etopmodel.AvailableShippingService{
		Name:               s.Name.String(),
		ServiceFee:         s.ServiceFee.Int(),
		ShippingFeeMain:    s.ServiceFeeMain.Int(),
		ProviderServiceID:  s.ServiceCode.String(),
		ExpectedPickAt:     s.ExpectedPickAt.ToTime(),
		ExpectedDeliveryAt: s.ExpectedDeliveryAt.ToTime(),
		Provider:           shipping_provider.Partner,
	}
}

func toAvailableShippingServices(services []*directclient.ShippingService) (res []*etopmodel.AvailableShippingService) {
	for _, s := range services {
		res = append(res, toAvailableShippingService(s))
	}
	return
}
