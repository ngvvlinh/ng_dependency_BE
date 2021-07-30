package sms_all

import (
	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/sms"
	imgroupsms "o.o/backend/pkg/integration/sms/imgroup"
	"o.o/backend/pkg/integration/sms/incom"
	"o.o/backend/pkg/integration/sms/mock"
	telegramsms "o.o/backend/pkg/integration/sms/telegram"
	"o.o/backend/pkg/integration/sms/vietguys"
	"o.o/common/l"
)

var ll = l.New()

func SupportedSMSDrivers(wlCfg cc.WhiteLabel, cfg sms.Config) []sms.DriverConfig {
	var imgroupSMSClient *imgroupsms.Client
	if wlCfg.IMGroup.SMS.APIKey != "" {
		imgroupSMSClient = imgroupsms.New(wlCfg.IMGroup.SMS)
	} else if !cmenv.IsDev() {
		ll.Panic("no sms config for whitelabel/imgroup")
	}
	var mainDriver sms.Driver
	if cfg.Mock {
		mainDriver = mock.GetMock()
	} else if cfg.Telegram {
		mainDriver = telegramsms.GetTelegram()
	} else if &cfg.Incom != nil {
		mainDriver = incom.New(cfg.Incom)
	} else {
		mainDriver = vietguys.New(cfg.Vietguys)
	}

	return []sms.DriverConfig{
		{"", mainDriver},
		{drivers.ITopXKey, imgroupSMSClient},
	}
}
