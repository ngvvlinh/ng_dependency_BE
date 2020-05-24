// +build wireinject

package address

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewQueryService, QueryServiceMessageBus,
)
