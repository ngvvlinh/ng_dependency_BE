package sms_all

import (
	"github.com/google/wire"

	_sms "o.o/backend/cogs/sms"
	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/sms"
	imgroupsms "o.o/backend/pkg/integration/sms/imgroup"
	"o.o/backend/pkg/integration/sms/mock"
	"o.o/backend/pkg/integration/sms/vietguys"
	"o.o/common/l"
)

var ll = l.New()

var WireSet = wire.NewSet(
	_sms.WireSet,
	SupportedSMSDrivers,
)

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
	} else {
		mainDriver = vietguys.New(cfg.Vietguys)
	}

	return []sms.DriverConfig{
		{"", mainDriver},
		{drivers.ITopXKey, imgroupSMSClient},
	}
}
