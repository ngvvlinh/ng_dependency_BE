package vtpay

import (
	"github.com/google/wire"

	"o.o/backend/pkg/integration/payment/vtpay/client"
)

var WireSet = wire.NewSet(
	client.New,
	New,
)
