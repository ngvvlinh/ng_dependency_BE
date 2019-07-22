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

type ListShopVariantsForImportArgs struct {
	Codes     []string
	AttrNorms []interface{}
}

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
	FtShopVariant ShopVariantFilters
	ftShopProduct ShopProductFilters // unexported

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopVariantStore) Paging(paging meta.Paging) *ShopVariantStore {
	s.paging = paging
	return s
}

func (s *ShopVariantStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopVariantStore) ID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByVariantID(id))
	return s
}

func (s *ShopVariantStore) IDs(ids ...int64) *ShopVariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.FtShopVariant.prefix, "variant_id", ids))
	return s
}

func (s *ShopVariantStore) ShopID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByShopID(id))
	return s
}

func (s *ShopVariantStore) OptionalShopID(id int64) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByShopID(id).Optional())
	return s
}

type ListVariantsForImportArgs struct {
	Codes     []string
	AttrNorms []interface{}
}

func (s *ShopVariantStore) FilterForImport(args ListVariantsForImportArgs) *ShopVariantStore {
	pred := sq.Or{
		sq.In("ed_code", args.Codes),
		sq.Ins([]string{"product_id", "attr_norm_kv"}, args.AttrNorms),
	}
	s.preds = append(s.preds, pred)
	return s
}

func (s *ShopVariantStore) GetShopVariantDB() (*catalogmodel.ShopVariant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())

	var variant catalogmodel.ShopVariant
	err := query.ShouldGet(&variant)
	return &variant, err
}

func (s *ShopVariantStore) GetShopVariant() (*catalog.ShopVariant, error) {
	variant, err := s.GetShopVariantDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopVariant(variant), nil
}

func (s *ShopVariantStore) GetShopVariantWithProductDB() (*catalogmodel.ShopVariantWithProduct, error) {
	query := s.query().Where(s.preds)
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

func (s *ShopVariantStore) ListShopVariantsDB() ([]*catalogmodel.ShopVariant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortVariant)
	if err != nil {
		return nil, err
	}

	var variants catalogmodel.ShopVariants
	err = query.Find(&variants)
	return variants, err
}

func (s *ShopVariantStore) ListShopVariants() ([]*catalog.ShopVariant, error) {
	variants, err := s.ListShopVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopVariants(variants), nil

}

func (s *ShopVariantStore) ListShopVariantsWithProductDB() ([]*catalogmodel.ShopVariantWithProduct, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
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
