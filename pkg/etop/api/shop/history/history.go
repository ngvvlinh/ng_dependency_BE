package history

import (
	"context"

	"o.o/api/top/int/etop"
	api "o.o/api/top/int/shop"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/model"
)

type HistoryService struct {
	session.Session
}

func (s *HistoryService) Clone() api.HistoryService { res := *s; return &res }

func (s *HistoryService) GetFulfillmentHistory(ctx context.Context, r *api.GetFulfillmentHistoryRequest) (*etop.HistoryResponse, error) {

	filters := map[string]interface{}{
		"shop_id": s.SS.Shop().ID,
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
		return nil, cm.Error(cm.InvalidArgument, "Must provide either all, id or order_id", nil)
	}

	paging := cmapi.CMPaging(r.Paging, "-rid")
	query := &model.GetHistoryQuery{
		Paging:  paging,
		Table:   "fulfillment",
		Filters: filters,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &etop.HistoryResponse{
		Paging: cmapi.PbPageInfo(paging),
		Data:   cmapi.RawJSONObjectMsg(query.Result.Data),
	}
	return result, nil
}
