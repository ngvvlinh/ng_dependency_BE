package query

import (
	"context"

	"o.o/api/external/payment"
	"o.o/backend/com/external/payment/payment/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ payment.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.PaymentStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewPaymentStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) payment.QueryBus {
	b := bus.New()
	return payment.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetPaymentByID(ctx context.Context, id dot.ID) (*payment.Payment, error) {
	return q.store(ctx).ID(id).GetPayment()
}

func (q *QueryService) GetPaymentByExternalTransID(ctx context.Context, id string) (*payment.Payment, error) {
	return q.store(ctx).ExternalTransactionID(id).GetPayment()
}
