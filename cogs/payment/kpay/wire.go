// +build wireinject

package kpay

import (
	"github.com/google/wire"
	"o.o/backend/com/external/payment/kpay"
	paymentkpay "o.o/backend/pkg/integration/payment/kpay"
)

var WireSet = wire.NewSet(
	kpay.WireSet,
	paymentkpay.WireSet,
)
