package shippingwebhook

import (
	"context"

	"o.o/backend/com/etc/logging/shippingwebhook/model"
	"o.o/backend/com/etc/logging/shippingwebhook/sqlstore"
	com "o.o/backend/com/main"
)

type Aggregate struct {
	store sqlstore.ShippingWebhookStoreFactory
}

func NewAggregate(db com.LogDB) *Aggregate {
	return &Aggregate{store: sqlstore.NewShippingWebhookStore(db)}
}

func (a *Aggregate) CreateShippingWebhookLog(ctx context.Context, args *model.ShippingProviderWebhook) error {
	return a.store(ctx).CreateShippingWebhookLog(args)
}
