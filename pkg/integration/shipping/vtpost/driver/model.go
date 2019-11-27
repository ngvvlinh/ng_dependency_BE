package driver

import (
	"etop.vn/api/main/location"
	vtpostclient "etop.vn/backend/pkg/integration/shipping/vtpost/client"
)

type CalcShippingFeeArgs struct {
	ArbitraryID  int64
	FromProvince *location.Province
	FromDistrict *location.District
	ToProvince   *location.Province
	ToDistrict   *location.District
	Request      *vtpostclient.CalcShippingFeeAllServicesRequest
}
