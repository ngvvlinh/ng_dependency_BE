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

func (h *Handler) HandleShopProductEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopProductEvent", l.Object("pgevent", event))
	var history catalogmodel.ShopProductHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("ShopProduct not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopProductHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("product_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ProductID().ID().Apply(0)
	query := &catalog.GetShopProductQuery{
		ProductID: id,
	}
	if err := h.catalogQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("ShopProduct not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	shopProduct := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Product: convertpb.PbShopProduct(shopProduct),
	}
	change.Changed = &types.ChangeOneOf{
		Product: changed,
	}
	accountIDs := []dot.ID{shopProduct.ShopID, shopProduct.PartnerID}
	return h.sender.CollectPb(ctx, entity_type.Product, id, shopProduct.ShopID, accountIDs, change)
}
