package config

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/captcha"
	cc "etop.vn/backend/pkg/common/config"
	ahamoveclient "etop.vn/backend/pkg/integration/ahamove/client"
	"etop.vn/backend/pkg/integration/email"
	"etop.vn/backend/pkg/integration/ghn"
	"etop.vn/backend/pkg/integration/ghtk"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/backend/pkg/integration/vtpost"
)

const (
	ChannelWebhook              = "webhook"
	ChannelImport               = "import"
	PathAhamoveUserVerification = "/ahamove/user_verification"
)

type Upload struct {
	DirImportShopOrder   string `yaml:"dir_import_shop_order"`
	DirImportShopProduct string `yaml:"dir_import_shop_product"`
}

type Export struct {
	URLPrefix string `yaml:"url_prefix"`
	DirExport string `yaml:"dir_export"`
}

type EmailConfig struct {
	Enabled bool `yaml:"enabled"`

	ResetPasswordURL     string `valid:"url,required" yaml:"reset_password_url"`
	EmailVerificationURL string `valid:"url,required" yaml:"email_verification_url"`
}

// Config ...
type Config struct {
	Postgres         cc.Postgres      `yaml:"postgres"`
	PostgresLogs     cc.Postgres      `yaml:"postgres_logs"`
	PostgresNotifier cc.Postgres      `yaml:"postgres_notifier"`
	Redis            cc.Redis         `yaml:"redis"`
	HTTP             cc.HTTP          `yaml:"http"`
	Kafka            cc.Kafka         `yaml:"kafka"`
	Upload           Upload           `yaml:"upload"`
	Export           Export           `yaml:"export"`
	TelegramBot      cc.TelegramBot   `yaml:"telegram_bot"`
	SMTP             email.SMTPConfig `yaml:"smtp"`
	Email            EmailConfig      `yaml:"email"`
	SMS              sms.Config       `yaml:"sms"`
	Captcha          captcha.Config   `yaml:"captcha"`

	GHN            ghn.Config           `yaml:"ghn"`
	GHNWebhook     cc.HTTP              `yaml:"ghn_webhook"`
	GHTK           ghtk.Config          `yaml:"ghtk"`
	GHTKWebhook    cc.HTTP              `yaml:"ghtk_webhook"`
	VTPost         vtpost.Config        `yaml:"vtpost"`
	VTPostWebhook  cc.HTTP              `yaml:"vtpost_webhook"`
	Ahamove        ahamoveclient.Config `yaml:"ahamove"`
	AhamoveWebhook cc.HTTP              `yaml:"ahamove_webhook"`

	SAdminToken string `yaml:"sadmin_token"`
	ServeDoc    bool   `yaml:"serve_doc"`
	Env         string `yaml:"env"`

	URL struct {
		Auth     string `yaml:"auth"`
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Postgres:         cc.DefaultPostgres(),
		PostgresNotifier: cc.DefaultPostgres(),
		PostgresLogs:     cc.DefaultPostgresEtopLog(),
		Redis:            cc.DefaultRedis(),
		HTTP:             cc.HTTP{Port: 8080},
		Upload: Upload{
			DirImportShopOrder:   "/tmp",
			DirImportShopProduct: "/tmp",
		},
		Export: Export{
			DirExport: "/tmp",
			URLPrefix: "http://localhost:8080",
		},

		GHN:            ghn.DefaultConfig(),
		GHNWebhook:     cc.HTTP{Port: 9022},
		GHTK:           ghtk.DefaultConfig(),
		GHTKWebhook:    cc.HTTP{Port: 9032},
		VTPost:         vtpost.DefaultConfig(),
		VTPostWebhook:  cc.HTTP{Port: 9042},
		Ahamove:        ahamoveclient.DefaultConfig(),
		AhamoveWebhook: cc.HTTP{Port: 9052},

		SAdminToken: "PZJvDAY2.sadmin.HXnnEkdV",
		ServeDoc:    true,
		Captcha: captcha.Config{
			Secret:        "6LcVOnkUAAAAALKlDJY_IYfQUmBfD_36azKtCv9P",
			LocalPasscode: "recaptcha_token",
		},
		Env: cm.EnvDev,
	}
	cfg.Postgres.Database = "etop_dev"
	cfg.Email = EmailConfig{
		Enabled:              false,
		ResetPasswordURL:     "https://etop.d.etop.vn/reset-password",
		EmailVerificationURL: "https://etop.d.etop.vn/verify-email",
	}
	cfg.URL.Auth = "http://localhost:8080"
	cfg.URL.MainSite = "http://localhost:8080"
	return cfg
}

// DefaultTest returns default config for testing
func DefaultTest() Config {
	cfg := Default()
	cfg.Postgres.Database = "test"
	cfg.Postgres.Port = 5432
	return cfg
}

// Load loads config from file
func Load(isTest bool) (Config, error) {
	var cfg, defCfg Config
	if isTest {
		defCfg = DefaultTest()
	} else {
		defCfg = Default()
	}
	err := cc.LoadWithDefault(&cfg, defCfg)
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresLogs, "ET_POSTGRES_LOGS")
	cc.PostgresMustLoadEnv(&cfg.PostgresNotifier, "ET_POSTGRES_NOTIFIER")
	cfg.Redis.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.SMS.MustLoadEnv()
	cfg.SMTP.MustLoadEnv()
	cfg.Captcha.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
	cfg.Ahamove.MustLoadEnv()
	cc.MustLoadEnv("ET_SADMIN_TOKEN", &cfg.SAdminToken)
	return cfg, err
}
