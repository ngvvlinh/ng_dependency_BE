package sqlstore

import (
	"context"

	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/etop/model"
)

type FulfillmentStore struct {
	ctx   context.Context
	ft    FulfillmentFilters
	preds []interface{}
}

func Fulfillment(ctx context.Context) *FulfillmentStore {
	return &FulfillmentStore{ctx: ctx}
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

func (s *FulfillmentStore) Get() (*model.Fulfillment, error) {
	var ffm model.Fulfillment
	err := x.Where(s.preds...).ShouldGet(&ffm)
	return &ffm, err
}

func (s *FulfillmentStore) Insert(ffm *model.Fulfillment) error {
	return x.ShouldInsert(ffm)
}
