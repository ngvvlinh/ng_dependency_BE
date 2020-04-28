package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/webserver"
	"o.o/backend/com/web/webserver/convert"
	"o.o/backend/com/web/webserver/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type WsPageStoreFactory func(ctx context.Context) *WsPageStore

func NewWsPageStore(db *cmsql.Database) WsPageStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *WsPageStore {
		return &WsPageStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type WsPageStore struct {
	ft WsPageFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
	includeDeleted sqlstore.IncludeDeleted
}

func (s *WsPageStore) WithPaging(paging meta.Paging) *WsPageStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *WsPageStore) Filters(filters meta.Filters) *WsPageStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *WsPageStore) ID(id dot.ID) *WsPageStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *WsPageStore) IDs(ids ...dot.ID) *WsPageStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *WsPageStore) ShopID(id dot.ID) *WsPageStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *WsPageStore) GetWsPageDB() (*model.WsPage, error) {
	query := s.query().Where(s.preds)
	var wsPage model.WsPage
	err := query.ShouldGet(&wsPage)
	return &wsPage, err
}

func (s *WsPageStore) GetWsPage() (wsPageResult *webserver.WsPage, err error) {
	wsPage, err := s.GetWsPageDB()
	if err != nil {
		return nil, err
	}
	wsPageResult = convert.Convert_webservermodel_WsPage_webserver_WsPage(wsPage, wsPageResult)
	return wsPageResult, nil
}

func (s *WsPageStore) Create(args *webserver.WsPage) error {
	var voucherDB = model.WsPage{}
	convert.Convert_webserver_WsPage_webservermodel_WsPage(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *WsPageStore) CreateDB(WsPage *model.WsPage) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(WsPage)
}

func (s *WsPageStore) UpdateWsPageDB(args *model.WsPage) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *WsPageStore) UpdateWsPageAll(args *webserver.WsPage) error {
	var result = &model.WsPage{}
	result = convert.Convert_webserver_WsPage_webservermodel_WsPage(args, result)
	return s.UpdateWsPageDB(result)
}

func (s *WsPageStore) ListWsPagesDB() ([]*model.WsPage, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortWsPage)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterWsPage)
	if err != nil {
		return nil, err
	}

	var wsPages model.WsPages
	err = query.Find(&wsPages)
	return wsPages, err
}

func (s *WsPageStore) ListWsPages() ([]*webserver.WsPage, error) {
	wsPageDB, err := s.ListWsPagesDB()
	if err != nil {
		return nil, err
	}
	wsPage := convert.Convert_webservermodel_WsPages_webserver_WsPages(wsPageDB)
	return wsPage, nil
}

func (s *WsPageStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("ws_page").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}
