package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/api/services/affiliate"
	"o.o/backend/com/services/affiliate/convert"
	"o.o/backend/com/services/affiliate/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ProductPromotionStoreFactory func(ctx context.Context) *ProductPromotionStore

func NewProductPromotionStore(db *cmsql.Database) ProductPromotionStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ProductPromotionStore {
		return &ProductPromotionStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type ProductPromotionStore struct {
	query cmsql.QueryFactory
	preds []interface{}

	ft ProductPromotionFilters

	sqlstore.Paging
	filters meta.Filters
}

func (s *ProductPromotionStore) ID(id dot.ID) *ProductPromotionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ProductPromotionStore) ShopID(id dot.ID) *ProductPromotionStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ProductPromotionStore) ProductID(id dot.ID) *ProductPromotionStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *ProductPromotionStore) ProductIDs(ids ...dot.ID) *ProductPromotionStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *ProductPromotionStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ProductPromotion)(nil))
}

func (s *ProductPromotionStore) WithPaging(paging meta.Paging) *ProductPromotionStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *ProductPromotionStore) Filters(filters meta.Filters) *ProductPromotionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ProductPromotionStore) GetProductPromotionDB() (*model.ProductPromotion, error) {
	var productPromotion model.ProductPromotion
	err := s.query().Where(s.preds).ShouldGet(&productPromotion)
	return &productPromotion, err
}

func (s *ProductPromotionStore) GetProductPromotion() (*affiliate.ProductPromotion, error) {
	productPromotion, err := s.GetProductPromotionDB()
	if err != nil {
		return nil, err
	}
	return convert.ProductPromotion(productPromotion), nil
}

func (s *ProductPromotionStore) GetProductPromotions() ([]*affiliate.ProductPromotion, error) {
	var results model.ProductPromotions
	err := s.query().Where(s.preds).Find(&results)
	return convert.ProductPromotions(results), err

}

func (s *ProductPromotionStore) CreateProductPromotion(productPromotion *model.ProductPromotion) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.query().Insert(productPromotion)
	return err
}

func (s *ProductPromotionStore) UpdateProductPromotion(productPromotion *model.ProductPromotion) error {
	sqlstore.MustNoPreds(s.preds)
	_, err := s.ID(productPromotion.ID).query().Where(s.preds).Update(productPromotion)
	return err
}
