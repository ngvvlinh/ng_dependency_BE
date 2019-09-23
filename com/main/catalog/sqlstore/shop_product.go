package sqlstore

import (
	"context"
	"strings"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type ShopProductStoreFactory func(context.Context) *ShopProductStore

func NewShopProductStore(db cmsql.Database) ShopProductStoreFactory {
	return func(ctx context.Context) *ShopProductStore {
		return &ShopProductStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShopProductStore struct {
	FtShopProduct ShopProductFilters
	ftShopVariant ShopVariantFilters // unexported

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopProductStore) Paging(paging meta.Paging) *ShopProductStore {
	s.paging = paging
	return s
}

func (s *ShopProductStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ShopProductStore) Where(pred sq.FilterQuery) *ShopProductStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ShopProductStore) Filters(filters meta.Filters) *ShopProductStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ShopProductStore) ID(id int64) *ShopProductStore {
	s.preds = append(s.preds, s.FtShopProduct.ByProductID(id))
	return s
}

func (s *ShopProductStore) CategoryID(id int64) *ShopProductStore {
	s.preds = append(s.preds, s.FtShopProduct.ByCategoryID(id))
	return s
}

func (s *ShopProductStore) IDs(ids ...int64) *ShopProductStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *ShopProductStore) ShopID(id int64) *ShopProductStore {
	s.preds = append(s.preds, s.FtShopProduct.ByShopID(id))
	return s
}

func (s *ShopProductStore) OptionalShopID(id int64) *ShopProductStore {
	s.preds = append(s.preds, s.FtShopProduct.ByShopID(id).Optional())
	return s
}

func (s *ShopProductStore) Code(code string) *ShopProductStore {
	s.preds = append(s.preds, s.FtShopProduct.ByCode(code))
	return s
}

func (s *ShopProductStore) Codes(codes ...string) *ShopProductStore {
	s.preds = append(s.preds, sq.In("code", codes))
	return s
}

func (s *ShopProductStore) ByNameNormUas(names ...string) *ShopProductStore {
	s.preds = append(s.preds, sq.In("name_norm_ua", names))
	return s
}

func (s *ShopProductStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	return query.Count((*model.ShopProduct)(nil))
}

func (s *ShopProductStore) CreateShopProduct(product *catalog.ShopProduct) error {
	sqlstore.MustNoPreds(s.preds)
	productDB := convert.ShopProductDB(product)
	_, err := s.query().Insert(productDB)
	return checkProductOrVariantError(err, productDB.Code)
}

func (s *ShopProductStore) GetShopProductDB() (*model.ShopProduct, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())

	var product model.ShopProduct
	err := query.ShouldGet(&product)
	return &product, err
}

func (s *ShopProductStore) GetShopProduct() (*catalog.ShopProduct, error) {
	product, err := s.GetShopProductDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProduct(product), nil
}

func (s *ShopProductStore) GetShopProductWithVariantsDB() (*model.ShopProductWithVariants, error) {
	product, err := s.GetShopProductDB()
	if err != nil {
		return nil, err
	}
	var variants model.ShopVariants
	{
		q := s.query().OrderBy("variant_id").
			Where(s.ftShopVariant.ByProductID(product.ProductID))
		q = s.includeDeleted.Check(q, s.ftShopVariant.NotDeleted())
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}
	return &model.ShopProductWithVariants{
		ShopProduct: product,
		Variants:    variants,
	}, nil
}

func (s *ShopProductStore) GetShopProductWithVariants() (*catalog.ShopProductWithVariants, error) {
	product, err := s.GetShopProductWithVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProductWithVariants(product), nil
}

func (s *ShopProductStore) ListShopProductsDB() ([]*model.ShopProduct, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortShopProduct)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProduct)
	if err != nil {
		return nil, err
	}

	var products model.ShopProducts
	err = query.OrderBy("created_at DESC").Find(&products)
	return products, err
}

func (s *ShopProductStore) ListShopProducts() ([]*catalog.ShopProduct, error) {
	products, err := s.ListShopProductsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProducts(products), nil
}

func (s *ShopProductStore) ListShopProductsWithVariantsDB() ([]*model.ShopProductWithVariants, error) {
	products, err := s.ListShopProductsDB()
	if err != nil {
		return nil, err
	}

	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.ProductID
	}

	var variants model.ShopVariants
	{
		q := s.query().In("product_id", productIDs)
		q = s.includeDeleted.Check(q, s.ftShopVariant.NotDeleted())
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}

	mapProducts := make(map[int64]*model.ShopProductWithVariants)
	result := make([]*model.ShopProductWithVariants, len(products))
	for i, p := range products {
		result[i] = &model.ShopProductWithVariants{
			ShopProduct: p,
		}
		mapProducts[p.ProductID] = result[i]
	}
	for _, v := range variants {
		p := mapProducts[v.ProductID]
		if p != nil {
			p.Variants = append(p.Variants, v)
		}
	}
	return result, nil
}

func (s *ShopProductStore) ListShopProductsWithVariants() ([]*catalog.ShopProductWithVariants, error) {
	products, err := s.ListShopProductsWithVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProductsWithVariants(products), nil
}

func (s *ShopProductStore) UpdateShopProduct(product *model.ShopProduct) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().In("product_id", product.ProductID).UpdateAll().ShouldUpdate(product)
	return checkProductOrVariantError(err, product.Code)
}

func (s *ShopProductStore) UpdateShopProductCategory(product *model.ShopProduct) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Where(s.FtShopProduct.ByProductID((product.ProductID))).Update(product)
	return err
}

func (s *ShopProductStore) UpdateStatusShopProducts(status int16) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	updateStatus, err := query.Table("shop_product").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return int(updateStatus), err
}

func (s *ShopProductStore) UpdateImageShopProduct(product *catalog.ShopProduct) error {
	query := s.query().Where(s.preds)
	producttDB := convert.ShopProductDB(product)
	err := query.ShouldUpdate(producttDB)
	return err
}

func (s *ShopProductStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	_deleted, err := query.Table("shop_product").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func checkProductOrVariantError(e error, code string) error {
	if e != nil {
		errMsg := e.Error()
		switch {
		case strings.Contains(errMsg, "shop_product_shop_id_code_idx"):
			e = cm.Errorf(cm.FailedPrecondition, nil, "Mã sản phẩm %v đã tồn tại. Vui lòng chọn mã khác.", code)
		case strings.Contains(errMsg, "shop_variant_shop_id_code_idx"):
			e = cm.Errorf(cm.FailedPrecondition, nil, "Mã phiên bản sản phẩm %v đã tồn tại. Vui lòng chọn mã khác.", code)
		}
	}
	return e
}

func (s *ShopProductStore) RemoveShopProductCategory() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	_deleted, err := query.Table("shop_product").UpdateMap(map[string]interface{}{
		"category_id": nil,
	})
	return int(_deleted), err
}
