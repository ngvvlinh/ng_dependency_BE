package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type ShopBrandStoreFactory func(context.Context) *ShopBrandStore

func NewShopBrandStore(db *cmsql.Database) ShopBrandStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopBrandStore {
		return &ShopBrandStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopBrandStore struct {
	ftShopBrand ShopBrandFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopBrandStore) Paging(paging meta.Paging) *ShopBrandStore {
	s.paging = paging
	return s
}

func (s *ShopBrandStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopBrandStore) ID(id dot.ID) *ShopBrandStore {
	s.preds = append(s.preds, s.ftShopBrand.ByID(id))
	return s
}

func (s *ShopBrandStore) IDs(ids ...dot.ID) *ShopBrandStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShopBrandStore) ShopID(id dot.ID) *ShopBrandStore {
	s.preds = append(s.preds, s.ftShopBrand.ByShopID(id))
	return s
}

func (s *ShopBrandStore) BrandName(name string) *ShopBrandStore {
	s.preds = append(s.preds, s.ftShopBrand.ByBrandName(name))
	return s
}

func (s *ShopBrandStore) CreateShopBrand(brand *catalog.ShopBrand) error {
	sqlstore.MustNoPreds(s.preds)
	brandDB := convert.Convert_catalog_ShopBrand_catalogmodel_ShopBrand(brand, nil)
	err := s.query().ShouldInsert(brandDB)
	return err
}

func (s *ShopBrandStore) UpdateShopBrand(brand *catalog.ShopBrand) error {
	brandDB := convert.Convert_catalog_ShopBrand_catalogmodel_ShopBrand(brand, nil)

	err := s.query().Where(s.preds).ShouldUpdate(brandDB)
	return err
}

func (s *ShopBrandStore) GetShopBrandDB() (*model.ShopBrand, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopBrand.NotDeleted())

	var brand model.ShopBrand
	err := query.ShouldGet(&brand)
	return &brand, err
}

func (s *ShopBrandStore) GetShopBrand() (*catalog.ShopBrand, error) {
	brand, err := s.GetShopBrandDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_catalogmodel_ShopBrand_catalog_ShopBrand(brand, nil), nil
}

func (s *ShopBrandStore) ListShopBrandsDB() ([]*model.ShopBrand, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopBrand.NotDeleted())
	if len(s.paging.Sort) == 0 {
		s.paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.paging, SortShopBrand, s.ftShopBrand.prefix)
	if err != nil {
		return nil, err
	}

	var brands model.ShopBrands
	err = query.Find(&brands)
	return brands, err
}
func (s *ShopBrandStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopBrand.NotDeleted())
	return query.Count((*model.ShopBrand)(nil))
}

func (s *ShopBrandStore) ListShopBrands() ([]*catalog.ShopBrand, error) {
	products, err := s.ListShopBrandsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_catalogmodel_ShopBrands_catalog_ShopBrands(products), nil
}

func (s *ShopBrandStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftShopBrand.NotDeleted())
	_deleted, err := query.Table("shop_brand").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}
