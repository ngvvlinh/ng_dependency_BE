package sms_min

import (
	"github.com/google/wire"

	_sms "o.o/backend/cogs/sms"
	"o.o/backend/pkg/integration/sms"
	"o.o/backend/pkg/integration/sms/mock"
	"o.o/backend/pkg/integration/sms/telegram"
	"o.o/backend/pkg/integration/sms/vietguys"
)

var WireSet = wire.NewSet(
	_sms.WireSet,
	SupportedSMSDrivers,
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
