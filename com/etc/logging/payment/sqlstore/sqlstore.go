package sqlstore

import (
	"context"

	"o.o/backend/com/etc/logging/payment/model"
	"o.o/backend/pkg/common/sql/cmsql"
)

type PaymentLogStoreFactory func(context.Context) *PaymentLogStore

func NewPaymentLogStore(db *cmsql.Database) PaymentLogStoreFactory {
	return func(ctx context.Context) *PaymentLogStore {
		return &PaymentLogStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type PaymentLogStore struct {
	query cmsql.QueryFactory
	ft    PaymentFilters
}

func (s *PaymentLogStore) CreatePaymentLog(payment *model.Payment) error {
	return s.query().ShouldInsert(payment)
}

func (s *PaymentLogStore) UpdatePaymentLog(payment *model.Payment) error {
	return s.query().Where(s.ft.ByID(payment.ID)).ShouldUpdate(payment)
}
