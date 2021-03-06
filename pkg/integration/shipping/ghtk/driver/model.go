package driver

import (
	"o.o/api/main/location"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
)

type CalcShippingFeeArgs struct {
	ArbitraryID  int64
	FromProvince *location.Province
	FromDistrict *location.District
	ToProvince   *location.Province
	ToDistrict   *location.District
	Request      *ghtkclient.CalcShippingFeeRequest
}
