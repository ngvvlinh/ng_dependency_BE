package config

import (
	"errors"
	"io/ioutil"
	"strings"

	"etop.vn/api/main/invitation"
	"etop.vn/backend/com/supporting/crm/vtiger/mapping"
	ecomconfig "etop.vn/backend/com/web/ecom/config"
	"etop.vn/backend/pkg/common/apifw/captcha"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/integration/email"
	vtpayclient "etop.vn/backend/pkg/integration/payment/vtpay/client"
	ahamoveclient "etop.vn/backend/pkg/integration/shipnow/ahamove/client"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	"etop.vn/backend/pkg/integration/shipping/vtpost"
	"etop.vn/backend/pkg/integration/sms"
	imgroupsms "etop.vn/backend/pkg/integration/sms/imgroup"
	"etop.vn/common/jsonx"
)

const (
	ChannelWebhook       = "webhook"
	ChannelImport        = "import"
	ChannelSMS           = "sms"
	ChannelDataWarehouse = "etl"

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
	Postgres          cc.Postgres      `yaml:"postgres"`
	PostgresLogs      cc.Postgres      `yaml:"postgres_logs"`
	PostgresNotifier  cc.Postgres      `yaml:"postgres_notifier"`
	PostgresAffiliate cc.Postgres      `yaml:"postgres_affiliate"`
	Redis             cc.Redis         `yaml:"redis"`
	HTTP              cc.HTTP          `yaml:"http"`
	Kafka             cc.Kafka         `yaml:"kafka"`
	Upload            Upload           `yaml:"upload"`
	Export            Export           `yaml:"export"`
	TelegramBot       cc.TelegramBot   `yaml:"telegram_bot"`
	SMTP              email.SMTPConfig `yaml:"smtp"`
	Email             EmailConfig      `yaml:"email"`
	SMS               sms.Config       `yaml:"sms"`
	Captcha           captcha.Config   `yaml:"captcha"`

	GHN            ghn.Config           `yaml:"ghn"`
	GHNWebhook     ghn.WebhookConfig    `yaml:"ghn_webhook"`
	GHTK           ghtk.Config          `yaml:"ghtk"`
	GHTKWebhook    cc.HTTP              `yaml:"ghtk_webhook"`
	VTPost         vtpost.Config        `yaml:"vtpost"`
	VTPostWebhook  cc.HTTP              `yaml:"vtpost_webhook"`
	Ahamove        ahamoveclient.Config `yaml:"ahamove"`
	AhamoveWebhook cc.HTTP              `yaml:"ahamove_webhook"`
	Ecom           ecomconfig.Config    `yaml:"ecom"`

	VTPay vtpayclient.Config `yaml:"vtpay"`

	SAdminToken string `yaml:"sadmin_token"`
	ServeDoc    bool   `yaml:"serve_doc"`
	Env         string `yaml:"env"`

	URL struct {
		Auth     string `yaml:"auth"`
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`

	ThirdPartyHost string `yaml:"third_party_host"`
	Secret         string `yaml:"secret"`

	Invitation invitation.Config

	WhiteLabel struct {
		IMGroup struct {
			SMS imgroupsms.Config `yaml:"sms"`
		} `yaml:"imgroup"`
	} `yaml:"white_label"`

	FlagEnablePermission string `yaml:"flag_enable_permission"`
}

// Default ...
func Default() Config {
	cfg := Config{
		Postgres:          cc.DefaultPostgres(),
		PostgresNotifier:  cc.DefaultPostgres(),
		PostgresLogs:      cc.DefaultPostgres(),
		PostgresAffiliate: cc.DefaultPostgres(),
		Redis:             cc.DefaultRedis(),
		HTTP:              cc.HTTP{Port: 8080},
		Kafka: cc.Kafka{
			Enabled:     false,
			Brokers:     nil,
			TopicPrefix: "etop",
		},
		Upload: Upload{
			DirImportShopOrder:   "/tmp",
			DirImportShopProduct: "/tmp",
		},
		Export: Export{
			DirExport: "/tmp",
			URLPrefix: "http://localhost:8080",
		},

		GHN:            ghn.DefaultConfig(),
		GHNWebhook:     ghn.DefaultWebhookConfig(),
		GHTK:           ghtk.DefaultConfig(),
		GHTKWebhook:    cc.HTTP{Port: 9032},
		VTPost:         vtpost.DefaultConfig(),
		VTPostWebhook:  cc.HTTP{Port: 9042},
		Ahamove:        ahamoveclient.DefaultConfig(),
		AhamoveWebhook: cc.HTTP{Port: 9052},
		Ecom: ecomconfig.Config{
			HTTP:     cc.HTTP{Port: 8100},
			MainSite: "http://localhost:8100",
		},
		VTPay: vtpayclient.DefaultConfig(),
		SMS: sms.Config{
			Mock:    true,
			Enabled: true,
		},
		SAdminToken: "PZJvDAY2.sadmin.HXnnEkdV",
		ServeDoc:    true,
		Captcha: captcha.Config{
			Secret:        "6LcVOnkUAAAAALKlDJY_IYfQUmBfD_36azKtCv9P",
			LocalPasscode: "recaptcha_token",
		},
		Env:            cmenv.EnvDev.String(),
		Secret:         "secret",
		ThirdPartyHost: "https://etop.d.etop.vn",

		Invitation: invitation.Config{
			Secret: "IBVEhECSHtJiBoxQKOVafHW58zt9qRK7",
		},
		FlagEnablePermission: "all",
	}
	cfg.Postgres.Database = "etop_dev"
	cfg.PostgresAffiliate.Database = "etop_dev"
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
	cc.PostgresMustLoadEnv(&cfg.PostgresAffiliate, "ET_POSTGRES_AFFILIATE")
	cfg.Redis.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cfg.SMS.MustLoadEnv()
	cfg.SMTP.MustLoadEnv()
	cfg.Captcha.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
	cfg.Ahamove.MustLoadEnv()
	cfg.VTPay.MustLoadEnv()
	cc.MustLoadEnv("ET_SADMIN_TOKEN", &cfg.SAdminToken)

	if cfg.ThirdPartyHost == "" && !cmenv.IsDev() {
		return cfg, errors.New("Empty third_party_host")
	}
	if cfg.GHNWebhook.Endpoint == "" {
		return cfg, errors.New("Empty GHN webhook endpoint")
	}
	cfg.ThirdPartyHost = strings.TrimSuffix(cfg.ThirdPartyHost, "/")
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	return cfg, err
}

// ReadMappingFile read mapping json file for mapping fields between vtiger and etop
func ReadMappingFile(filename string) (configMap mapping.ConfigMap, _ error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = jsonx.Unmarshal(body, &configMap)
	return
}
