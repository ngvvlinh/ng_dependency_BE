package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sq"

	"etop.vn/api/meta"

	"etop.vn/backend/pkg/common/sqlstore"

	"etop.vn/backend/com/services/affiliate/convert"

	"etop.vn/api/services/affiliate"

	"etop.vn/backend/com/services/affiliate/model"

	"etop.vn/backend/pkg/common/cmsql"
)

type ProductPromotionStoreFactory func(ctx context.Context) *ProductPromotionStore

func NewProductPromotionStore(db cmsql.Database) ProductPromotionStoreFactory {
	return func(ctx context.Context) *ProductPromotionStore {
		return &ProductPromotionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ProductPromotionStore struct {
	query func() cmsql.QueryInterface
	preds []interface{}

	ft ProductPromotionFilters

	paging  meta.Paging
	filters meta.Filters
}

func (s *ProductPromotionStore) ID(id int64) *ProductPromotionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ProductPromotionStore) ShopID(id int64) *ProductPromotionStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ProductPromotionStore) ProductID(id int64) *ProductPromotionStore {
	s.preds = append(s.preds, s.ft.ByProductID(id))
	return s
}

func (s *ProductPromotionStore) ProductIDs(ids ...int64) *ProductPromotionStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *ProductPromotionStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	return query.Count((*model.ProductPromotion)(nil))
}

func (s *ProductPromotionStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *ProductPromotionStore) Paging(paging meta.Paging) *ProductPromotionStore {
	s.paging = paging
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