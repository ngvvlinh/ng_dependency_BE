package captcha

import (
	recaptcha "github.com/dpapathanasiou/go-recaptcha"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	Secret        string `yaml:"secret"`
	LocalPasscode string `yaml:"local_passcode"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := cc.EnvPrefix(prefix, "ET_CAPTCHA")
	cc.EnvMap{
		p + "_SECRET":         &c.Secret,
		p + "_LOCAL_PASSCODE": &c.LocalPasscode,
	}.MustLoad()
}

var gcfg Config

func Init(cfg Config) {
	gcfg = cfg
	if !cmenv.IsDev() && cfg.Secret == "" {
		ll.Fatal("Missing Captcha Secret Code")
	}
	if cmenv.IsProd() && cfg.LocalPasscode != "" {
		ll.Fatal("Do not use local passcode on production")
	}
	if cfg.Secret != "" {
		recaptcha.Init(cfg.Secret)
	}
}

func Verify(token string) error {
	if token == "" {
		return cm.Error(cm.CaptchaRequired, "", nil)
	}
	if !cmenv.IsProd() && gcfg.LocalPasscode != "" && gcfg.LocalPasscode == token {
		return nil
	}
	if ok, err := recaptcha.Confirm("", token); err != nil {
		return cm.Error(cm.Internal, "", err)
	} else if !ok {
		return cm.Error(cm.CaptchaInvalid, "", nil)
	}
	return nil
}
