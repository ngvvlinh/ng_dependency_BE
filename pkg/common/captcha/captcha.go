package captcha

import (
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/common/l"
	recaptcha "github.com/dpapathanasiou/go-recaptcha"
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
	if !cm.IsDev() && cfg.Secret == "" {
		ll.Fatal("Missing Captcha Secret Code")
	}
	if cm.IsProd() && cfg.LocalPasscode != "" {
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
	if !cm.IsProd() && gcfg.LocalPasscode != "" && gcfg.LocalPasscode == token {
		return nil
	}
	if ok, err := recaptcha.Confirm("", token); err != nil {
		return cm.Error(cm.Internal, "", err)
	} else if !ok {
		return cm.Error(cm.CaptchaInvalid, "", nil)
	}
	return nil
}
