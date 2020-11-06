package query

import (
	"context"

	"o.o/api/main/contact"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/contact/sqlstore"
	"o.o/backend/pkg/common/bus"
)

var _ contact.QueryService = &ContactQuery{}

type ContactQuery struct {
	store sqlstore.ContactStoreFactory
}

func NewContactQuery(db com.MainDB) *ContactQuery {
	return &ContactQuery{
		store: sqlstore.NewContactStore(db),
	}
}

func ContactQueryMessageBus(q *ContactQuery) contact.QueryBus {
	b := bus.New()
	return contact.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *ContactQuery) GetContactByID(
	ctx context.Context, args *contact.GetContactByIDArgs,
) (*contact.Contact, error) {
	return q.store(ctx).ID(args.ID).ShopID(args.ShopID).GetContact()
}
