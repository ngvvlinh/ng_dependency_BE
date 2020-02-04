package handler

import (
	"context"

	"etop.vn/api/shopping/addressing"
	"etop.vn/api/top/external/types"
	"etop.vn/backend/com/handler/pgevent"
	customermodel "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/mq"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

func (h *Handler) HandleShopTraderAddressEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopTraderAddressEvent", l.Object("pgevent", event))
	var history customermodel.ShopTraderAddressHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("CustomerAddress not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopTraderAddressHistory(ctx, history, h.locationQuery)
	if !changed.HasChanged() {
		ll.Debug("skip uninteresting changes", l.ID("id", changed.Id))
		return mq.CodeOK, nil
	}

	id := history.ID().ID().Apply(0)
	shopID := history.ShopID().ID().Apply(0)
	query := &addressing.GetAddressByIDQuery{
		ID:     id,
		ShopID: shopID,
	}
	if err := h.addressQuery.Dispatch(ctx, query); err != nil {
		ll.Warn("CustomerAddress not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	address := query.Result

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		CustomerAddress: convertpb.PbShopTraderAddress(ctx, address, h.locationQuery),
	}
	change.Changed = &types.ChangeOneOf{
		CustomerAddress: changed,
	}
	accountIDs := []dot.ID{address.ShopID, address.PartnerID}
	return h.sender.CollectPb(ctx, "customer_address", id, address.ShopID, accountIDs, change)
}
