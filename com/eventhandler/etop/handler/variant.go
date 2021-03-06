package handler

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/external/types"
	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/pgevent"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleShopVariantEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopVariantEvent", l.Object("pgevent", event))
	var history catalogmodel.ShopVariantHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("ShopVariant not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopVariantHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("variant_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.VariantID().ID().Apply(0)
	shopID := history.ShopID().ID().Apply(0)
	query := &catalog.GetShopVariantQuery{
		VariantID: id,
		ShopID:    shopID,
	}
	if err := h.catalogQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("ShopVariant not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	shopVariant := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Variant: convertpb.PbShopVariant(shopVariant),
	}
	change.Changed = &types.ChangeOneOf{
		Variant: changed,
	}
	accountIDs := []dot.ID{shopVariant.ShopID}
	return h.sender.CollectPb(ctx, entity_type.Variant, id, shopVariant.ShopID, accountIDs, change)
}
