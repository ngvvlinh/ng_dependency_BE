package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/com/supporting/crm/vht/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/backend/pkg/common/validate"
)

type VhtCallHistoriesFactory func(context.Context) *VhtCallHistoryStore

type VhtCallHistoryStore struct {
	query cmsql.QueryFactory
	preds []interface{}
	ft    VhtCallHistoryFilters
	sqlstore.Paging
	OrderBy string
}

func NewVhtCallHistoryStore(db *cmsql.Database) VhtCallHistoriesFactory {
	return func(ctx context.Context) *VhtCallHistoryStore {
		return &VhtCallHistoryStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

var SortVhtCallHistories = map[string]string{
	"time_started": "-time_started",
}

func (s *VhtCallHistoryStore) WithPaging(paging meta.Paging) *VhtCallHistoryStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *VhtCallHistoryStore) ByStatus(value string) *VhtCallHistoryStore {
	s.preds = append(s.preds, s.ft.BySyncStatus(value))
	return s
}

func (s *VhtCallHistoryStore) ByCallID(value string) *VhtCallHistoryStore {
	s.preds = append(s.preds, s.ft.ByCallID(value))
	return s
}

func (s *VhtCallHistoryStore) SortBy(order string) *VhtCallHistoryStore {
	s.OrderBy = order
	return s
}

func (s *VhtCallHistoryStore) BySdkCallID(value string) *VhtCallHistoryStore {
	s.preds = append(s.preds, s.ft.BySdkCallID(value))
	return s
}

func (s *VhtCallHistoryStore) GetCallHistory() (*model.VhtCallHistory, error) {
	query := s.query().Where(s.preds)
	var vthCallHistory model.VhtCallHistory
	err := query.ShouldGet(&vthCallHistory)
	if err != nil {
		return nil, err
	}
	return &vthCallHistory, nil
}

func (s *VhtCallHistoryStore) GetCallHistories() ([]*model.VhtCallHistory, error) {
	query := s.query().Where(s.preds)
	if s.OrderBy != "" {
		query.OrderBy(s.OrderBy)
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortVhtCallHistories)
	if err != nil {
		return nil, err
	}
	var vthCallHistories []*model.VhtCallHistory
	err = query.Find((*model.VhtCallHistories)(&vthCallHistories))
	return vthCallHistories, err
}

func (s *VhtCallHistoryStore) SearchVhtCallHistories(value string) ([]*model.VhtCallHistory, error) {
	query := s.query().Where(`search_norm @@ ?::tsquery`, validate.NormalizeSearchQueryAnd(value))
	if s.OrderBy != "" {
		query.OrderBy(s.OrderBy)
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortVhtCallHistories)
	if err != nil {
		return nil, err
	}
	var vthCallHistories []*model.VhtCallHistory
	err = query.Find((*model.VhtCallHistories)(&vthCallHistories))
	return vthCallHistories, err
}

func (s *VhtCallHistoryStore) CreateVhtCallHistory(contact *model.VhtCallHistory) error {
	query := s.query().Where(s.preds)
	err := query.ShouldInsert(contact)
	return err
}

func (s *VhtCallHistoryStore) UpdateVhtCallHistory(contact *model.VhtCallHistory) error {
	query := s.query().Where(s.preds)
	err := query.ShouldUpdate(contact)
	return err
}
