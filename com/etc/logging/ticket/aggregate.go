package ticket

import (
	"context"

	"o.o/backend/com/etc/logging/ticket/model"
	"o.o/backend/com/etc/logging/ticket/sqlstore"
	com "o.o/backend/com/main"
)

type Aggregate struct {
	store sqlstore.TicketLogStoreFactory
}

func NewAggregate(db com.LogDB) *Aggregate {
	return &Aggregate{store: sqlstore.NewTicketLogStore(db)}
}

func (a *Aggregate) CreateTicketWebhookLog(ctx context.Context, args *model.TicketProviderWebhook) error {
	return a.store(ctx).CreateTicketLog(args)
}
