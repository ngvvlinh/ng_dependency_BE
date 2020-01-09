package txmodely

import (
	"etop.vn/backend/com/main/moneytx/model"
	shippingmodely "etop.vn/backend/com/main/shipping/modely"
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
