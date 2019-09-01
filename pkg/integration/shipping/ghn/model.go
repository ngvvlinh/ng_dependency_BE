package ghn

import (
	"etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
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
		Env: cm.PartnerEnvTest,
		AccountDefault: Account{
			Token:     "5b20c7c194c06b03b2010913",
			AccountID: 503809,
		},
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
	Result []*model.AvailableShippingService
}

type RequestCalculateFeeCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghnclient.CalculateFeeRequest
	Result    *ghnclient.CalculateFeeResponse
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