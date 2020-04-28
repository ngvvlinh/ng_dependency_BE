package query

import (
	"context"

	"o.o/api/meta"
	crmvht "o.o/api/supporting/crm/vht"
	"o.o/backend/com/supporting/crm/vht/convert"
	"o.o/backend/com/supporting/crm/vht/model"
	"o.o/backend/com/supporting/crm/vht/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

type QueryService struct {
	VhtStore sqlstore.VhtCallHistoriesFactory
}

var _ crmvht.QueryService = &QueryService{}

func New(db *cmsql.Database) *QueryService {
	return &QueryService{
		VhtStore: sqlstore.NewVhtCallHistoryStore(db),
	}
}

func (q *QueryService) MessageBus() crmvht.QueryBus {
	b := bus.New()
	return crmvht.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q QueryService) GetCallHistories(ctx context.Context, req *crmvht.GetCallHistoriesArgs) (*crmvht.GetCallHistoriesResponse, error) {

	var paging meta.Paging
	paging = *req.Paging
	textSearch := req.TextSearch
	var dbResult []*model.VhtCallHistory
	var err error
	if textSearch != "" {
		dbResult, err = q.VhtStore(ctx).WithPaging(paging).SearchVhtCallHistories(textSearch)
	} else {
		dbResult, err = q.VhtStore(ctx).WithPaging(paging).GetCallHistories()
	}
	if err != nil {
		return nil, err
	}
	var vhtCallHistoryResult []*crmvht.VhtCallLog
	for i := 0; i < len(dbResult); i++ {
		dataRow := dbResult[i]
		callHistoryRow := convert.ConvertFromModel(dataRow)
		vhtCallHistoryResult = append(vhtCallHistoryResult, callHistoryRow)
	}
	return &crmvht.GetCallHistoriesResponse{
		VhtCallLog: vhtCallHistoryResult,
	}, nil
}

func (q *QueryService) GetLastCallHistory(ctx context.Context, paging meta.Paging) (*crmvht.VhtCallLog, error) {
	result, err := q.VhtStore(ctx).WithPaging(paging).ByStatus("Done").SortBy("time_started desc").GetCallHistories()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return convert.ConvertFromModel(result[0]), nil
}
