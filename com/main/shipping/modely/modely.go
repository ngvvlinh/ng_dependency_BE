package modely

import (
	identitymodel "o.o/backend/com/main/identity/model"
	txmodel "o.o/backend/com/main/moneytx/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
)

// +sqlgen:     Fulfillment as f
// +sqlgen:join:      Shop  as s on s.id = f.shop_id
// +sqlgen:join:      Order as o on o.id = f.order_id
// +sqlgen:left-join: MoneyTransactionShipping as m on f.money_transaction_id = m.id
type FulfillmentExtended struct {
	*shipmodel.Fulfillment
	Shop                     *identitymodel.Shop
	Order                    *ordermodel.Order
	MoneyTransactionShipping *txmodel.MoneyTransactionShipping
}
