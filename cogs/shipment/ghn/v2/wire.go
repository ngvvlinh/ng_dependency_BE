// +build wireinject

package v2

import (
	"github.com/google/wire"

	ghnwebhook "o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
)

var WireSet = wire.NewSet(
	ghnwebhook.WireSet,
	NewGHNWebhookServer,
)
