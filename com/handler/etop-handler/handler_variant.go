package handler

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/external/types"
	"etop.vn/backend/com/handler/pgevent"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
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
	return h.sender.CollectPb(ctx, "variant", id, shopVariant.ShopID, accountIDs, change)
}
