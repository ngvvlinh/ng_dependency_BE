package sqlstore

import (
	"context"

	"etop.vn/api/top/types/etc/status5"
	"etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
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

func (s *FulfillmentStore) OrderID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *FulfillmentStore) GetFfmDB() (*model.Fulfillment, error) {
	var ffm model.Fulfillment
	err := s.query().Where(s.preds...).ShouldGet(&ffm)
	return &ffm, err
}

func (s *FulfillmentStore) ListFfmsDB() ([]*model.Fulfillment, error) {
	var ffms model.Fulfillments
	err := s.query().Where(s.preds...).Find(&ffms)
	return ffms, err
}

func (s *FulfillmentStore) CreateFulfillmentDB(ctx context.Context, ffm *model.Fulfillment) (*model.Fulfillment, error) {
	if ffm.ID == 0 {
		ffm.ID = cm.NewID()
	}
	if err := ffm.BeforeInsert(); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(ffm); err != nil {
		return nil, err
	}
	return s.ID(ffm.ID).GetFfmDB()
}

func (s *FulfillmentStore) CreateFulfillmentsDB(ctx context.Context, ffms []*model.Fulfillment) error {
	for _, ffm := range ffms {
		if err := ffm.BeforeInsert(); err != nil {
			return err
		}
		if err := s.query().ShouldInsert(ffm); err != nil {
			return err
		}
	}
	return nil
}

func (s *FulfillmentStore) UpdateFulfillmentDB(ctx context.Context, ffm *model.Fulfillment) (*model.Fulfillment, error) {
	if err := s.query().Where(s.ft.ByID(ffm.ID)).Where("status not in (?, ?, ?)", status5.N, status5.NS, status5.P).ShouldUpdate(ffm); err != nil {
		return nil, err
	}
	return s.ID(ffm.ID).GetFfmDB()
}

func (s *FulfillmentStore) UpdateFulfillmentsDB(ctx context.Context, ffms []*model.Fulfillment) error {
	for _, ffm := range ffms {
		if err := ffm.BeforeInsert(); err != nil {
			return err
		}
		if err := s.query().Where(s.ft.ByID(ffm.ID)).Where("status not in (?, ?, ?)", status5.N, status5.NS, status5.P).ShouldUpdate(ffm); err != nil {
			return err
		}
	}
	return nil
}
