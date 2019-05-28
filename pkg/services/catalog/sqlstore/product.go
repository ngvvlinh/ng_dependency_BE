package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sql"
)

type ProductStoreFactory func(context.Context) *ProductStore

func NewProductStore(db cmsql.Database) ProductStoreFactory {
	return func(ctx context.Context) *ProductStore {
		return &ProductStore{
			query: db.WithContext(ctx),
		}
	}
}

type ProductStore struct {
	FtProduct ProductFilters
	FtVariant VariantFilters

	query cmsql.QueryInterface
	preds []interface{}

	includeDeleted bool
}

func (s *ProductStore) Where(pred sql.FilterQuery) *ProductStore {
	s.preds = append(s.preds, pred)
	return s
}

func (s *ProductStore) ID(id int64) {
	s.preds = append(s.preds, s.FtProduct.ByID(id))
}

func (s *ProductStore) GetProduct() {

}

func (s *ProductStore) GetProductWithVariants() {

}

func (s *ProductStore) GetProducts() {

}

func (s *ProductStore) GetProductsWithVariants() {

}
