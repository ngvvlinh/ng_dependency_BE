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
	FtProduct     ProductFilters
	FtShopProduct ShopProductFilters

	// unexported
	ftVariant     VariantFilters
	ftShopVariant ShopVariantFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ShopProductStore) extend() *ShopProductStore {
	s.FtProduct.prefix = "p"
	s.FtShopProduct.prefix = "sp"
	s.ftVariant.prefix = "v"
	s.ftShopVariant.prefix = "sv"
	return s
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
	s.preds = append(s.preds, s.FtProduct.ByID(id))
	return s
}

func (s *ShopProductStore) IDs(ids ...int64) *ShopProductStore {
	s.preds = append(s.preds, sq.In("p.id", ids))
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

func (s *ShopProductStore) Count() (uint64, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	return query.Count((*catalogmodel.ShopProductExtendeds)(nil))
}

func (s *ShopProductStore) GetShopProductDB() (*catalogmodel.ShopProductExtended, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())

	var product catalogmodel.ShopProductExtended
	err := query.ShouldGet(&product)
	return &product, err
}

func (s *ShopProductStore) GetShopProduct() (*catalog.ShopProductExtended, error) {
	product, err := s.GetShopProductDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProductExtended(product), nil
}

func (s *ShopProductStore) GetShopProductWithVariantsDB() (*catalogmodel.ShopProductFtVariant, error) {
	product, err := s.GetShopProductDB()
	if err != nil {
		return nil, err
	}

	var variants catalogmodel.ShopVariantExtendeds
	{
		q := s.extend().query().OrderBy("v.id").
			Where(s.ftVariant.ByProductID(product.ID))
		q = s.includeDeleted.Check(q, s.ftVariant.NotDeleted())
		q = s.includeDeleted.Check(q, s.ftShopVariant.NotDeleted())
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}
	return &catalogmodel.ShopProductFtVariant{
		ShopProduct: product.ShopProduct,
		Product:     product.Product,
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

func (s *ShopProductStore) ListShopProductsDB() ([]*catalogmodel.ShopProductExtended, error) {
	query := s.extend().query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	query = s.includeDeleted.Check(query, s.FtShopProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortShopProduct)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterShopProductWhitelist)
	if err != nil {
		return nil, err
	}

	var products catalogmodel.ShopProductExtendeds
	err = query.Find(&products)
	return products, err
}

func (s *ShopProductStore) ListShopProducts() ([]*catalog.ShopProductExtended, error) {
	products, err := s.ListShopProductsDB()
	if err != nil {
		return nil, err
	}
	return convert.ShopProductExtendeds(products), nil
}

func (s *ShopProductStore) ListShopProductsWithVariantsDB() ([]*catalogmodel.ShopProductFtVariant, error) {
	products, err := s.ListShopProductsDB()
	if err != nil {
		return nil, err
	}

	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}

	var variants catalogmodel.ShopVariantExtendeds
	{
		q := s.extend().query().In("v.product_id", productIDs)
		q = s.includeDeleted.Check(q, s.ftVariant.NotDeleted())
		q = s.includeDeleted.Check(q, s.ftShopVariant.NotDeleted())
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}

	mapProducts := make(map[int64]*catalogmodel.ShopProductFtVariant)
	result := make([]*catalogmodel.ShopProductFtVariant, len(products))
	for i, p := range products {
		result[i] = &catalogmodel.ShopProductFtVariant{
			ShopProduct: p.ShopProduct,
			Product:     p.Product,
		}
		mapProducts[p.ID] = result[i]
	}
	for _, v := range variants {
		p := mapProducts[v.Variant.ProductID]
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
