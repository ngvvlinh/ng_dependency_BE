package handler

import (
	"context"

	"o.o/api/top/external/types"
	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/pgevent"
	ordermodel "o.o/backend/com/main/ordering/model"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleOrderEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	var history ordermodel.OrderHistory
	if ok, err := h.db.Where("rid = ?", event.RID).Get(&history); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbOrderHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.ID("order_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	var order ordermodel.Order
	if ok, err := h.db.Where("id = ?", id).Get(&order); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("order not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}

	convertpb.PbOrderHistory(history)
	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Order: convertpb.PbOrder(&order),
	}
	change.Changed = &types.ChangeOneOf{
		Order: changed,
	}
	accountIDs := []dot.ID{order.ShopID, order.PartnerID}
	return h.sender.CollectPb(ctx, entity_type.Order, id, order.ShopID, accountIDs, change)
}
