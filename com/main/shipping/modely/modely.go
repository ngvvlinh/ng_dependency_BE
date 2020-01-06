package modely

import (
	identitymodel "etop.vn/backend/com/main/identity/model"
	txmodel "etop.vn/backend/com/main/moneytx/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/sql/sq"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFulfillmentExtended(
	&FulfillmentExtended{}, &shipmodel.Fulfillment{}, "f",
	sq.JOIN, &identitymodel.Shop{}, "s", "s.id = f.shop_id",
	sq.JOIN, &ordermodel.Order{}, "o", "o.id = f.order_id",
	sq.LEFT_JOIN, &txmodel.MoneyTransactionShipping{}, "m", "f.money_transaction_id = m.id",
)

type FulfillmentExtended struct {
	*shipmodel.Fulfillment
	Shop                     *identitymodel.Shop
	Order                    *ordermodel.Order
	MoneyTransactionShipping *txmodel.MoneyTransactionShipping
}
