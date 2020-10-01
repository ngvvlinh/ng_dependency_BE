// +build wireinject

package shippingcode

import (
	"github.com/google/wire"
	"o.o/backend/com/main/shippingcode/query"
)

var WireSet = wire.NewSet(
	query.QueryServiceMessageBus, query.NewQueryService,
)
