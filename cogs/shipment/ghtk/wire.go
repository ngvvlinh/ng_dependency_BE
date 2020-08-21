// +build wireinject

package _ghtk

import (
	"github.com/google/wire"

	"o.o/backend/pkg/integration/shipping/ghtk/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	NewGHTKWebhookServer,
)
