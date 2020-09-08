// +build wireinject

package _vtpost

import (
	"github.com/google/wire"

	"o.o/backend/pkg/integration/shipping/vtpost/driver"
	"o.o/backend/pkg/integration/shipping/vtpost/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	NewVTPostWebhookServer,
	driver.NewShippingCodeGenerator,
)
