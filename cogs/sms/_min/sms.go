package sms_min

import (
	"o.o/backend/pkg/integration/sms"
	"o.o/backend/pkg/integration/sms/mock"
	"o.o/backend/pkg/integration/sms/telegram"
	"o.o/backend/pkg/integration/sms/vietguys"
)

func SupportedSMSDrivers(cfg sms.Config) []sms.DriverConfig {
	var mainDriver sms.Driver
	if cfg.Mock {
		mainDriver = mock.GetMock()
	} else if cfg.Telegram {
		mainDriver = telegram.GetTelegram()
	} else {
		mainDriver = vietguys.New(cfg.Vietguys)
	}
	return []sms.DriverConfig{
		{"", mainDriver},
	}
}
