package handler

import (
	"context"

	"o.o/api/top/external/types"
	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/pgevent"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFulfillmentEvent", l.Object("pgevent", event))
	var history shipmodel.FulfillmentHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbFulfillmentHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.ID("fulfillment_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	var ffm shipmodel.Fulfillment
	if ok, err := h.db.Where("id = ?", id).Get(&ffm); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}

	convertpb.PbFulfillmentHistory(history)
	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Fulfillment: convertpb.PbFulfillment(&ffm),
	}
	change.Changed = &types.ChangeOneOf{
		Fulfillment: changed,
	}
	accountIDs := []dot.ID{ffm.ShopID, ffm.PartnerID}
	return h.sender.CollectPb(ctx, entity_type.Fulfillment, id, ffm.ShopID, accountIDs, change)
}
