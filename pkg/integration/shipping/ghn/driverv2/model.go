package driverv2

import (
	"o.o/api/main/location"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
)

type CalcShippingFeeArgs struct {
	ArbitraryID  int64
	FromProvince *location.Province
	FromDistrict *location.District
	ToProvince   *location.Province
	ToDistrict   *location.District
	Request      *ghnclient.FindAvailableServicesRequest
}
