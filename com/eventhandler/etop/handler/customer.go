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

func (h *Handler) HandleShopCustomerEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopCustomerEvent", l.Object("pgevent", event))
	var history customermodel.ShopCustomerHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Customer not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopCustomerHistory(history)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("customer_id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	shopID := history.ShopID().ID().Apply(0)
	query := &customering.GetCustomerByIDQuery{
		ID:     id,
		ShopID: shopID,
	}
	if err := h.customerQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("Customer not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	customer := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		Customer: convertpb.PbShopCustomer(customer),
	}
	change.Changed = &types.ChangeOneOf{
		Customer: changed,
	}
	accountIDs := []dot.ID{customer.ShopID, customer.PartnerID}
	return h.sender.CollectPb(ctx, entity_type.Customer, id, customer.ShopID, accountIDs, change)
}
