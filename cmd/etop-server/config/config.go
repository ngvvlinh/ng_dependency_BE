package config

import (
	"errors"
	"io/ioutil"
	"strings"

	"etop.vn/api/main/invitation"
	crmsyncconfig "etop.vn/backend/cmd/supporting/crm-sync-service/config"
	"etop.vn/backend/com/supporting/crm/vtiger/mapping"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/captcha"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/integration/email"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
	vtpayclient "etop.vn/backend/pkg/integration/payment/vtpay/client"
	ahamoveclient "etop.vn/backend/pkg/integration/shipnow/ahamove/client"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	"etop.vn/backend/pkg/integration/shipping/vtpost"
	"etop.vn/backend/pkg/integration/sms"
	"etop.vn/common/jsonx"
)

const (
	ChannelWebhook = "webhook"
	ChannelImport  = "import"
	ChannelSMS     = "sms"

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
	PostgresCRM       cc.Postgres      `yaml:"postgres_crm"`
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
	GHNWebhook     cc.HTTP              `yaml:"ghn_webhook"`
	GHTK           ghtk.Config          `yaml:"ghtk"`
	GHTKWebhook    cc.HTTP              `yaml:"ghtk_webhook"`
	VTPost         vtpost.Config        `yaml:"vtpost"`
	VTPostWebhook  cc.HTTP              `yaml:"vtpost_webhook"`
	Ahamove        ahamoveclient.Config `yaml:"ahamove"`
	AhamoveWebhook cc.HTTP              `yaml:"ahamove_webhook"`

	Haravan haravanclient.Config `yaml:"haravan"`
	VTPay   vtpayclient.Config   `yaml:"vtpay"`

	SAdminToken string `yaml:"sadmin_token"`
	ServeDoc    bool   `yaml:"serve_doc"`
	Env         string `yaml:"env"`

	URL struct {
		Auth     string `yaml:"auth"`
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`

	ThirdPartyHost string `yaml:"third_party_host"`
	Secret         string `yaml:"secret"`

	Vtiger crmsyncconfig.Vtiger `yaml:"vtiger"`
	Vht    crmsyncconfig.Vht    `yaml:"vht"`

	Invitation invitation.Config
}

// Default ...
func Default() Config {
	cfg := Config{
		Postgres:          cc.DefaultPostgres(),
		PostgresNotifier:  cc.DefaultPostgres(),
		PostgresLogs:      cc.DefaultPostgres(),
		PostgresCRM:       cc.DefaultPostgres(),
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
		GHNWebhook:     cc.HTTP{Port: 9022},
		GHTK:           ghtk.DefaultConfig(),
		GHTKWebhook:    cc.HTTP{Port: 9032},
		VTPost:         vtpost.DefaultConfig(),
		VTPostWebhook:  cc.HTTP{Port: 9042},
		Ahamove:        ahamoveclient.DefaultConfig(),
		AhamoveWebhook: cc.HTTP{Port: 9052},
		Haravan:        haravanclient.DefaultConfig(),
		VTPay:          vtpayclient.DefaultConfig(),
		SMS: sms.Config{
			Mock: true,
		},
		SAdminToken: "PZJvDAY2.sadmin.HXnnEkdV",
		ServeDoc:    true,
		Captcha: captcha.Config{
			Secret:        "6LcVOnkUAAAAALKlDJY_IYfQUmBfD_36azKtCv9P",
			LocalPasscode: "recaptcha_token",
		},
		Env:            cm.EnvDev,
		Secret:         "secret",
		ThirdPartyHost: "https://etop.d.etop.vn",
		Vtiger:         crmsyncconfig.DefaultVtiger(),
		Vht:            crmsyncconfig.DefaultVht(),

		Invitation: invitation.Config{
			Secret: "IBVEhECSHtJiBoxQKOVafHW58zt9qRK7",
		},
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
	cc.PostgresMustLoadEnv(&cfg.PostgresCRM, "ET_POSTGRES_CRM")
	cc.PostgresMustLoadEnv(&cfg.PostgresAffiliate, "ET_POSTGRES_AFFILIATE")
	cfg.Redis.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.SMS.MustLoadEnv()
	cfg.SMTP.MustLoadEnv()
	cfg.Captcha.MustLoadEnv()
	cfg.Haravan.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
	cfg.Ahamove.MustLoadEnv()
	cfg.VTPay.MustLoadEnv()
	cfg.Vtiger.MustLoadEnv()
	cfg.Vht.MustLoadEnv()
	cc.MustLoadEnv("ET_SADMIN_TOKEN", &cfg.SAdminToken)

	if cfg.Haravan.Secret == "" {
		return cfg, errors.New("Empty Haravan secret")
	}
	if cfg.ThirdPartyHost == "" && !cm.IsDev() {
		return cfg, errors.New("Empty third_party_host")
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
