package vtpost

import (
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
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
		Env: cmenv.PartnerEnvTest,
		AccountDefault: vtpostclient.ConfigAccount{
			Username: "tuan@eye-solution.vn",
			Password: "1234@5678",
		},
	}
}
