package shop

import (
	"context"

	"o.o/api/top/int/etop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/model"
)

type HistoryService struct{}

func (s *HistoryService) Clone() *HistoryService { res := *s; return &res }

func (s *HistoryService) GetFulfillmentHistory(ctx context.Context, r *GetFulfillmentHistoryEndpoint) error {

	filters := map[string]interface{}{
		"shop_id": r.Context.Shop.ID,
	}
	count := 0
	if r.All {
		count++
	}
	if r.OrderId != 0 {
		count++
		filters["order_id"] = r.OrderId
	}
	if r.Id != 0 {
		count++
		filters["id"] = r.Id
	}
	if count != 1 {
		return cm.Error(cm.InvalidArgument, "Must provide either all, id or order_id", nil)
	}

	paging := cmapi.CMPaging(r.Paging, "-rid")
	query := &model.GetHistoryQuery{
		Paging:  paging,
		Table:   "fulfillment",
		Filters: filters,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &etop.HistoryResponse{
		Paging: cmapi.PbPageInfo(paging),
		Data:   cmapi.RawJSONObjectMsg(query.Result.Data),
	}
	return nil
}
