package aggregate

import (
	"context"

	"etop.vn/backend/com/etc/log/payment/model"
	"etop.vn/backend/com/etc/log/payment/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
)

type Aggregate struct {
	store sqlstore.PaymentLogStoreFactory
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{store: sqlstore.NewPaymentLogStore(db)}
}

func (a *Aggregate) CreatePaymentLog(ctx context.Context, args *model.Payment) error {
	return a.store(ctx).CreatePaymentLog(args)
}

func (a *Aggregate) UpdatePaymentLog(ctx context.Context, args *model.Payment) error {
	return a.store(ctx).UpdatePaymentLog(args)
}
