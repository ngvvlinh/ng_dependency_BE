package ahamove

import (
	"time"

	"etop.vn/api/main/shipnow"

	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	ahamoveClient "etop.vn/backend/pkg/integration/ahamove/client"
)

type Config struct {
	Env string `yaml:"env"`

	AccountDefault ahamoveClient.AhamoveAccount `yaml:"account_default"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_AHAMOVE"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENV":                &c.Env,
		p + "_DEFAULT_ACCOUNT_ID": &c.AccountDefault.AccountID,
		p + "_DEFAULT_TOKEN":      &c.AccountDefault.Token,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env: "test",
		AccountDefault: ahamoveClient.AhamoveAccount{
			AccountID: "ahamove_default",
			Token:     "5cd40775eb44d75715c3b97a",
		},
	}
}

type CalcShippingFeeCommand struct {
	ArbitraryID int64 // This is provided as a seed, for stable randomization

	Request *ahamoveClient.CalcShippingFeeRequest
	Result  []*model.AvailableShippingService
}

type CalcSingleShippingFeeCommand struct {
	ServiceID string

	FromDistrictCode string
	ToDistrictCode   string

	Request *ahamoveClient.CalcShippingFeeRequest

	Result *model.AvailableShippingService
}

type CreateOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveClient.CreateOrderRequest
	Result  *ahamoveClient.CreateOrderResponse
}

type GetOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveClient.GetOrderRequest
	Result  *ahamoveClient.Order
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveClient.CancelOrderRequest
}

func ToShippingService(sfResp *ahamoveClient.CalcShippingFeeResponse, serviceID ServiceCode, providerServiceID string) *model.AvailableShippingService {
	if sfResp == nil {
		return nil
	}
	service := ServicesIndexID[serviceID]
	return &model.AvailableShippingService{
		Name:               service.Name,
		ServiceFee:         sfResp.TotalFee,
		ShippingFeeMain:    0,
		Provider:           shipnow.Ahamove,
		ProviderServiceID:  providerServiceID,
		ExpectedPickAt:     time.Now(),
		ExpectedDeliveryAt: time.Now(),
	}
}
