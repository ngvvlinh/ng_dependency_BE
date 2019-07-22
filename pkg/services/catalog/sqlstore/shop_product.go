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
	s.preds = append(s.preds, sq.In("ed_code", codes))
	return s
}

func (s *ShopProductStore) ByNameNormUas(names ...string) *ShopProductStore {
	s.preds = append(s.preds, sq.In("name_norm_ua", names))
	return s
}

func (s *ShopProductStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	return query.Count((*catalogmodel.ShopProduct)(nil))
}

func (s *ShopProductStore) GetShopProductDB() (*catalogmodel.ShopProduct, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())

	var product catalogmodel.ShopProduct
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

func (s *ShopProductStore) GetShopProductWithVariantsDB() (*catalogmodel.ShopProductWithVariants, error) {
	product, err := s.GetShopProductDB()
	if err != nil {
		return nil, err
	}

	var variants catalogmodel.ShopVariants
	{
		q := s.query().OrderBy("variant_id").
			Where(s.ftShopVariant.ByProductID(product.ProductID))
		q = s.includeDeleted.Check(q, s.ftShopVariant.NotDeleted())
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}
	return &catalogmodel.ShopProductWithVariants{
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

func (s *ShopProductStore) ListShopProductsDB() ([]*catalogmodel.ShopProduct, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortShopProduct)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProductWhitelist)
	if err != nil {
		return nil, err
	}

	var products catalogmodel.ShopProducts
	err = query.Find(&products)
	return products, err
}

func (s *ShopProductStore) ListShopProducts() ([]*catalog.ShopProduct, error) {
	products, err := s.ListShopProductsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProducts(products), nil
}

func (s *ShopProductStore) ListShopProductsWithVariantsDB() ([]*catalogmodel.ShopProductWithVariants, error) {
	products, err := s.ListShopProductsDB()
	if err != nil {
		return nil, err
	}

	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.ProductID
	}

	var variants catalogmodel.ShopVariants
	{
		q := s.query().In("product_id", productIDs)
		q = s.includeDeleted.Check(q, s.ftShopVariant.NotDeleted())
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}

	mapProducts := make(map[int64]*catalogmodel.ShopProductWithVariants)
	result := make([]*catalogmodel.ShopProductWithVariants, len(products))
	for i, p := range products {
		result[i] = &catalogmodel.ShopProductWithVariants{
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
