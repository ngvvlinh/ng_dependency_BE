package handler

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/external/types"
	"o.o/backend/com/handler/pgevent"
	catalogmodel "o.o/backend/com/main/catalog/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (h *Handler) HandleShopProductionCollectionRelationshipEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopProductionCollectionRelationshipEvent", l.Object("pgevent", event))
	var history catalogmodel.ProductShopCollectionHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, err
	} else if !ok {
		ll.Warn("ShopProductCollectionRelationship not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopProductionCollectionRelationshipHistory(history)

	deleted := false
	productID := history.ProductID().ID().Apply(0)
	collectionID := history.CollectionID().ID().Apply(0)
	shopID := history.ShopID().ID().Apply(0)
	query := &catalog.GetShopProductCollectionQuery{
		ProductID:    productID,
		CollectionID: collectionID,
	}
	err := h.catalogQuery.Dispatch(ctx, query)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		// no-op
	case cm.NotFound:
		deleted = true
	default:
		return mq.CodeStop, err
	}
	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		ProductCollectionRelationship: &types.ProductCollectionRelationship{
			ProductId:    productID,
			CollectionId: collectionID,
			Deleted:      deleted,
		},
	}
	change.Changed = &types.ChangeOneOf{
		ProductCollectionRelationship: changed,
	}

	accountIDs := []dot.ID{shopID}
	return h.sender.CollectPb(ctx, "product_collection_relationship", productID+collectionID, shopID, accountIDs, change)
}
