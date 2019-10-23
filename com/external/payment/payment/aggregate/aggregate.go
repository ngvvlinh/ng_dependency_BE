package aggregate

import (
	"context"

	"etop.vn/api/external/payment"
	"etop.vn/backend/com/external/payment/payment/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ payment.Aggregate = &Aggregate{}

type Aggregate struct {
	db    cmsql.Transactioner
	store sqlstore.PaymentStoreFactory
}

func NewAggregate(db *cmsql.Database) *Aggregate {
	return &Aggregate{
		db:    db,
		store: sqlstore.NewPaymentStore(db),
	}
}

func (a *Aggregate) MessageBus() payment.CommandBus {
	b := bus.New()
	return payment.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateOrUpdatePayment(ctx context.Context, args *payment.CreatePaymentArgs) (*payment.Payment, error) {
	_payment, err := a.store(ctx).ExternalTransactionID(args.ExternalTransID).PaymentProvider(string(args.PaymentProvider)).GetPayment()
	if cm.ErrorCode(err) == cm.NotFound {
		// create
		cmd := &sqlstore.CreatePaymentArgs{
			Amount:          args.Amount,
			Status:          args.Status,
			State:           args.State,
			PaymentProvider: args.PaymentProvider,
			ExternalTransID: args.ExternalTransID,
			ExternalData:    args.ExternalData,
		}
		return a.store(ctx).CreatePayment(cmd)
	}
	if err != nil {
		return nil, err
	}
	// update
	update := &sqlstore.UpdateExternalPaymentInfoArgs{
		ID:           _payment.ID,
		Amount:       args.Amount,
		Status:       args.Status,
		State:        args.State,
		ExternalData: args.ExternalData,
	}
	return a.store(ctx).UpdateExternalPaymentInfo(update)
}

func (a *Aggregate) UpdateExternalPaymentInfo(ctx context.Context, args *payment.UpdateExternalPaymentInfoArgs) (*payment.Payment, error) {
	cmd := &sqlstore.UpdateExternalPaymentInfoArgs{
		ID:              args.ID,
		Amount:          args.Amount,
		Status:          args.Status,
		State:           args.State,
		ExternalData:    args.ExternalData,
		ExternalTransID: args.ExternalTransID,
	}
	return a.store(ctx).UpdateExternalPaymentInfo(cmd)
}
