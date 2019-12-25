package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/tradering/convert"
	"etop.vn/backend/com/shopping/tradering/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type TraderStoreFactory func(ctx context.Context) *TraderStore

func NewTraderStore(db *cmsql.Database) TraderStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *TraderStore {
		return &TraderStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TraderStore struct {
	ft ShopTraderFilters

	query          cmsql.QueryFactory
	preds          []interface{}
	filters        meta.Filters
	paging         meta.Paging
	includeDeleted sqlstore.IncludeDeleted
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

func (s *TraderStore) ID(id dot.ID) *TraderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TraderStore) IDs(ids ...dot.ID) *TraderStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *TraderStore) ShopID(id dot.ID) *TraderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *TraderStore) OptionalShopID(id dot.ID) *TraderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *TraderStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ShopTrader)(nil))
}

func (s *TraderStore) GetTraderDB() (*model.ShopTrader, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
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
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
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

func (s *TraderStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_trader").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}
