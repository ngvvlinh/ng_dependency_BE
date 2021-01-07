// +build wireinject

package fabo

import (
	"github.com/google/wire"

	"o.o/backend/pkg/fabo/faboinfo"
)

var WireSet = wire.NewSet(
	wire.Struct(new(faboinfo.FaboPagesKit), "*"),
	wire.Struct(new(CustomerService), "*"),
	wire.Struct(new(CustomerConversationService), "*"),
	wire.Struct(new(PageService), "*"),
	wire.Struct(new(ShopService), "*"),
	wire.Struct(new(ExtraShipmentService), "*"),
	wire.Struct(new(SummaryService), "*"),
	wire.Struct(new(DemoService), "*"),
	NewServers,
)
