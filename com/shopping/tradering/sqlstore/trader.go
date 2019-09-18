package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/tradering/convert"
	"etop.vn/backend/com/shopping/tradering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type TraderStoreFactory func(ctx context.Context) *TraderStore

func NewTraderStore(db cmsql.Database) TraderStoreFactory {
	return func(ctx context.Context) *TraderStore {
		return &TraderStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type TraderStore struct {
	ft ShopTraderFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging
}

func (s *TraderStore) Paging(paging meta.Paging) *TraderStore {
	s.paging = paging
	return s
}

func (s *TraderStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *TraderStore) Filters(filters meta.Filters) *TraderStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *TraderStore) ID(id int64) *TraderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TraderStore) IDs(ids ...int64) *TraderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *TraderStore) ShopID(id int64) *TraderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *TraderStore) OptionalShopID(id int64) *TraderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *TraderStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ShopTrader)(nil))
}

func (s *TraderStore) GetTraderDB() (*model.ShopTrader, error) {
	query := s.query().Where(s.preds)

	var trader model.ShopTrader
	err := query.ShouldGet(&trader)
	return &trader, err
}

func (s *TraderStore) GetTrader() (traderResult *tradering.ShopTrader, _ error) {
	trader, err := s.GetTraderDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_traderingmodel_ShopTrader_tradering_ShopTrader(trader, traderResult), nil
}

func (s *TraderStore) ListTradersDB() ([]*model.ShopTrader, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, &s.paging, SortTrader)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterTrader)
	if err != nil {
		return nil, err
	}

	var traders model.ShopTraders
	err = query.Find(&traders)
	return traders, err
}

func (s *TraderStore) ListTraders() ([]*tradering.ShopTrader, error) {
	traders, err := s.ListTradersDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_traderingmodel_ShopTraders_tradering_ShopTraders(traders), nil
}
