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

var keyMapTableEntities = map[string]string{
	"shop_product": "product",
	"shop_variant": "variant",
}

func (h *Handler) HandleShopProductEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopProductEvent", l.Object("pgevent", event))
	var history catalogmodel.ShopProductHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopProductHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninsteresting changes", l.ID("product_id", changed.Id))
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
	entity := GetEntityFromTableName(event.Table)
	return h.sender.CollectPb(ctx, entity, id, accountIDs, change)
}

func GetEntityFromTableName(tableName string) string {
	res := keyMapTableEntities[tableName]
	if res == "" {
		res = tableName
	}
	return res
}
