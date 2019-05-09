package ghtk

import (
	"strings"

	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/etop/model"
	ghtkClient "etop.vn/backend/pkg/integration/ghtk/client"
)

type Config struct {
	Env string `yaml:"env"`

	AccountDefault    ghtkClient.GhtkAccount `yaml:"account_default"`
	AccountSamePrice  ghtkClient.GhtkAccount `yaml:"account_same_price"`
	AccountSamePrice2 ghtkClient.GhtkAccount `yaml:"account_same_price_2"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_GHTK"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENV":                     &c.Env,
		p + "_DEFAULT_ACCOUNT_ID":      &c.AccountDefault.AccountID,
		p + "_DEFAULT_TOKEN":           &c.AccountDefault.Token,
		p + "_SAME_PRICE_ACCOUNT_ID":   &c.AccountSamePrice.AccountID,
		p + "_SAME_PRICE_TOKEN":        &c.AccountSamePrice.Token,
		p + "_SAME_PRICE_2_ACCOUNT_ID": &c.AccountSamePrice2.AccountID,
		p + "_SAME_PRICE_2_TOKEN":      &c.AccountSamePrice2.Token,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env: "test",
		AccountDefault: ghtkClient.GhtkAccount{
			AccountID: "S1858017",
			Token:     "877F63B8596d32e1B0b2B0FcB0cF8E2980B28777",
		},
		AccountSamePrice: ghtkClient.GhtkAccount{
			AccountID: "S1858017",
			Token:     "877F63B8596d32e1B0b2B0FcB0cF8E2980B28777",
		},
		AccountSamePrice2: ghtkClient.GhtkAccount{
			AccountID: "S1858017",
			Token:     "877F63B8596d32e1B0b2B0FcB0cF8E2980B28777",
		},
	}
}

type Connection struct {
	ClientCode string
}

type CalcShippingFeeCommand struct {
	ArbitraryID int64 // This is provided as a seed, for stable randomization

	FromDistrictCode string
	ToDistrictCode   string

	Request *ghtkClient.CalcShippingFeeRequest
	Result  []*model.AvailableShippingService
}

type CalcSingleShippingFeeCommand struct {
	ServiceID string

	FromDistrictCode string
	ToDistrictCode   string

	Request *ghtkClient.CalcShippingFeeRequest

	Result *model.AvailableShippingService
}

type CreateOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghtkClient.CreateOrderRequest
	Result    *ghtkClient.CreateOrderResponse
}

type GetOrderCommand struct {
	ServiceID string // Required for detecting which client
	LabelID   string // Mã đơn hàng của hệ thống GHTK
	PartnerID string // Mã đơn hàng thuộc hệ thống của đối tác

	Result *ghtkClient.GetOrderResponse
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client
	LabelID   string // Mã đơn hàng của hệ thống GHTK

	Result *ghtkClient.CommonResponse
}

// GHTK code format: "S1858017.SG5.19D.299241528". Normalize it to "299241528".
func NormalizeGHTKCode(code string) string {
	arrs := strings.Split(code, ".")
	return arrs[len(arrs)-1]
}
