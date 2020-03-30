package modely

import (
	"etop.vn/backend/com/main/moneytx/model"
	"etop.vn/backend/com/main/shipping/modely"
)

type MoneyTransactionExtended struct {
	*model.MoneyTransactionShipping
	Fulfillments []*modely.FulfillmentExtended
}

type MoneyTransactionShippingEtopExtended struct {
	*model.MoneyTransactionShippingEtop
	MoneyTransactions []*MoneyTransactionExtended
}
