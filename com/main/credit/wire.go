package credit

import (
	"github.com/google/wire"
	creditpm "o.o/backend/com/main/credit/pm"
)

var WireSet = wire.NewSet(
	creditpm.New,
	NewAggregateCredit, CreditAggregateMessageBus,
	NewQueryCredit, CreditQueryServiceMessageBus,
)
