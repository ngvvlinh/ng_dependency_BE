package sqlstore

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sql"
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

	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted bool
}

func (s *ProductStore) Where(pred sql.FilterQuery) *ProductStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ProductStore) ID(id int64) *ProductStore {
	s.preds = append(s.preds, s.FtProduct.ByID(id))
	return s
}

func (s *ProductStore) IDs(ids ...int64) *ProductStore {
	s.preds = append(s.preds, sql.In("id", ids))
	return s
}

func (s *ProductStore) IncludeDeleted() *ProductStore {
	s.includeDeleted = true
	return s
}

func (s *ProductStore) GetProductDB() (*catalogmodel.Product, error) {
	if !s.includeDeleted {
		s.preds = append(s.preds, s.FtProduct.NotDeleted())
	}
	var product catalogmodel.Product
	err := s.query().Where(s.preds).ShouldGet(&product)
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
	if !s.includeDeleted {
		s.preds = append(s.preds, s.FtProduct.NotDeleted())
	}
	var product catalogmodel.Product
	if err := s.query().Where(s.preds...).ShouldGet(&product); err != nil {
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

func (s *ProductStore) GetProductsDB() ([]*catalogmodel.Product, error) {
	if !s.includeDeleted {
		s.preds = append(s.preds, "deleted_at IS NULL")
	}
	var products catalogmodel.Products
	err := s.query().Where(s.preds).Find(&products)
	return products, err
}

func (s *ProductStore) GetProducts() ([]*catalog.Product, error) {
	products, err := s.GetProductsDB()
	if err != nil {
		return nil, err
	}
	return convert.Products(products), nil
}

func (s *ProductStore) GetProductsWithVariantsDB() ([]*catalogmodel.ProductFtVariant, error) {
	if !s.includeDeleted {
		s.preds = append(s.preds, s.FtProduct.NotDeleted())
	}
	var products catalogmodel.Products
	if err := s.query().Where(s.preds...).Find(&products); err != nil {
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

func (s *ProductStore) GetProductsWithVariants() ([]*catalog.ProductWithVariants, error) {
	products, err := s.GetProductsWithVariantsDB()
	if err != nil {
		return nil, err
	}
	return convert.ProductsWithVariants(products), nil
}
