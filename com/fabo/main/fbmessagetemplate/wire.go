// +build wireinject

package fbmessagetemplate

import (
	"github.com/google/wire"
	"o.o/backend/com/fabo/main/fbmessagetemplate/pm"
)

var WireSet = wire.NewSet(
	FbMessagingQueryMessageBus,
	NewFbMessagingQuery,
	FbMessageTemplateAggregateMessageBus,
	NewFbMessageTemplateAggregate,
	pm.NewProcessManager,
)
