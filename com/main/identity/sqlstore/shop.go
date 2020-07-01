package sqlstore

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/meta"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type ShopStoreFactory func(context.Context) *ShopStore

func NewShopStore(db *cmsql.Database) ShopStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopStore {
		return &ShopStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type ShopStore struct {
	query  cmsql.QueryFactory
	shopFt ShopFilters
	preds  []interface{}
	sqlstore.Paging
	filter         meta.Filters
	ctx            context.Context
	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopStore) extend() *ShopStore {
	s.shopFt.prefix = "s"
	return s
}

func (s *ShopStore) WithPaging(paging meta.Paging) *ShopStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *ShopStore) Filters(filters meta.Filters) *ShopStore {
	if s.filter == nil {
		s.filter = filters
	} else {
		s.filter = append(s.filter, filters...)
	}
	return s
}

func (s *ShopStore) ByID(id dot.ID) *ShopStore {
	s.preds = append(s.preds, s.shopFt.ByID(id))
	return s
}

func (s *ShopStore) ByIDs(ids ...dot.ID) *ShopStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShopStore) ByOwnerID(id dot.ID) *ShopStore {
	s.preds = append(s.preds, s.shopFt.ByOwnerID(id))
	return s
}

func (s *ShopStore) GetShopDB() (*identitymodel.Shop, error) {
	var shop identitymodel.Shop
	query := s.query().Where(s.preds)

	// FIX(Tuan): comment vụ check wlPartnerID
	//
	// Webhook NVC cần biết đơn thuộc wlPartnerID nào, hiện tại chưa lưu thông tin này trong order/ffm
	// Tạm thời gọi api GetShopByID từ ffm.ShopID ra để lấy wlPartnerID

	// query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	err := query.ShouldGet(&shop)
	return &shop, err
}

func (s *ShopStore) GetShop() (*identity.Shop, error) {
	shop, err := s.GetShopDB()
	if err != nil {
		return nil, err
	}
	return convert.Shop(shop), nil
}

func (s *ShopStore) ListShopDBs() (res []*identitymodel.Shop, err error) {
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query := s.query().Where(s.preds)
	query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	query, err = sqlstore.LimitSort(query, &s.Paging, map[string]string{"created_at": "created_at"})
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filter, filterShopExtendedWhitelist)

	err = s.query().Where(s.preds).Find((*identitymodel.Shops)(&res))
	return
}

func (s *ShopStore) ListShops() ([]*identity.Shop, error) {
	shops, err := s.ListShopDBs()
	if err != nil {
		return nil, err
	}
	var res []*identity.Shop
	if err := scheme.Convert(shops, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ShopStore) ListShopExtendedDBs() (res []*identitymodel.ShopExtended, err error) {
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query := s.extend().query().Where(s.preds)
	query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	query = s.includeDeleted.Check(query, s.shopFt.NotDeleted())
	query, err = sqlstore.LimitSort(query, &s.Paging, map[string]string{"created_at": "created_at"}, s.shopFt.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filter, filterShopExtendedWhitelist)

	err = query.Find((*identitymodel.ShopExtendeds)(&res))
	return
}
func (s *ShopStore) FilterByWhiteLabelPartner(query cmsql.Query, wlPartnerID dot.ID) cmsql.Query {
	if wlPartnerID != 0 {
		return query.Where(s.shopFt.ByWLPartnerID(wlPartnerID))
	}
	return query.Where(s.shopFt.NotBelongWLPartner())
}

func (s *ShopStore) ListShopExtendeds() ([]*identity.ShopExtended, error) {
	shops, err := s.ListShopExtendedDBs()
	if err != nil {
		return nil, err
	}
	var res []*identity.ShopExtended
	if err := scheme.Convert(shops, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ShopStore) UpdateShop(args *identity.Shop) error {
	var result = &identitymodel.Shop{}
	err := scheme.Convert(args, result)
	if err != nil {
		return err
	}
	err = s.UpdateShopDB(result)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShopStore) UpdateShopDB(args *identitymodel.Shop) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *ShopStore) NotDeleted() *ShopStore {
	s.preds = append(s.preds, sq.NewExpr("deleted_at is null"))
	return s
}

// Only use this function when get model.ShopExtended
func (s *ShopStore) FullTextSearchName(name filter.FullTextSearch) *ShopStore {
	s.preds = append(s.preds, s.shopFt.Filter(`ss.name_norm @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(name)))
	return s
}
