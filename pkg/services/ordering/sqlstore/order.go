package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
)

type OrderStore struct {
	db    cmsql.Database
	ctx   context.Context
	ft    OrderFilters
	preds []interface{}
}

func NewOrderStore(db cmsql.Database) *OrderStore {
	return &OrderStore{db: db, ctx: context.Background()}
}

func (s *OrderStore) WithContext(ctx context.Context) *OrderStore {
	return &OrderStore{db: s.db, ctx: ctx}
}

func (s *OrderStore) ID(id int64) *OrderStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *OrderStore) ShopID(id int64) *OrderStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *OrderStore) PartnerID(id int64) *OrderStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *OrderStore) ExternalID(id string) *OrderStore {
	s.preds = append(s.preds, s.ft.ByExternalOrderID(id))
	return s
}

func (s *OrderStore) Code(code string) *OrderStore {
	s.preds = append(s.preds, s.ft.ByCode(code))
	return s
}

func (s *OrderStore) ExternalShopID(shopID int64, externalID string) *OrderStore {
	s.preds = append(s.preds,
		s.ft.ByShopID(shopID),
		s.ft.ByExternalOrderID(externalID),
	)
	return s
}

func (s *OrderStore) ExternalPartnerID(partnerID int64, externalID string) *OrderStore {
	s.preds = append(s.preds,
		s.ft.ByPartnerID(partnerID),
		s.ft.ByExternalOrderID(externalID),
	)
	return s
}

func (s *OrderStore) Get() (*ordermodel.Order, error) {
	var order ordermodel.Order
	err := s.db.Where(s.preds...).ShouldGet(&order)
	return &order, err
}
