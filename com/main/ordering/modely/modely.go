package modely

import (
	"etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/sq"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenOrderExtended(
	&OrderExtended{}, &model.Order{}, sq.AS("o"),
	sq.LEFT_JOIN, &shipmodel.Fulfillment{}, sq.AS("f"), "o.id = f.order_id",
)

type OrderExtended struct {
	*model.Order
	Fulfillment *shipmodel.Fulfillment
}
