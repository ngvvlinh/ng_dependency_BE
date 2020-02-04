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

func (h *Handler) HandleShopCustomerGroupCustomerEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleShopCustomerGroupCustomerEvent", l.Object("pgevent", event))
	var history customermodel.ShopCustomerGroupCustomerHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("CustomerGroupCustomer not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	changed := convertpb.PbShopCustomerGroupCustomerHistory(history)

	var haveError bool
	var index int
	customerID := history.CustomerID().ID().Apply(0)
	groupID := history.GroupID().ID().Apply(0)
	query := &customering.ListCustomerGroupsCustomersQuery{
		CustomerIDs: []dot.ID{customerID},
	}
	if err := h.customerQuery.Dispatch(ctx, query); err != nil {
		haveError = true
	}
	if query.Result == nil {
		haveError = true
	} else {
		haveError = true
		for i, relationship := range query.Result.CustomerGroupsCustomers {
			if relationship.GroupID == groupID {
				haveError = false
				index = i
				break
			}
		}
	}

	if haveError {
		ll.Warn("CustomerGroupCustomer not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	relationship := query.Result.CustomerGroupsCustomers[index]

	change := pbChange(event)
	change.Latest = &types.LatestOneOf{
		CustomerGroupRelationship: convertpb.PbCustomerGroupRelationship(relationship),
	}
	change.Changed = &types.ChangeOneOf{
		CustomerGroupRelationship: changed,
	}

	customerGroupQuery := &customering.GetCustomerGroupQuery{
		ID: groupID,
	}
	if err := h.customerQuery.Dispatch(ctx, customerGroupQuery); err != nil {
		ll.Warn("CustomerGroupCustomer not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	accountIDs := []dot.ID{customerGroupQuery.Result.ShopID}
	return h.sender.CollectPb(ctx, "customer_group_relationship", customerID+groupID, customerGroupQuery.Result.ShopID, accountIDs, change)
}
