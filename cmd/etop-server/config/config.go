package config

import (
	"errors"
	"strings"

	"o.o/api/main/invitation"
	_telebot "o.o/backend/cogs/base/telebot"
	config_server "o.o/backend/cogs/config/_server"
	database_all "o.o/backend/cogs/database/_all"
	shipment_all "o.o/backend/cogs/shipment/_all"
	_uploader "o.o/backend/cogs/uploader"
	"o.o/backend/com/main/invitation/aggregate"
	"o.o/backend/com/main/shipping/carrier"
	ecomconfig "o.o/backend/com/web/ecom/config"
	"o.o/backend/pkg/common/apifw/captcha"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/etop/api/export"
	"o.o/backend/pkg/etop/apix/partner"
	orderS "o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/etop/upload"
	"o.o/backend/pkg/integration/email"
	vtpayclient "o.o/backend/pkg/integration/payment/vtpay/client"
	ahamoveclient "o.o/backend/pkg/integration/shipnow/ahamove/client"
	ahamoveserver "o.o/backend/pkg/integration/shipnow/ahamove/server"
	"o.o/backend/pkg/integration/sms"
)

const (
	ChannelWebhook         = "webhook"
	ChannelImport          = "import"
	ChannelSMS             = "sms"
	ChannelDataWarehouse   = "etl"
	ChannelShipmentCarrier = "shipment_carrier"
)

// Config ...
type Config struct {
	SharedConfig config_server.SharedConfig `yaml:",inline"`
	Databases    database_all.Config        `yaml:",inline"`
	Shipment     shipment_all.Config        `yaml:",inline"`

	Redis cc.Redis `yaml:"redis"`

	Kafka       cc.Kafka         `yaml:"kafka"`
	Upload      upload.Config    `yaml:"upload"`
	Export      export.Config    `yaml:"export"`
	TelegramBot cc.TelegramBot   `yaml:"telegram_bot"`
	SMTP        email.SMTPConfig `yaml:"smtp"`
	Email       cc.EmailConfig   `yaml:"email"`
	SMS         sms.Config       `yaml:"sms"`
	Captcha     captcha.Config   `yaml:"captcha"`

	Ahamove        ahamoveclient.Config        `yaml:"ahamove"`
	AhamoveWebhook ahamoveserver.WebhookConfig `yaml:"ahamove_webhook"`
	Ecom           ecomconfig.Config           `yaml:"ecom"`

	VTPay vtpayclient.Config `yaml:"vtpay"`

	URL struct {
		Auth     partner.AuthURL `yaml:"auth"`
		MainSite string          `yaml:"main_site"`
	} `yaml:"url"`

	ThirdPartyHost string         `yaml:"third_party_host"`
	Secret         cc.SecretToken `yaml:"secret"`

	Invitation invitation.Config

	WhiteLabel cc.WhiteLabel `yaml:"white_label"`

	FlagEnableNewLinkInvitation           aggregate.FlagEnableNewLinkInvitation        `yaml:"flag_enable_new_link_invitation"`
	FlagApplyShipmentPrice                carrier.FlagApplyShipmentPrice               `yaml:"flag_apply_shipment_price"`
	FlagFaboOrderAutoConfirmPaymentStatus orderS.FlagFaboOrderAutoConfirmPaymentStatus `yaml:"flag_fabo_auto_confirm_payment_status"`
}

// Default ...
func Default() Config {
	cfg := Config{
		SharedConfig: config_server.DefaultConfig(),
		Databases:    database_all.DefaultConfig(),
		Redis:        cc.DefaultRedis(),
		Kafka: cc.Kafka{
			Enabled:     false,
			Brokers:     nil,
			TopicPrefix: "etop",
		},
		Upload: _uploader.DefaultConfig(),
		Export: export.Config{
			DirExport: "/tmp",
			URLPrefix: "http://localhost:8080",
		},
		TelegramBot:    _telebot.DefaultConfig(),
		Shipment:       shipment_all.DefaultConfig(),
		Ahamove:        ahamoveclient.DefaultConfig(),
		AhamoveWebhook: ahamoveserver.WebhookConfig{Port: 9052},
		Ecom: ecomconfig.Config{
			HTTP:     cc.HTTP{Port: 8100},
			MainSite: "http://localhost:8100",
		},
		VTPay: vtpayclient.DefaultConfig(),
		SMS: sms.Config{
			Mock:    true,
			Enabled: true,
		},
		Captcha: captcha.Config{
			Secret:        "6LcVOnkUAAAAALKlDJY_IYfQUmBfD_36azKtCv9P",
			LocalPasscode: "recaptcha_token",
		},
		Secret:         "secret",
		ThirdPartyHost: "https://etop.d.etop.vn",

		Invitation: invitation.Config{
			Secret: "IBVEhECSHtJiBoxQKOVafHW58zt9qRK7",
		},
	}
	cfg.Email = cc.EmailConfig{
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
	if err != nil {
		return cfg, err
	}

	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.Databases.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cfg.SMS.MustLoadEnv()
	cfg.SMTP.MustLoadEnv()
	cfg.Captcha.MustLoadEnv()
	cfg.Shipment.GHN.MustLoadEnv()
	cfg.Shipment.GHTK.MustLoadEnv()
	cfg.Shipment.VTPost.MustLoadEnv()
	cfg.Ahamove.MustLoadEnv()
	cfg.VTPay.MustLoadEnv()
	cc.MustLoadEnv("ET_SADMIN_TOKEN", &cfg.SharedConfig.SAdminToken)

	if cfg.ThirdPartyHost == "" && cfg.SharedConfig.Env != cmenv.EnvDev.String() {
		return cfg, errors.New("Empty third_party_host")
	}
	cfg.ThirdPartyHost = strings.TrimSuffix(cfg.ThirdPartyHost, "/")
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	return cfg, err
}
