package credit

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAggregateCredit, CreditAggregateMessageBus,
	NewQueryCredit, CreditQueryServiceMessageBus,
)
