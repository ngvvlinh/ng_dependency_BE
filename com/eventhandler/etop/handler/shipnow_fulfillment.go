package handler

import (
	"context"

	"o.o/api/main/shipnow"
	typesx "o.o/api/top/external/types"
	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/pgevent"
	shipnowmodel "o.o/backend/com/main/shipnow/model"
	"o.o/backend/pkg/common/mq"
	convertx "o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleShipnowFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShipnowFulfillmentEvent", l.Object("pgevent", event))
	var history shipnowmodel.ShipnowFulfillmentHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Shipnow Fulfillment not found", l.Int64("rid", event.RID))
	}

	changed := convertx.Convert_shipnowmodel_ShipnowFulfillmentHistory_To_apix_ShipnowFulfillment(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("shipnow_fulfillment_id", changed.ID))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	query := &shipnow.GetShipnowFulfillmentQuery{
		ID: id,
	}
	if err := h.shipnowQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("Shipnow Fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	shipnowFfm := query.Result.ShipnowFulfillment
	change := pbChange(event)
	change.Latest = &typesx.LatestOneOf{
		ShipnowFulfillment: convertx.Convert_core_ShipnowFulfillment_To_apix_ShipnowFulfillment(shipnowFfm, nil),
	}
	change.Changed = &typesx.ChangeOneOf{
		ShipnowFulfillment: changed,
	}
	accountIDs := []dot.ID{shipnowFfm.ShopID, shipnowFfm.PartnerID}
	return h.sender.CollectPb(ctx, entity_type.ShipnowFulfillment, id, shipnowFfm.ShopID, accountIDs, change)
}
