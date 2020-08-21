// +build wireinject

package sms_min

import (
	"github.com/google/wire"

	_sms "o.o/backend/cogs/sms"
)

var WireSet = wire.NewSet(
	_sms.WireSet,
	SupportedSMSDrivers,
)
