// +build wireinject

package api

import "github.com/google/wire"

var WireSet = wire.NewSet(
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(WebhookService), "*"),
	NewServers,
)
