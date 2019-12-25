package query

import (
	"context"

	"etop.vn/api/external/payment"
	"etop.vn/backend/com/external/payment/payment/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
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

func (q *QueryService) MessageBus() payment.QueryBus {
	b := bus.New()
	return payment.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetPaymentByID(ctx context.Context, id dot.ID) (*payment.Payment, error) {
	return q.store(ctx).ID(id).GetPayment()
}

func (q *QueryService) GetPaymentByExternalTransID(ctx context.Context, id string) (*payment.Payment, error) {
	return q.store(ctx).ExternalTransactionID(id).GetPayment()
}
