package ghtk

import (
	"strings"

	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/etop/model"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	"o.o/capi/dot"
)

type Config struct {
	Env string `yaml:"env"`

	AccountDefault    ghtkclient.GhtkAccount `yaml:"account_default"`
	AccountSamePrice  ghtkclient.GhtkAccount `yaml:"account_same_price"`
	AccountSamePrice2 ghtkclient.GhtkAccount `yaml:"account_same_price_2"`
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
		Env: cmenv.PartnerEnvTest,
		AccountDefault: ghtkclient.GhtkAccount{
			AccountID: "S1858017",
			Token:     "877F63B8596d32e1B0b2B0FcB0cF8E2980B28777",
		},
		AccountSamePrice: ghtkclient.GhtkAccount{
			AccountID: "S1858017",
			Token:     "877F63B8596d32e1B0b2B0FcB0cF8E2980B28777",
		},
		AccountSamePrice2: ghtkclient.GhtkAccount{
			AccountID: "S1858017",
			Token:     "877F63B8596d32e1B0b2B0FcB0cF8E2980B28777",
		},
	}
}

type Connection struct {
	ClientCode string
}

type CalcShippingFeeCommand struct {
	ArbitraryID dot.ID // This is provided as a seed, for stable randomization

	FromDistrictCode string
	ToDistrictCode   string

	Request *ghtkclient.CalcShippingFeeRequest
	Result  []*model.AvailableShippingService
}

type CalcSingleShippingFeeCommand struct {
	ServiceID string

	FromDistrictCode string
	ToDistrictCode   string

	Request *ghtkclient.CalcShippingFeeRequest

	Result *model.AvailableShippingService
}

type CreateOrderCommand struct {
	ServiceID string // Required for detecting which client
	Request   *ghtkclient.CreateOrderRequest
	Result    *ghtkclient.CreateOrderResponse
}

type GetOrderCommand struct {
	ServiceID string // Required for detecting which client
	LabelID   string // Mã đơn hàng của hệ thống GHTK
	PartnerID string // Mã đơn hàng thuộc hệ thống của đối tác

	Result *ghtkclient.GetOrderResponse
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client
	LabelID   string // Mã đơn hàng của hệ thống GHTK

	Result *ghtkclient.CommonResponse
}

// GHTK code format: "S1858017.SG5.19D.299241528". Normalize it to "299241528".
func NormalizeGHTKCode(code string) string {
	arrs := strings.Split(code, ".")
	return arrs[len(arrs)-1]
}
