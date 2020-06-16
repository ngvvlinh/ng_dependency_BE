package ghn

import (
	"o.o/api/main/location"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
)

type Account struct {
	Token     string `yaml:"token"`
	AccountID int    `yaml:"account_id"`
}

type Config struct {
	Env string `yaml:"env"`

	AccountDefault Account `yaml:"account_default"`
	AccountExt     Account `yaml:"account_ext"`
}

type WebhookConfig struct {
	cc.HTTP  `yaml:",inline"`
	Endpoint string `yaml:"endpoint"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_GHN"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENV":                &c.Env,
		p + "_DEFAULT_ACCOUNT_ID": &c.AccountDefault.AccountID,
		p + "_DEFAULT_TOKEN":      &c.AccountDefault.Token,
		p + "_EXT_ACCOUNT_ID":     &c.AccountExt.AccountID,
		p + "_EXT_TOKEN":          &c.AccountExt.Token,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env: cmenv.PartnerEnvTest,
		AccountDefault: Account{
			Token:     "5b20c7c194c06b03b2010913",
			AccountID: 503809,
		},
	}
}

func DefaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		HTTP:     cc.HTTP{Port: 9022},
		Endpoint: "http://callback-url",
	}
}

type RequestCreateOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.CreateOrderRequest
	Result    *ghnclient.CreateOrderResponse
}

type RequestFindAvailableServicesCommand struct {
	ServiceID    string // Required for detecting which client
	FromDistrict *location.District
	ToDistrict   *location.District
	Request      *ghnclient.FindAvailableServicesRequest
	// Result  *ghnClient.FindAvailableServicesResponse
	Result []*shippingsharemodel.AvailableShippingService
}

type RequestCalculateFeeCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.CalculateFeeRequest
	Result    *ghnclient.CalculateFeeResponse
}

type CalcShippingFeeCommand struct {
	FromDistrict *location.District
	ToDistrict   *location.District
	Request      *ghnclient.FindAvailableServicesRequest
	Result       []*ghnclient.AvailableService
}

type RequestGetOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.OrderCodeRequest
	Result    *ghnclient.Order
}

type RequestCancelOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.OrderCodeRequest
}

type RequestReturnOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.OrderCodeRequest
}

type RequestGetOrderLogsCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.OrderLogsRequest
	Result    *ghnclient.OrderLogsResponse
}
