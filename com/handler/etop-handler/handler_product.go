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
	return h.sender.CollectPb(ctx, "product", id, shopProduct.ShopID, accountIDs, change)
}
