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

type VariantStoreFactory func(context.Context) *VariantStore

func NewVariantStore(db cmsql.Database) VariantStoreFactory {
	return func(ctx context.Context) *VariantStore {
		return &VariantStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type VariantStore struct {
	FtProduct ProductFilters

	// unexported
	ftVariant VariantFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *VariantStore) extend() *VariantStore {
	s.FtProduct.prefix = "p"
	s.ftVariant.prefix = "v"
	return s
}

func (s *VariantStore) Paging(paging meta.Paging) *VariantStore {
	s.paging = paging
	return s
}

func (s *VariantStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *VariantStore) ID(id int64) *VariantStore {
	s.preds = append(s.preds, s.ftVariant.ByID(id))
	return s
}

func (s *VariantStore) IDs(ids ...int64) *VariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ftVariant.prefix, "id", ids))
	return s
}

func (s *VariantStore) ProductSourceID(id int64) *VariantStore {
	s.preds = append(s.preds, s.ftVariant.ByProductSourceID(id))
	return s
}
func (s *VariantStore) ByAttrNorms(attrNorms ...string) *VariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ftVariant.prefix, "attr_norm_kv", attrNorms))
	return s
}

type ListVariantsForImportArgs struct {
	Codes     []string
	AttrNorms []interface{}
}

func (s *VariantStore) FilterForImport(args ListVariantsForImportArgs) *VariantStore {
	pred := sq.Or{
		sq.In("ed_code", args.Codes),
		sq.Ins([]string{"product_id", "attr_norm_kv"}, args.AttrNorms),
	}
	s.preds = append(s.preds, pred)
	return s
}

func (s *VariantStore) GetVariantDB() (*catalogmodel.Variant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftVariant.NotDeleted())

	var variant catalogmodel.Variant
	err := query.ShouldGet(&variant)
	return &variant, err
}

func (s *VariantStore) GetVariant() (*catalog.Variant, error) {
	variant, err := s.GetVariantDB()
	if err != nil {
		return nil, err
	}
	return convert.Variant(variant), nil
}

func (s *VariantStore) GetVariantWithProductDB() (*catalogmodel.VariantExtended, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())

	var variant catalogmodel.VariantExtended
	err := query.ShouldGet(&variant)
	return &variant, err
}

func (s *VariantStore) GetVariantWithProduct() (*catalog.VariantWithProduct, error) {
	variant, err := s.GetVariantWithProductDB()
	if err != nil {
		return nil, err
	}
	return convert.VariantWithProduct(variant), nil
}

func (s *VariantStore) ListVariantsDB() ([]*catalogmodel.Variant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftVariant.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortVariant)
	if err != nil {
		return nil, err
	}
	var variants catalogmodel.Variants
	err = query.Find(&variants)
	return variants, err
}

func (s *VariantStore) ListVariants() ([]*catalog.Variant, error) {
	variants, err := s.ListVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.Variants(variants), nil
}

func (s *VariantStore) ListVariantsWithProductDB() ([]*catalogmodel.VariantExtended, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ftVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortVariant)
	if err != nil {
		return nil, err
	}
	var variants catalogmodel.VariantExtendeds
	err = query.Find(&variants)
	return variants, err
}

func (s *VariantStore) ListVariantsWithProduct() ([]*catalog.VariantWithProduct, error) {
	variants, err := s.ListVariantsWithProductDB()
	if err != nil {
		return nil, err
	}
	return convert.VariantsWithProduct(variants), nil
}
