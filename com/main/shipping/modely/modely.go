package modely

import (
	identitymodel "etop.vn/backend/com/main/identity/model"
	txmodel "etop.vn/backend/com/main/moneytx/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
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
