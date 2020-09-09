// +build wireinject

package _all

import (
	"github.com/google/wire"

	ghnwebhookv1 "o.o/backend/pkg/integration/shipping/ghn/webhook/v1"
	ghnwebhookv2 "o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
	njvwebhook "o.o/backend/pkg/integration/shipping/ninjavan/webhook"
)

var WireSet = wire.NewSet(
	ghnwebhookv1.WireSet,
	ghnwebhookv2.WireSet,
	njvwebhook.WireSet,
	NewShipmentWebhookServer,
)
