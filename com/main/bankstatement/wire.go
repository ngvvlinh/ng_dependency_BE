package bankstatement

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAggregateBankStatement, BankStatementAggregateMessageBus,
	NewQueryBankStatement, BankStatementQueryServiceMessageBus,
)
