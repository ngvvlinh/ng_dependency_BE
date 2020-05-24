// +build wireinject

package fbpage

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewFbPageAggregate, FbExternalPageAggregateMessageBus,
	NewFbPageQuery, FbPageQueryMessageBus,
)
