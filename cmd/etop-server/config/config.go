package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
)

const (
	ChannelWebhook = "webhook"
	ChannelImport  = "import"
	ChannelSMS     = "sms"

	PathAhamoveUserVerification = "/ahamove/user_verification"
)

var exampleMappingFile = filepath.Join(
	os.Getenv("ETOPDIR"),
	"backend/cmd/etop-server/config/field_mapping_example.json",
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

type Vtiger struct {
	ServiceURL string `yaml:"service_url"`
	Username   string `yaml:"username"`
	APIKey     string `yaml:"api_key"`
}

type Vht struct {
	ServiceURL string `yaml:"service_url"`
	UserName   string `yaml:"user_name"`
	PassWord   string `yaml:"pass_word"`
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

	Vtiger      Vtiger `yaml:"vtiger"`
	MappingFile string `yaml:"mapping_file"`
	Vht         Vht    `yaml:vht`
}

// Default ...
func Default() Config {
	cfg := Config{
		Postgres:         cc.DefaultPostgres(),
		PostgresNotifier: cc.DefaultPostgres(),
		PostgresLogs:     cc.DefaultPostgres(),
		Redis:            cc.DefaultRedis(),
		HTTP:             cc.HTTP{Port: 8080},
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

		SAdminToken: "PZJvDAY2.sadmin.HXnnEkdV",
		ServeDoc:    true,
		Captcha: captcha.Config{
			Secret:        "6LcVOnkUAAAAALKlDJY_IYfQUmBfD_36azKtCv9P",
			LocalPasscode: "recaptcha_token",
		},
		Env:            cm.EnvDev,
		Secret:         "secret",
		ThirdPartyHost: "https://etop.d.etop.vn",
		Vtiger: Vtiger{
			ServiceURL: "http://vtiger/webservice.php",
			Username:   "admin",
			APIKey:     "q5dZOnJYGlmPY2nc",
		},
		Vht: Vht{
			UserName: "5635810cde4c14ebf6a41341f4e68395",
			PassWord: "36828473ce0d87db8cc29798f6b8aa1e",
		},
		MappingFile: exampleMappingFile,
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
	cfg.Haravan.MustLoadEnv()
	cfg.GHN.MustLoadEnv()
	cfg.GHTK.MustLoadEnv()
	cfg.VTPost.MustLoadEnv()
	cfg.Ahamove.MustLoadEnv()
	cfg.VTPay.MustLoadEnv()
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
	err = json.Unmarshal(body, &configMap)
	return
}
