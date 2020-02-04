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

func (h *Handler) HandleShopProductCollectionEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopProductCollectionEvent", l.Object("pgevent", event))
	var history catalogmodel.ShopCollectionHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("ShopProductCollection not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopProductCollectionHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("collection_id", changed.ID))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	shopID := history.ShopID().ID().Apply(0)
	query := &catalog.GetShopCollectionQuery{
		ID:     id,
		ShopID: shopID,
	}
	if err := h.catalogQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("ShopProductCollection not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	collection := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		ProductCollection: convertpb.PbShopProductCollection(collection),
	}
	change.Changed = &types.ChangeOneOf{
		ProductCollection: changed,
	}
	accountIDs := []dot.ID{collection.ShopID, collection.PartnerID}
	return h.sender.CollectPb(ctx, "product_collection", id, collection.ShopID, accountIDs, change)
}
