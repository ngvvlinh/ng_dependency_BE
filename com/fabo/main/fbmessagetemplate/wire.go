// +build wireinject

package fbmessagetemplate

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	FbMessagingQueryMessageBus,
	NewFbMessagingQuery,
	FbMessageTemplateAggregateMessageBus,
	NewFbMessageTemplateAggregate,
)
