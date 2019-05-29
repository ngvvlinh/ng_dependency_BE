package ghn

import (
	"etop.vn/api/main/location"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	ghnClient "etop.vn/backend/pkg/integration/ghn/client"
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
		Env: "test",
		AccountDefault: Account{
			Token:     "5b20c7c194c06b03b2010913",
			AccountID: 503809,
		},
	}
}

type RequestCreateOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnClient.CreateOrderRequest
	Result    *ghnClient.CreateOrderResponse
}

type RequestFindAvailableServicesCommand struct {
	ServiceID    string // Required for detecting which client
	FromDistrict *location.District
	ToDistrict   *location.District
	Request      *ghnClient.FindAvailableServicesRequest
	// Result  *ghnClient.FindAvailableServicesResponse
	Result []*model.AvailableShippingService
}

type RequestCalculateFeeCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnClient.CalculateFeeRequest
	Result    *ghnClient.CalculateFeeResponse
}

type RequestGetOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnClient.OrderCodeRequest
	Result    *ghnClient.Order
}

type RequestCancelOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnClient.OrderCodeRequest
}

type RequestReturnOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnClient.OrderCodeRequest
}

type RequestGetOrderLogsCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnClient.OrderLogsRequest
	Result    *ghnClient.OrderLogsResponse
}