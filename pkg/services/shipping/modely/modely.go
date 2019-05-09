package modely

import (
	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/etop/model"
	txmodel "etop.vn/backend/pkg/services/moneytx/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenFulfillmentExtended(
	&FulfillmentExtended{}, &model.Fulfillment{}, sq.AS("f"),
	sq.JOIN, &model.Shop{}, sq.AS("s"), "s.id = f.shop_id",
	sq.JOIN, &model.Order{}, sq.AS("o"), "o.id = f.order_id",
	sq.LEFT_JOIN, &txmodel.MoneyTransactionShipping{}, sq.AS("m"), "f.money_transaction_id = m.id",
)

type FulfillmentExtended struct {
	*model.Fulfillment
	Shop                     *model.Shop
	Order                    *model.Order
	MoneyTransactionShipping *txmodel.MoneyTransactionShipping
}
