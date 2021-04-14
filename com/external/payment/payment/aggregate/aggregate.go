package aggregate

import (
	"context"

	"o.o/api/external/payment"
	"o.o/backend/com/external/payment/payment/sqlstore"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
)

var _ payment.Aggregate = &Aggregate{}

type Aggregate struct {
	db       *cmsql.Database
	eventBus capi.EventBus
	store    sqlstore.PaymentStoreFactory
}

func NewAggregate(db com.MainDB, eventBus capi.EventBus) *Aggregate {
	return &Aggregate{
		db:       db,
		eventBus: eventBus,
		store:    sqlstore.NewPaymentStore(db),
	}
}

func AggregateMessageBus(a *Aggregate) payment.CommandBus {
	b := bus.New()
	return payment.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreatePayment(ctx context.Context, args *payment.CreatePaymentArgs) (*payment.Payment, error) {
	cmd := &sqlstore.CreatePaymentArgs{
		ShopID:          args.ShopID,
		Amount:          args.Amount,
		Status:          args.Status,
		State:           args.State,
		PaymentProvider: args.PaymentProvider,
		ExternalTransID: args.ExternalTransID,
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
	paymentUpdated, err := a.store(ctx).UpdateExternalPaymentInfo(update)
	if err != nil {
		return nil, err
	}

	return paymentUpdated, nil
}

func (a *Aggregate) UpdateExternalPaymentInfo(ctx context.Context, args *payment.UpdateExternalPaymentInfoArgs) (res *payment.Payment, err error) {
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		cmd := &sqlstore.UpdateExternalPaymentInfoArgs{
			ID:              args.ID,
			Amount:          args.Amount,
			Status:          args.Status,
			State:           args.State,
			ExternalData:    args.ExternalData,
			ExternalTransID: args.ExternalTransID,
		}
		res, err = a.store(ctx).UpdateExternalPaymentInfo(cmd)
		if err != nil {
			return err
		}

		paymentStatusUpdatedEvent := &payment.PaymentStatusUpdatedEvent{
			ID:            args.ID,
			PaymentStatus: args.Status,
		}
		if err = a.eventBus.Publish(ctx, paymentStatusUpdatedEvent); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
