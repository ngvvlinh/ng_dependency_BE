package kpay

import (
	"github.com/google/wire"
	"o.o/backend/com/external/payment/kpay/gateway/server"
)

var WireSet = wire.NewSet(
	server.New,
)
