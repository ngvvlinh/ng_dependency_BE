package aggregate

import (
	"context"

	"o.o/api/external/payment"
	"o.o/backend/com/external/payment/payment/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
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

func AggregateMessageBus(a *Aggregate) payment.CommandBus {
	b := bus.New()
	return payment.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreatePayment(ctx context.Context, args *payment.CreatePaymentArgs) (*payment.Payment, error) {
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

func (a *Aggregate) CreateOrUpdatePayment(ctx context.Context, args *payment.CreatePaymentArgs) (*payment.Payment, error) {
	_payment, err := a.store(ctx).OptionalExternalTransactionID(args.ExternalTransID).PaymentProvider(args.PaymentProvider).GetPayment()
	if cm.ErrorCode(err) == cm.NotFound {
		// create
		return a.CreatePayment(ctx, args)
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
