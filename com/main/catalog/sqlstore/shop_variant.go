package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type ListShopVariantsForImportArgs struct {
	Codes     []string
	AttrNorms []interface{}
}

type ShopVariantStoreFactory func(context.Context) *ShopVariantStore

func NewShopVariantStore(db *cmsql.Database) ShopVariantStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopVariantStore {
		return &ShopVariantStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ShopVariantStore struct {
	FtShopVariant ShopVariantFilters
	ftShopProduct ShopProductFilters // unexported

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopVariantStore) extend() *ShopVariantStore {
	s.ftShopProduct.prefix = "sp"
	s.FtShopVariant.prefix = "sv"
	return s
}

func (s *ShopVariantStore) WithPaging(paging meta.Paging) *ShopVariantStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *ShopVariantStore) IncludeDeleted() *ShopVariantStore {
	s.includeDeleted = true
	return s
}

func (s *ShopVariantStore) ID(id dot.ID) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByVariantID(id))
	return s
}
func (s *ShopVariantStore) Code(code string) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByCode(code))
	return s
}

func (s *ShopVariantStore) ProductIDs(productIds ...dot.ID) *ShopVariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.FtShopVariant.prefix, "product_id", productIds))
	return s
}

func (s *ShopVariantStore) IDs(ids ...dot.ID) *ShopVariantStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.FtShopVariant.prefix, "variant_id", ids))
	return s
}

func (s *ShopVariantStore) ShopID(id dot.ID) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByShopID(id))
	return s
}

func (s *ShopVariantStore) OptionalShopID(id dot.ID) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByShopID(id).Optional())
	return s
}

type ListVariantsForImportArgs struct {
	Codes     []string
	AttrNorms []interface{}
}

func (s *ShopVariantStore) FilterForImport(args ListVariantsForImportArgs) *ShopVariantStore {
	pred := sq.Or{
		sq.In("code", args.Codes),
		sq.Ins([]string{"product_id", "attr_norm_kv"}, args.AttrNorms),
	}
	s.preds = append(s.preds, pred)
	return s
}

func (s *ShopVariantStore) CreateShopVariant(variant *catalog.ShopVariant) error {
	sqlstore.MustNoPreds(s.preds)
	variantDB := &model.ShopVariant{}
	if err := scheme.Convert(variant, variantDB); err != nil {
		return err
	}
	_, err := s.query().Insert(variantDB)
	return checkProductOrVariantError(err, variantDB.Code)
}

func (s *ShopVariantStore) GetShopVariantDB() (*model.ShopVariant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())

	var variant model.ShopVariant
	err := query.ShouldGet(&variant)
	return &variant, err
}

func (s *ShopVariantStore) GetShopVariant() (*catalog.ShopVariant, error) {
	variantDB, err := s.GetShopVariantDB()
	if err != nil {
		return nil, err
	}
	variant := &catalog.ShopVariant{}
	err = scheme.Convert(variantDB, variant)
	return variant, err
}

func (s *ShopVariantStore) GetShopVariantWithProductDB() (*model.ShopVariantWithProduct, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())

	var variant model.ShopVariantWithProduct
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

func (s *ShopVariantStore) ListShopVariantsDB() ([]*model.ShopVariant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.Paging, SortShopVariant, s.FtShopVariant.prefix)
	if err != nil {
		return nil, err
	}

	var variants model.ShopVariants
	err = query.Find(&variants)
	return variants, err
}

func (s *ShopVariantStore) ListShopVariants() ([]*catalog.ShopVariant, error) {
	variantsModel, err := s.ListShopVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_catalogmodel_ShopVariants_catalog_ShopVariants(variantsModel), nil
}

func (s *ShopVariantStore) ListShopVariantsWithProductDB() ([]*model.ShopVariantWithProduct, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.ftShopProduct.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.PrefixedLimitSort(query, &s.Paging, SortShopVariant, s.FtShopVariant.prefix)
	if err != nil {
		return nil, err
	}

	var variants model.ShopVariantWithProducts
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

func (s *ShopVariantStore) UpdateShopVariant(variant *model.ShopVariant) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().In("variant_id", variant.VariantID).UpdateAll().ShouldUpdate(variant)
	return checkProductOrVariantError(err, variant.Code)
}

func (s *ShopVariantStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	_deleted, err := query.Table("shop_variant").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ShopVariantStore) UpdateStatusShopVariant(status int16) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	updateStatus, err := query.Table("shop_variant").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return updateStatus, err
}

func (s *ShopVariantStore) UpdateImageShopVariant(variant *catalog.ShopVariant) error {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	variantDB := &model.ShopVariant{}
	if err := scheme.Convert(variant, variantDB); err != nil {
		return err
	}
	err := query.ShouldUpdate(variantDB)
	return err
}
