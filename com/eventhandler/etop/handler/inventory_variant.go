package handler

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/top/external/types"
	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/pgevent"
	inventorymodel "o.o/backend/com/main/inventory/model"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleInventoryVariantEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleInventoryVariantEvent", l.Object("pgevent", event))
	var history inventorymodel.InventoryVariantHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("InventoryVariant not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbInventoryVariantHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("variant_id", changed.VariantId))
		return mq.CodeOK, nil
	}

	id := history.VariantID().ID().Apply(0)
	shopID := history.ShopID().ID().Apply(0)
	query := &inventory.GetInventoryVariantQuery{
		ShopID:    id,
		VariantID: shopID,
	}
	if err := h.inventoryQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("InventoryVariant not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	inventoryVariant := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		InventoryLevel: convertpb.PbInventoryLevel(inventoryVariant),
	}
	change.Changed = &types.ChangeOneOf{
		InventoryLevel: changed,
	}
	accountIDs := []dot.ID{inventoryVariant.ShopID}
	return h.sender.CollectPb(ctx, entity_type.InventoryLevel, id, inventoryVariant.ShopID, accountIDs, change)
}
