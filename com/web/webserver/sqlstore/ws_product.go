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

type WsProductStoreFactory func(ctx context.Context) *WsProductStore

func NewWsProductStore(db *cmsql.Database) WsProductStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *WsProductStore {
		return &WsProductStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type WsProductStore struct {
	ft WsProductFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *WsProductStore) WithPaging(paging meta.Paging) *WsProductStore {
	ss := *s
	ss.Paging.WithPaging(paging)
	return &ss
}

func (s *WsProductStore) Filters(filters meta.Filters) *WsProductStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *WsProductStore) ID(id dot.ID) *WsProductStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *WsProductStore) IDs(ids ...dot.ID) *WsProductStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *WsProductStore) ShopID(id dot.ID) *WsProductStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *WsProductStore) GetWsProductDB() (*model.WsProduct, error) {
	query := s.query().Where(s.preds)
	var wsProduct model.WsProduct
	err := query.ShouldGet(&wsProduct)
	return &wsProduct, err
}

func (s *WsProductStore) GetWsProduct() (wsProductResult *webserver.WsProduct, err error) {
	wsProduct, err := s.GetWsProductDB()
	if err != nil {
		return nil, err
	}
	wsProductResult = convert.Convert_webservermodel_WsProduct_webserver_WsProduct(wsProduct, wsProductResult)
	return wsProductResult, nil
}

func (s *WsProductStore) Create(args *webserver.WsProduct) error {
	var voucherDB = model.WsProduct{}
	convert.Convert_webserver_WsProduct_webservermodel_WsProduct(args, &voucherDB)
	return s.CreateDB(&voucherDB)
}

func (s *WsProductStore) CreateDB(WsProduct *model.WsProduct) error {
	sqlstore.MustNoPreds(s.preds)
	return s.query().ShouldInsert(WsProduct)
}

func (s *WsProductStore) UpdateWsProductDB(args *model.WsProduct) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(args)
}

func (s *WsProductStore) UpdateWsProductAll(args *webserver.WsProduct) error {
	var result = &model.WsProduct{}
	result = convert.Convert_webserver_WsProduct_webservermodel_WsProduct(args, result)
	return s.UpdateWsProductDB(result)
}

func (s *WsProductStore) ListWsProductsDB() ([]*model.WsProduct, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortWsProduct)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterWsProduct)
	if err != nil {
		return nil, err
	}

	var wsProducts model.WsProducts
	err = query.Find(&wsProducts)
	return wsProducts, err
}

func (s *WsProductStore) ListWsProducts() ([]*webserver.WsProduct, error) {
	wsProductDB, err := s.ListWsProductsDB()
	if err != nil {
		return nil, err
	}
	wsProduct := convert.Convert_webservermodel_WsProducts_webserver_WsProducts(wsProductDB)
	return wsProduct, nil
}
