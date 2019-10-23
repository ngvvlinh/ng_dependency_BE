package sqlstore

import (
	"context"

	"etop.vn/backend/com/etc/log/payment/model"

	"etop.vn/backend/pkg/common/cmsql"
)

type PaymentLogStoreFactory func(context.Context) *PaymentLogStore

func NewPaymentLogStore(db *cmsql.Database) PaymentLogStoreFactory {
	return func(ctx context.Context) *PaymentLogStore {
		return &PaymentLogStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type PaymentLogStore struct {
	query func() cmsql.QueryInterface
	ft    PaymentFilters
}

func (s *PaymentLogStore) CreatePaymentLog(payment *model.Payment) error {
	return s.query().ShouldInsert(payment)
}

func (s *PaymentLogStore) UpdatePaymentLog(payment *model.Payment) error {
	return s.query().Where(s.ft.ByID(payment.ID)).ShouldUpdate(payment)
}
