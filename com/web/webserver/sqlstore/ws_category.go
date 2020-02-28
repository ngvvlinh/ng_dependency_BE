package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/webserver"
	"etop.vn/backend/com/web/webserver/convert"
	"etop.vn/backend/com/web/webserver/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type WsCategoryStoreFactory func(ctx context.Context) *WsCategoryStore

func NewWsCategoryStore(db *cmsql.Database) WsCategoryStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *WsCategoryStore {
		return &WsCategoryStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type WsCategoryStore struct {
	ft WsCategoryFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *WsCategoryStore) WithPaging(paging meta.Paging) *WsCategoryStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *WsCategoryStore) Filters(filters meta.Filters) *WsCategoryStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *WsCategoryStore) ID(id dot.ID) *WsCategoryStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *WsCategoryStore) IDs(ids ...dot.ID) *WsCategoryStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *WsCategoryStore) ShopID(id dot.ID) *WsCategoryStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *WsCategoryStore) GetWsCategoryDB() (*model.WsCategory, error) {
	query := s.query().Where(s.preds)
	var wsCategory model.WsCategory
	err := query.ShouldGet(&wsCategory)
	return &wsCategory, err
}

func (s *WsCategoryStore) GetWsCategory() (wsCategoryResult *webserver.WsCategory, err error) {
	wsCategory, err := s.GetWsCategoryDB()
	if err != nil {
		return nil, err
	}
	wsCategoryResult = convert.Convert_webservermodel_WsCategory_webserver_WsCategory(wsCategory, wsCategoryResult)
	return wsCategoryResult, nil
}

func (s *WsCategoryStore) Create(args *webserver.WsCategory) error {
	var voucherDB = model.WsCategory{}
	convert.Convert_webserver_WsCategory_webservermodel_WsCategory(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *WsCategoryStore) CreateDB(wsCategory *model.WsCategory) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(wsCategory)
}

func (s *WsCategoryStore) UpdateWsCategoryDB(args *model.WsCategory) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *WsCategoryStore) UpdateWsCategoryAll(args *webserver.WsCategory) error {
	var result = &model.WsCategory{}
	result = convert.Convert_webserver_WsCategory_webservermodel_WsCategory(args, result)
	return s.UpdateWsCategoryDB(result)
}

func (s *WsCategoryStore) ListWsCategorysDB() ([]*model.WsCategory, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortWsCategory)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterWsCategory)
	if err != nil {
		return nil, err
	}

	var wsCategories model.WsCategories
	err = query.Find(&wsCategories)
	return wsCategories, err
}

func (s *WsCategoryStore) ListWsCategories() ([]*webserver.WsCategory, error) {
	wsCategoriesDB, err := s.ListWsCategorysDB()
	if err != nil {
		return nil, err
	}
	wsCategories := convert.Convert_webservermodel_WsCategories_webserver_WsCategories(wsCategoriesDB)
	return wsCategories, nil
}
