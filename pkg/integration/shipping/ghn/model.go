package ghn

import (
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
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
