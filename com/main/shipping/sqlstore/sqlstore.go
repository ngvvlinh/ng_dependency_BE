package sqlstore

import (
	"context"

	"etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/capi/dot"
)

type FulfillmentStoreFactory func(context.Context) *FulfillmentStore

func NewFulfillmentStore(db *cmsql.Database) FulfillmentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FulfillmentStore {
		return &FulfillmentStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FulfillmentStore struct {
	ft FulfillmentFilters

	query cmsql.QueryFactory
	preds []interface{}

	includeDeleted bool
}

func (s *FulfillmentStore) ID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FulfillmentStore) ShippingCode(code string) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code))
	return s
}

func (s *FulfillmentStore) IDOrShippingCode(id dot.ID, shippingCode string) *FulfillmentStore {
	s.preds = append(s.preds, sq.Once{
		s.ft.ByID(id),
		s.ft.ByShippingCode(shippingCode),
	})
	return s
}

func (s *FulfillmentStore) ShopID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *FulfillmentStore) PartnerID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *FulfillmentStore) Get() (*model.Fulfillment, error) {
	var ffm model.Fulfillment
	err := s.query().Where(s.preds...).ShouldGet(&ffm)
	return &ffm, err
}

func (s *FulfillmentStore) Insert(ffm *model.Fulfillment) error {
	return s.query().ShouldInsert(ffm)
}
