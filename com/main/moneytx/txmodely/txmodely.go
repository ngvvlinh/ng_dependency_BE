package txmodely

import (
	"o.o/backend/com/main/moneytx/model"
	shippingmodely "o.o/backend/com/main/shipping/modely"
)

// +convert:type=moneytx.MoneyTransactionShippingExtended
type MoneyTransactionExtended struct {
	*model.MoneyTransactionShipping
	Fulfillments []*shippingmodely.FulfillmentExtended
}

type MoneyTransactionShippingEtopExtended struct {
	*model.MoneyTransactionShippingEtop
	MoneyTransactions []*MoneyTransactionExtended
}
