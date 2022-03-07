// +build wireinject

package _ntx

import (
	"github.com/google/wire"

	"o.o/backend/pkg/integration/shipping/ntx/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	NewNTXWebhookServer,
)
