package vtpost

import (
	"time"

	mdlocation "etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	vtpostclient "etop.vn/backend/pkg/integration/shipping/vtpost/client"
	"etop.vn/capi/dot"
)

type Config struct {
	Env string `yaml:"env"`

	AccountDefault vtpostclient.ConfigAccount `yaml:"account_default"`
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
		Env: cm.PartnerEnvTest,
		AccountDefault: vtpostclient.ConfigAccount{
			Username: "tuan@eye-solution.vn",
			Password: "1234@5678",
		},
	}
}

type Connection struct {
	ClientCode string
}

type CalcShippingFeeAllServicesArgs struct {
	ArbitraryID  dot.ID // This is provided as a seed, for stable randomization
	FromProvince *mdlocation.Province
	FromDistrict *mdlocation.District
	ToProvince   *mdlocation.Province
	ToDistrict   *mdlocation.District

	Request *vtpostclient.CalcShippingFeeAllServicesRequest
	Result  []*model.AvailableShippingService
}

type GetShippingFeeLinesCommand struct {
	ServiceID    string // Required for detecting which client
	FromProvince *mdlocation.Province
	FromDistrict *mdlocation.District
	ToProvince   *mdlocation.Province
	ToDistrict   *mdlocation.District

	Request *vtpostclient.CalcShippingFeeRequest
	Result  *GetShippingFeeLineResponse
}

type GetShippingFeeLineResponse struct {
	ShippingFeeLines   []*model.ShippingFeeLine
	ExpectedPickAt     time.Time
	ExpectedDeliveryAt time.Time
}

type CreateOrderArgs struct {
	ServiceID string // Required for detecting which client
	Request   *vtpostclient.CreateOrderRequest
	Result    *vtpostclient.CreateOrderResponse
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *vtpostclient.CancelOrderRequest
	Result    *vtpostclient.CommonResponse
}
