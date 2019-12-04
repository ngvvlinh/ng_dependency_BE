package modely

import (
	txmodel "etop.vn/backend/com/main/moneytx/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFulfillmentExtended(
	&FulfillmentExtended{}, &shipmodel.Fulfillment{}, "f",
	sq.JOIN, &model.Shop{}, "s", "s.id = f.shop_id",
	sq.JOIN, &ordermodel.Order{}, "o", "o.id = f.order_id",
	sq.LEFT_JOIN, &txmodel.MoneyTransactionShipping{}, "m", "f.money_transaction_id = m.id",
)

type FulfillmentExtended struct {
	*shipmodel.Fulfillment
	Shop                     *model.Shop
	Order                    *ordermodel.Order
	MoneyTransactionShipping *txmodel.MoneyTransactionShipping
}
