package sqlstore

import (
	"context"

	model2 "etop.vn/backend/com/supporting/crm/vht/model"

	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/backend/pkg/common/validate"
)

type VhtCallHistoriesFactory func(context.Context) *VhtCallHistoryStore

type VhtCallHistoryStore struct {
	query   func() cmsql.QueryInterface
	preds   []interface{}
	ft      VhtCallHistoryFilters
	paging  meta.Paging
	OrderBy string
}

func NewVhtCallHistoryStore(db cmsql.Database) VhtCallHistoriesFactory {
	return func(ctx context.Context) *VhtCallHistoryStore {
		return &VhtCallHistoryStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

var SortVhtCallHistories = map[string]string{
	"time_started": "-time_started",
}

func (s *VhtCallHistoryStore) Paging(paging meta.Paging) *VhtCallHistoryStore {
	s.paging = paging
	return s
}

func (s *VhtCallHistoryStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (v *VhtCallHistoryStore) ByStatus(value string) *VhtCallHistoryStore {
	v.preds = append(v.preds, v.ft.BySyncStatus(value))
	return v
}

func (v *VhtCallHistoryStore) ByCallID(value string) *VhtCallHistoryStore {
	v.preds = append(v.preds, v.ft.ByCallID(value))
	return v
}

func (v *VhtCallHistoryStore) SortBy(order string) *VhtCallHistoryStore {
	v.OrderBy = order
	return v
}

func (v *VhtCallHistoryStore) BySdkCallID(value string) *VhtCallHistoryStore {
	v.preds = append(v.preds, v.ft.BySdkCallID(value))
	return v
}

func (v *VhtCallHistoryStore) GetCallHistory() (*model2.VhtCallHistory, error) {
	query := v.query().Where(v.preds)
	var vthCallHistory model2.VhtCallHistory
	err := query.ShouldGet(&vthCallHistory)
	if err != nil {
		return nil, err
	}
	return &vthCallHistory, nil
}

func (v *VhtCallHistoryStore) GetCallHistories() ([]*model2.VhtCallHistory, error) {
	query := v.query().Where(v.preds)
	if v.OrderBy != "" {
		query.OrderBy(v.OrderBy)
	}
	query, err := sqlstore.LimitSort(query, &v.paging, SortVhtCallHistories)
	if err != nil {
		return nil, err
	}
	var vthCallHistories []*model2.VhtCallHistory
	err = query.Find((*model2.VhtCallHistories)(&vthCallHistories))
	return vthCallHistories, err
}

func (v *VhtCallHistoryStore) SearchVhtCallHistories(value string) ([]*model2.VhtCallHistory, error) {
	query := v.query().Where(`search_norm @@ ?::tsquery`, validate.NormalizeSearchQueryAnd(value))
	if v.OrderBy != "" {
		query.OrderBy(v.OrderBy)
	}
	query, err := sqlstore.LimitSort(query, &v.paging, SortVhtCallHistories)
	if err != nil {
		return nil, err
	}
	var vthCallHistories []*model2.VhtCallHistory
	err = query.Find((*model2.VhtCallHistories)(&vthCallHistories))
	return vthCallHistories, err
}

func (v *VhtCallHistoryStore) CreateVhtCallHistory(contact *model2.VhtCallHistory) error {
	query := v.query().Where(v.preds)
	err := query.ShouldInsert(contact)
	return err
}

func (v *VhtCallHistoryStore) UpdateVhtCallHistory(contact *model2.VhtCallHistory) error {
	query := v.query().Where(v.preds)
	err := query.ShouldUpdate(contact)
	return err
}
