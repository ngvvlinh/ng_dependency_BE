package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	sq "etop.vn/backend/pkg/common/sql"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
)

type FulfillmentStoreFactory func(context.Context) *FulfillmentStore

func NewFulfillmentStore(db cmsql.Database) FulfillmentStoreFactory {
	return func(ctx context.Context) *FulfillmentStore {
		return &FulfillmentStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type FulfillmentStore struct {
	ft FulfillmentFilters

	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted bool
}

func (s *FulfillmentStore) ID(id int64) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FulfillmentStore) ShippingCode(code string) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code))
	return s
}

func (s *FulfillmentStore) IDOrShippingCode(id int64, shippingCode string) *FulfillmentStore {
	s.preds = append(s.preds, sq.Once{
		s.ft.ByID(id),
		s.ft.ByShippingCode(shippingCode),
	})
	return s
}

func (s *FulfillmentStore) ShopID(id int64) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *FulfillmentStore) PartnerID(id int64) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *FulfillmentStore) Get() (*shipmodel.Fulfillment, error) {
	var ffm shipmodel.Fulfillment
	err := s.query().Where(s.preds...).ShouldGet(&ffm)
	return &ffm, err
}

func (s *FulfillmentStore) Insert(ffm *shipmodel.Fulfillment) error {
	return s.query().ShouldInsert(ffm)
}