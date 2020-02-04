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
	return h.sender.CollectPb(ctx, "customer", id, customer.ShopID, accountIDs, change)
}
