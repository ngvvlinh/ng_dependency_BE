package sqlstore

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/backend/pkg/services/catalog/convert"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

type ShopVariantStoreFactory func(context.Context) *ShopVariantStore

func NewShopVariantStore(db cmsql.Database) ShopVariantStoreFactory {
	return func(ctx context.Context) *ShopVariantStore {
		return &ShopVariantStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShopVariantStore struct {
	FtVariant     VariantFilters
	FtShopVariant ShopVariantFilters

	// unexported
	ftProduct     ProductFilters
	ftShopProduct ShopProductFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopVariantStore) extend() *ShopVariantStore {
	s.ftProduct.prefix = "p"
	s.FtVariant.prefix = "v"
	s.ftShopProduct.prefix = "sp"
	s.FtShopVariant.prefix = "sv"
	return s
}

func (s *ShopVariantStore) Paging(paging meta.Paging) *ShopVariantStore {
	s.paging = paging
	return s
}

func (s *ShopVariantStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopVariantStore) ID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, id)
	return s
}

func (s *ShopVariantStore) IDs(ids ...int64) *ShopVariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.FtVariant.prefix, "id", ids))
	return s
}

func (s *ShopVariantStore) ProductSourceID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, s.FtVariant.ByProductSourceID(id))
	return s
}

func (s *ShopVariantStore) ShopID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByShopID(id))
	return s
}

func (s *ShopVariantStore) OptionalProductSourceID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, s.FtVariant.ByProductSourceID(id).Optional())
	return s
}

func (s *ShopVariantStore) GetShopVariantDB() (*catalogmodel.ShopVariantExtended, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())

	var variant catalogmodel.ShopVariantExtended
	err := query.ShouldGet(&variant)
	return &variant, err
}

func (s *ShopVariantStore) GetShopVariant() (*catalog.ShopVariantExtended, error) {
	variant, err := s.GetShopVariantDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopVariantExtended(variant), nil
}

func (s *ShopVariantStore) GetShopVariantWithProductDB() (*catalogmodel.ShopVariantWithProduct, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())

	var variant catalogmodel.ShopVariantWithProduct
	err := query.ShouldGet(&variant)
	return &variant, err
}

func (s *ShopVariantStore) GetShopVariantWithProduct() (*catalog.ShopVariantWithProduct, error) {
	variant, err := s.GetShopVariantWithProductDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopVariantWithProduct(variant), nil
}

func (s *ShopVariantStore) ListShopVariantsDB() ([]*catalogmodel.ShopVariantExtended, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortVariant)
	if err != nil {
		return nil, err
	}

	var variants catalogmodel.ShopVariantExtendeds
	err = query.Find(&variants)
	return variants, err
}

func (s *ShopVariantStore) ListShopVariants() ([]*catalog.ShopVariantExtended, error) {
	variants, err := s.ListShopVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopVariantsExtended(variants), nil

}

func (s *ShopVariantStore) ListShopVariantsWithProductDB() ([]*catalogmodel.ShopVariantWithProduct, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.ftProduct.NotDeleted())
	query = s.includeDeleted.Check(query, s.ftShopProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortVariant)
	if err != nil {
		return nil, err
	}

	var variants catalogmodel.ShopVariantWithProducts
	err = query.Find(&variants)
	return variants, err
}

func (s *ShopVariantStore) ListShopVariantsWithProduct() ([]*catalog.ShopVariantWithProduct, error) {
	variants, err := s.ListShopVariantsWithProductDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopVariantsWithProduct(variants), nil
}
