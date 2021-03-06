package aggregate

import (
	"context"

	"o.o/backend/com/etc/logging/payment/model"
	"o.o/backend/com/etc/logging/payment/sqlstore"
	com "o.o/backend/com/main"
)

type Aggregate struct {
	store sqlstore.PaymentLogStoreFactory
}

func New(db com.LogDB) *Aggregate {
	return &Aggregate{store: sqlstore.NewPaymentLogStore(db)}
}

func (a *Aggregate) CreatePaymentLog(ctx context.Context, args *model.Payment) error {
	return a.store(ctx).CreatePaymentLog(args)
}

func (a *Aggregate) UpdatePaymentLog(ctx context.Context, args *model.Payment) error {
	return a.store(ctx).UpdatePaymentLog(args)
}
