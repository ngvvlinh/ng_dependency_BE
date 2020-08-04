package config

import (
	"o.o/api/main/invitation"
	_telebot "o.o/backend/cogs/base/telebot"
	config_server "o.o/backend/cogs/config/_server"
	database_min "o.o/backend/cogs/database/_min"
	shipment_fabo "o.o/backend/cogs/shipment/_fabo"
	storage_all "o.o/backend/cogs/storage/_all"
	_uploader "o.o/backend/cogs/uploader"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/main/invitation/aggregate"
	"o.o/backend/pkg/common/apifw/captcha"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/storage"
	"o.o/backend/pkg/etop/api/export"
	orderS "o.o/backend/pkg/etop/logic/orders"
	"o.o/backend/pkg/integration/email"
	"o.o/backend/pkg/integration/sms"
)

type Config struct {
	SharedConfig config_server.SharedConfig `yaml:",inline"`
	Databases    database_min.Config        `yaml:",inline"`
	Shipment     shipment_fabo.Config       `yaml:",inline"`

	Redis cc.Redis `yaml:"redis"`

	Kafka cc.Kafka `yaml:"kafka"`

	UploadDirs    storage.DirConfigs       `yaml:"upload_dirs"`
	ExportDirs    export.ConfigDirs        `yaml:"export_dirs"`
	StorageDriver storage_all.DriverConfig `yaml:"storage_driver"`

	TelegramBot cc.TelegramBot   `yaml:"telegram_bot"`
	SMTP        email.SMTPConfig `yaml:"smtp"`
	Email       cc.EmailConfig   `yaml:"email"`
	SMS         sms.Config       `yaml:"sms"`
	Captcha     captcha.Config   `yaml:"captcha"`

	URL struct {
		MainSite string `yaml:"main_site"`
	} `yaml:"url"`

	Secret cc.SecretToken `yaml:"secret"`

	Invitation invitation.Config

	WhiteLabel cc.WhiteLabel `yaml:"white_label"`

	FlagEnableNewLinkInvitation           aggregate.FlagEnableNewLinkInvitation        `yaml:"flag_enable_new_link_invitation"`
	FlagFaboOrderAutoConfirmPaymentStatus orderS.FlagFaboOrderAutoConfirmPaymentStatus `yaml:"flag_fabo_auto_confirm_payment_status"`

	FacebookApp fbclient.AppConfig `yaml:"facebook_app"`
	Webhook     WebhookConfig      `yaml:"webhook"`
}

type WebhookConfig struct {
	HTTP        cc.HTTP `yaml:"http"`
	VerifyToken string  `yaml:"verify_token"`
}

func Default() Config {
	cfg := Config{
		SharedConfig:  config_server.DefaultConfig(),
		Databases:     database_min.DefaultConfig(),
		Redis:         cc.DefaultRedis(),
		Kafka:         cc.DefaultKafka(),
		UploadDirs:    _uploader.DefaultConfig(),
		StorageDriver: storage_all.DefaultDriver(),
		ExportDirs: export.ConfigDirs{
			Export: storage.DirConfig{
				Path:      "export/dl",
				URLPath:   "/export/dl",
				URLPrefix: "http://localhost:8080/export/dl",
			},
		},
		TelegramBot: _telebot.DefaultConfig(),
		Shipment:    shipment_fabo.DefaultConfig(),
		SMS: sms.Config{
			Mock:    true,
			Enabled: true,
		},
		Captcha: captcha.Config{
			Secret:        "6LcVOnkUAAAAALKlDJY_IYfQUmBfD_36azKtCv9P",
			LocalPasscode: "recaptcha_token",
		},
		Secret: "secret",

		Invitation: invitation.Config{
			Secret: "IBVEhECSHtJiBoxQKOVafHW58zt9qRK7",
		},
	}
	cfg.URL.MainSite = "http://localhost:8080"
	cfg.Email = cc.EmailConfig{
		Enabled:              false,
		ResetPasswordURL:     "https://etop.d.etop.vn/reset-password",
		EmailVerificationURL: "https://etop.d.etop.vn/verify-email",
	}
	cfg.FacebookApp = fbclient.AppConfig{
		ID:          "1581362285363031",
		Secret:      "b3962ddf033b295c2bd0b543fff904f7",
		AccessToken: "1581362285363031|eLuNU9-1KNA0AMNucV9PQIHCF1A",
	}
	cfg.Webhook = WebhookConfig{
		HTTP: cc.HTTP{
			Host: "",
			Port: 8081,
		},
	}
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	if err != nil {
		return cfg, err
	}

	cc.RedisMustLoadEnv(&cfg.Redis)
	cfg.Databases.MustLoadEnv()
	cfg.TelegramBot.MustLoadEnv()
	cfg.SMS.MustLoadEnv()
	cfg.SMTP.MustLoadEnv()
	cfg.Captcha.MustLoadEnv()

	// TODO(vu): remove this
	cfg.Shipment.GHN.MustLoadEnv()
	cfg.Shipment.GHTK.MustLoadEnv()
	cfg.Shipment.VTPost.MustLoadEnv()

	cc.MustLoadEnv("ET_SADMIN_TOKEN", &cfg.SharedConfig.SAdminToken)
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	cfg.FacebookApp.MustLoadEnv()
	return cfg, err
}
