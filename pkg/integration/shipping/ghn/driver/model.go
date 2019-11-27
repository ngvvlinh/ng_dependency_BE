package driver

import (
	"etop.vn/api/main/location"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
)

type CalcShippingFeeArgs struct {
	ArbitraryID  int64
	FromProvince *location.Province
	FromDistrict *location.District
	ToProvince   *location.Province
	ToDistrict   *location.District
	Request      *ghnclient.FindAvailableServicesRequest
}

type RegisterWebhookForClientArgs struct {
	TokenClients []string
	// URLCallback: url must has ssl
	URLCallback string
}
