// +build wireinject

package fbpage

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewFbPageUtil,
	NewFbPageAggregate, FbExternalPageAggregateMessageBus,
	NewFbPageQuery, FbPageQueryMessageBus,
	NewProcessManager,
)
