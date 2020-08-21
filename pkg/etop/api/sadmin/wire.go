// +build wireinject

package sadmin

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewWebhookCallbackService,
	wire.Struct(new(MiscService), "*"),
	wire.Struct(new(UserService), "*"),
	wire.Struct(new(WebhookService), "*"),
	NewServers,
)
