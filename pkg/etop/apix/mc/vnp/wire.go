// +build wireinject

package vnp

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewServers,
	wire.Struct(new(VNPostService), "*"),
	wire.Struct(new(VNPostWebhookService), "*"),
)
