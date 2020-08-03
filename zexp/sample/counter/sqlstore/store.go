package sqlstore

import (
	"context"

	"o.o/api/meta"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/zexp/sample/counter/api"
	convert "o.o/backend/zexp/sample/counter/convert"
	"o.o/backend/zexp/sample/counter/model"
)

type CounterStoreFactory func(ctx context.Context) *CounterStore

func NewCounterStore(db *cmsql.Database) CounterStoreFactory {
	return func(ctx context.Context) *CounterStore {
		return &CounterStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type CounterStore struct {
	ft      CounterFilters
	query   cmsql.QueryFactory
	preds   []interface{}
	ctx     context.Context
	filters meta.Filters
	sqlstore.Paging
}

func (s *CounterStore) WithPaging(paging meta.Paging) *CounterStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *CounterStore) Label(label string) *CounterStore {
	s.preds = append(s.preds, s.ft.ByLabel(label))
	return s
}

func (s *CounterStore) CreateCounterDB(co *model.Counter) error {
	return s.query().ShouldInsert(co)
}

func (s *CounterStore) CreateCounter(co *api.Counter) error {
	var out *model.Counter
	out = convert.Convert_api_Counter_countermodel_Counter(co, out)
	return s.CreateCounterDB(out)
}

func (s *CounterStore) GetCounterDB() (*model.Counter, error) {
	eq := &model.Counter{}
	err := s.query().Where(s.preds).ShouldGet(eq)
	return eq, err
}

func (s *CounterStore) UpdateLabelvalue(value int) (updated int, err error) {
	if len(s.preds) == 0 {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	update := &model.Counter{
		Value: value,
	}
	return s.query().Where(s.preds).Update(update)
}

func (s *CounterStore) GetCounter() (*api.Counter, error) {
	coDB, err := s.GetCounterDB()
	if err != nil {
		return nil, err
	}
	var out *api.Counter
	out = convert.Convert_countermodel_Counter_api_Counter(coDB, out)

	return out, nil
}
