package handler

import (
	"context"

	"etop.vn/api/shopping/customering"
	"etop.vn/api/top/external/types"
	"etop.vn/backend/com/handler/pgevent"
	customermodel "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

func (h *Handler) HandleShopCustomerGroupEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopCustomerGroupEvent", l.Object("pgevent", event))
	var history customermodel.ShopCustomerGroupHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("ShopCustomerGroup not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbCustomerGroupHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	query := &customering.GetCustomerGroupQuery{
		ID: id,
	}
	if err := h.customerQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("ShopCustomerGroup not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	customerGroup := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		CustomerGroup: convertpb.PbCustomerGroup(customerGroup),
	}
	change.Changed = &types.ChangeOneOf{
		CustomerGroup: changed,
	}
	accountIDs := []dot.ID{customerGroup.ShopID, customerGroup.PartnerID}
	return h.sender.CollectPb(ctx, "customer_group", id, customerGroup.ShopID, accountIDs, change)
}
