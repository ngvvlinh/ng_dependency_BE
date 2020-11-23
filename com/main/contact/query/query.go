package query

import (
	"context"

	"o.o/api/main/contact"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/contact/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/filter"
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

func (q *ContactQuery) GetContacts(
	ctx context.Context, args *contact.GetContactsArgs,
) (*contact.GetContactsResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "shop_id is missing")
	}
	query := q.store(ctx).ShopID(args.ShopID).WithPaging(args.Paging)
	if len(args.IDs) != 0 {
		query = query.IDs(args.IDs...)
	}
	if args.Phone != "" {
		query = query.FullTextSearchFullPhone(filter.FullTextSearch(args.Phone))
	}

	contacts, err := query.ListContacts()
	if err != nil {
		return nil, err
	}

	return &contact.GetContactsResponse{
		Contacts: contacts,
		Paging:   query.GetPaging(),
	}, nil
}
