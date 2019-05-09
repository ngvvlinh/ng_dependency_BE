package vtpost

import (
	"time"

	mdlocation "etop.vn/api/main/location"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	vtpostClient "etop.vn/backend/pkg/integration/vtpost/client"
)

type Config struct {
	Env string `yaml:"env"`

	AccountDefault vtpostClient.ConfigAccount `yaml:"account_default"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_VTPOST"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENV":              &c.Env,
		p + "_DEFAULT_USERNAME": &c.AccountDefault.Username,
		p + "_DEFAULT_PASSWORD": &c.AccountDefault.Password,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env: "test",
		AccountDefault: vtpostClient.ConfigAccount{
			Username: "tuan@eye-solution.vn",
			Password: "1234@5678",
		},
	}
}

type Connection struct {
	ClientCode string
}

type CalcShippingFeeAllServicesArgs struct {
	ArbitraryID  int64 // This is provided as a seed, for stable randomization
	FromProvince *mdlocation.Province
	FromDistrict *mdlocation.District
	ToProvince   *mdlocation.Province
	ToDistrict   *mdlocation.District

	Request *vtpostClient.CalcShippingFeeAllServicesRequest
	Result  []*model.AvailableShippingService
}

type GetShippingFeeLinesCommand struct {
	ServiceID    string // Required for detecting which client
	FromProvince *mdlocation.Province
	FromDistrict *mdlocation.District
	ToProvince   *mdlocation.Province
	ToDistrict   *mdlocation.District

	Request *vtpostClient.CalcShippingFeeRequest
	Result  *GetShippingFeeLineResponse
}

type GetShippingFeeLineResponse struct {
	ShippingFeeLines   []*model.ShippingFeeLine
	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
}

type CreateOrderArgs struct {
	ServiceID string // Required for detecting which client
	Request   *vtpostClient.CreateOrderRequest
	Result    *vtpostClient.CreateOrderResponse
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *vtpostClient.CancelOrderRequest
	Result    *vtpostClient.CommonResponse
}
