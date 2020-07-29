package sqlstore

import (
	"context"
	"strings"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	"o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/catalog/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/capi/filter"
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

func (s *ShopVariantStore) FullTextSearchName(name filter.FullTextSearch) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.Filter(`name_norm @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(name)))
	return s
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

func (s *ShopVariantStore) ExternalID(externalID string) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByExternalID(externalID))
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

func (s *ShopVariantStore) AttributeNorm(attr string) *ShopVariantStore {
	s.preds = append(s.preds, s.FtShopVariant.ByAttrNormKv(attr))
	return s
}

func (s *ShopVariantStore) FilterForImport(args ListShopVariantsForImportArgs) *ShopVariantStore {
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

func (s *ShopVariantStore) CreateShopVariantImport(variant *catalog.ShopVariant) error {
	sqlstore.MustNoPreds(s.preds)
	variantDB := &model.ShopVariant{}
	if err := scheme.Convert(variant, variantDB); err != nil {
		return err
	}
	_, err := s.query().Insert(variantDB)
	return err
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
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortShopVariant, s.FtShopVariant.prefix)
	if err != nil {
		return nil, err
	}

	var variants model.ShopVariants
	err = query.Find(&variants)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(variants)
	return variants, nil
}

func (s *ShopVariantStore) ListShopVariants() ([]*catalog.ShopVariant, error) {
	variantsModel, err := s.ListShopVariantsDB()
	if err != nil {
		return nil, err
	}
	variants := convert.Convert_catalogmodel_ShopVariants_catalog_ShopVariants(variantsModel)
	for i := 0; i < len(variants); i++ {
		variants[i].Deleted = !variantsModel[i].DeletedAt.IsZero()
	}
	return variants, nil
}

func (s *ShopVariantStore) ListShopVariantsWithProductDB() ([]*model.ShopVariantWithProduct, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopVariant.NotDeleted())
	query = s.includeDeleted.Check(query, s.ftShopProduct.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortShopVariant, s.FtShopVariant.prefix)
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

func (s *ShopVariantStore) UpdateShopVariantImport(variant *catalog.ShopVariant) error {
	sqlstore.MustNoPreds(s.preds)
	variantDB := &model.ShopVariant{}
	if err := scheme.Convert(variant, variantDB); err != nil {
		return err
	}
	err := s.query().In("variant_id", variantDB.VariantID).UpdateAll().ShouldUpdate(variantDB)
	return err
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

func CheckShopVariantExternalError(e error, externalID, externalCode string) error {
	if e != nil {
		errMsg := e.Error()
		switch {
		case strings.Contains(errMsg, "shop_variant_shop_id_external_id_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "external_id %v đã tồn tại", externalID)
		case strings.Contains(errMsg, "shop_variant_shop_id_external_code_idx"):
			e = cm.Errorf(cm.FailedPrecondition, e, "external_code %v đã tồn tại", externalCode)
		}
	}
	return e
}

func (s *ShopVariantStore) GetVariantByMaximumCodeNorm(productID dot.ID) (*model.ShopVariant, error) {
	query := s.query().Where(s.preds).Where("code_norm != 0 AND product_id = ?", productID)
	query = query.OrderBy("code_norm desc").Limit(1)

	var variant model.ShopVariant
	if err := query.ShouldGet(&variant); err != nil {
		return nil, err
	}
	return &variant, nil
}
