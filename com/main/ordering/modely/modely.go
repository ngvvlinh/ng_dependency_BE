package modely

import (
	"etop.vn/backend/com/main/ordering/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
)

// +sqlgen:           Order as order
// +sqlgen:left-join: Fulfillment as f on "order".id = f.order_id
type OrderExtended struct {
	*model.Order
	Fulfillment *shipmodel.Fulfillment
}
