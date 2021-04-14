package kpay

import (
	"github.com/google/wire"
	"o.o/backend/pkg/integration/payment/kpay/client"
)

var WireSet = wire.NewSet(
	client.New,
	New,
)
