// +build wireinject

package _fabo

import (
	"github.com/google/wire"
	"o.o/backend/pkg/etop/api/sadmin"
)

var WireSet = wire.NewSet(
	sadmin.NewWebhookCallbackService,
	wire.Struct(new(sadmin.WebhookService), "*"),
	NewServers,
)
