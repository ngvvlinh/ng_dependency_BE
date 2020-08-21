// +build wireinject

package _sms

import (
	"github.com/google/wire"

	"o.o/backend/com/etc/logging/smslog"
	"o.o/backend/pkg/integration/sms"
)

var WireSet = wire.NewSet(
	smslog.WireSet,
	sms.WireSet,
)
