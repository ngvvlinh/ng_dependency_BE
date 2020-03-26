package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/api/webserver"
	"o.o/backend/com/web/webserver/convert"
	"o.o/backend/com/web/webserver/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type WsWebsiteStoreFactory func(ctx context.Context) *WsWebsiteStore

func NewWsWebsiteStore(db *cmsql.Database) WsWebsiteStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *WsWebsiteStore {
		return &WsWebsiteStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type WsWebsiteStore struct {
	ft WsWebsiteFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *WsWebsiteStore) WithPaging(paging meta.Paging) *WsWebsiteStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *WsWebsiteStore) Filters(filters meta.Filters) *WsWebsiteStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *WsWebsiteStore) ID(id dot.ID) *WsWebsiteStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *WsWebsiteStore) IDs(ids ...dot.ID) *WsWebsiteStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *WsWebsiteStore) ShopID(id dot.ID) *WsWebsiteStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *WsWebsiteStore) SiteSubdomain(site string) *WsWebsiteStore {
	s.preds = append(s.preds, s.ft.BySiteSubdomain(site))
	return s
}

func (s *WsWebsiteStore) GetWsWebsiteDB() (*model.WsWebsite, error) {
	query := s.query().Where(s.preds)
	var wsWebsite model.WsWebsite
	err := query.ShouldGet(&wsWebsite)
	return &wsWebsite, err
}

func (s *WsWebsiteStore) GetWsWebsite() (wsWebsiteResult *webserver.WsWebsite, err error) {
	wsWebsite, err := s.GetWsWebsiteDB()
	if err != nil {
		return nil, err
	}
	wsWebsiteResult = convert.Convert_webservermodel_WsWebsite_webserver_WsWebsite(wsWebsite, wsWebsiteResult)
	return wsWebsiteResult, nil
}

func (s *WsWebsiteStore) Create(args *webserver.WsWebsite) error {
	var voucherDB = model.WsWebsite{}
	convert.Convert_webserver_WsWebsite_webservermodel_WsWebsite(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *WsWebsiteStore) CreateDB(WsWebsite *model.WsWebsite) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(WsWebsite)
}

func (s *WsWebsiteStore) UpdateWsWebsiteDB(args *model.WsWebsite) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *WsWebsiteStore) UpdateWsWebsiteAll(args *webserver.WsWebsite) error {
	var result = &model.WsWebsite{}
	result = convert.Convert_webserver_WsWebsite_webservermodel_WsWebsite(args, result)
	return s.UpdateWsWebsiteDB(result)
}

func (s *WsWebsiteStore) ListWsWebsitesDB() ([]*model.WsWebsite, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortWsWebsite)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterWsWebsite)
	if err != nil {
		return nil, err
	}

	var wsWebsites model.WsWebsites
	err = query.Find(&wsWebsites)
	return wsWebsites, err
}

func (s *WsWebsiteStore) ListWsWebsites() ([]*webserver.WsWebsite, error) {
	wsWebsiteDB, err := s.ListWsWebsitesDB()
	if err != nil {
		return nil, err
	}
	wsWebsite := convert.Convert_webservermodel_WsWebsites_webserver_WsWebsites(wsWebsiteDB)
	return wsWebsite, nil
}
