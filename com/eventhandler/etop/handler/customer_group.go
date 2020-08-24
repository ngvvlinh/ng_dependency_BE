package handler

import (
	"context"

	"o.o/api/shopping/customering"
	"o.o/api/top/external/types"
	"o.o/api/top/types/etc/entity_type"
	"o.o/backend/com/eventhandler/pgevent"
	customermodel "o.o/backend/com/shopping/customering/model"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
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
	return h.sender.CollectPb(ctx, entity_type.CustomerGroup, id, customerGroup.ShopID, accountIDs, change)
}
