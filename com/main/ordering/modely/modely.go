package modely

import (
	"etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/sql/sq"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenOrderExtended(
	&OrderExtended{}, &model.Order{}, `"order"`,
	sq.LEFT_JOIN, &shipmodel.Fulfillment{}, "f", `"order".id = f.order_id`,
)

type OrderExtended struct {
	*model.Order
	Fulfillment *shipmodel.Fulfillment
}
