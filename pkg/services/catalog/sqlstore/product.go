package sqlstore

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/etop"
	"etop.vn/api/meta"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	etopconvert "etop.vn/backend/pkg/etop/convert"
	"etop.vn/backend/pkg/services/catalog/convert"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

type ProductStoreFactory func(context.Context) *ProductStore

func NewProductStore(db cmsql.Database) ProductStoreFactory {
	return func(ctx context.Context) *ProductStore {
		return &ProductStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ProductStore struct {
	FtProduct ProductFilters

	// unexported
	ftVariant VariantFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters

	includeDeleted sqlstore.IncludeDeleted
}

func (s *ProductStore) Where(pred sq.FilterQuery) *ProductStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ProductStore) Filters(filters meta.Filters) *ProductStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ProductStore) ID(id int64) *ProductStore {
	s.preds = append(s.preds, s.FtProduct.ByID(id))
	return s
}

func (s *ProductStore) IDs(ids ...int64) *ProductStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ProductStore) Status(status etop.Status3) *ProductStore {
	s.preds = append(s.preds,
		s.FtProduct.ByStatus(etopconvert.Status3ToModel(status)))
	return s
}

func (s *ProductStore) ProductSourceID(id int64) *ProductStore {
	s.preds = append(s.preds, s.FtProduct.ByProductSourceID(id))
	return s
}

func (s *ProductStore) Code(code string) *ProductStore {
	s.preds = append(s.preds, s.FtProduct.ByCode(code))
	return s
}

func (s *ProductStore) Codes(codes ...string) *ProductStore {
	s.preds = append(s.preds, sq.In("ed_code", codes))
	return s
}

func (s *ProductStore) ByNameNormUas(names ...string) *ProductStore {
	s.preds = append(s.preds, sq.In("name_norm_ua", names))
	return s
}

func (s *ProductStore) IncludeDeleted() *ProductStore {
	s.includeDeleted = true
	return s
}

func (s *ProductStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	return query.Count((*catalogmodel.Product)(nil))
}

func (s *ProductStore) GetProductDB() (*catalogmodel.Product, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())

	var product catalogmodel.Product
	err := query.ShouldGet(&product)
	return &product, err
}

func (s *ProductStore) GetProduct() (*catalog.Product, error) {
	product, err := s.GetProductDB()
	if err != nil {
		return nil, err
	}
	return convert.Product(product), nil
}

func (s *ProductStore) GetProductWithVariantsDB() (*catalogmodel.ProductFtVariant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())

	var product catalogmodel.Product
	if err := query.ShouldGet(&product); err != nil {
		return nil, err
	}

	var variants catalogmodel.Variants
	{
		q := s.query().Where(s.ftVariant.ByProductID(product.ID))
		if !s.includeDeleted {
			q = q.Where(s.ftVariant.NotDeleted())
		}
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}
	return &catalogmodel.ProductFtVariant{
		Product:  &product,
		Variants: variants,
	}, nil
}

func (s *ProductStore) GetProductWithVariants() (*catalog.ProductWithVariants, error) {
	product, err := s.GetProductWithVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.ProductWithVariants(product), nil
}

func (s *ProductStore) ListProductsDB(paging meta.Paging) ([]*catalogmodel.Product, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &paging, SortProduct)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterProductWhitelist)
	if err != nil {
		return nil, err
	}

	var products catalogmodel.Products
	err = query.Find(&products)
	return products, err
}

func (s *ProductStore) ListProducts(paging meta.Paging) ([]*catalog.Product, error) {
	products, err := s.ListProductsDB(paging)
	if err != nil {
		return nil, err
	}
	return convert.Products(products), nil
}

func (s *ProductStore) ListProductsWithVariantsDB(paging meta.Paging) ([]*catalogmodel.ProductFtVariant, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.FtProduct.NotDeleted())
	query, err := sqlstore.LimitSort(query, &paging, SortProduct)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterProductWhitelist)
	if err != nil {
		return nil, err
	}

	var products catalogmodel.Products
	if err := query.Find(&products); err != nil {
		return nil, err
	}
	productIDs := make([]int64, len(products))
	for i, p := range products {
		productIDs[i] = p.ID
	}

	var variants catalogmodel.Variants
	{
		q := s.query().In("product_id", productIDs)
		if !s.includeDeleted {
			q = q.Where(s.ftVariant.NotDeleted())
		}
		if err := q.Find(&variants); err != nil {
			return nil, err
		}
	}

	mapProducts := make(map[int64]*catalogmodel.ProductFtVariant)
	result := make([]*catalogmodel.ProductFtVariant, len(products))
	for i, p := range products {
		result[i] = &catalogmodel.ProductFtVariant{Product: p}
		mapProducts[p.ID] = result[i]
	}
	for _, v := range variants {
		p := mapProducts[v.ProductID]
		if p != nil {
			p.Variants = append(p.Variants, v)
		}
	}
	return result, nil
}

func (s *ProductStore) ListProductsWithVariants(paging meta.Paging) ([]*catalog.ProductWithVariants, error) {
	products, err := s.ListProductsWithVariantsDB(paging)
	if err != nil {
		return nil, err
	}
	return convert.ProductsWithVariants(products), nil
}
