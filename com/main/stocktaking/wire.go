// +build wireinject

package stocktaking

import (
	"github.com/google/wire"
	"o.o/backend/com/main/stocktaking/aggregate"
	"o.o/backend/com/main/stocktaking/query"
)

var WireSet = wire.NewSet(
	aggregate.NewAggregateStocktake, aggregate.StocktakeAggregateMessageBus,
	query.NewQueryStocktake, query.StocktakeQueryMessageBus,
)
