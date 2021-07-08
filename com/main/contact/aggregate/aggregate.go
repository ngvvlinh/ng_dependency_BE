package aggregate

import (
	"context"

	"o.o/api/main/contact"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/contact/convert"
	"o.o/backend/com/main/contact/model"
	"o.o/backend/com/main/contact/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

var ll = l.New()
var _ contact.Aggregate = &ContactAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type ContactAggregate struct {
	store sqlstore.ContactStoreFactory
}

func NewContactAggregate(
	database com.MainDB,
) *ContactAggregate {
	return &ContactAggregate{
		store: sqlstore.NewContactStore(database),
	}
}

func ContactAggregateMessageBus(a *ContactAggregate) contact.CommandBus {
	b := bus.New()
	return contact.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *ContactAggregate) CreateContact(
	ctx context.Context, args *contact.CreateContactArgs,
) (*contact.Contact, error) {
	if args.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tên không được để trống")
	}
	if args.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không được để trống")
	}
	contact := new(contact.Contact)
	if err := scheme.Convert(args, contact); err != nil {
		return nil, err
	}
	if _, ok := validate.NormalizePhone(args.Phone); !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}

	contact.PhoneNorm = validate.NormalizeSearchPhone(args.Phone)

	if err := a.store(ctx).CreateContact(contact); err != nil {
		return nil, err
	}
	return contact, nil
}

func (a *ContactAggregate) UpdateContact(
	ctx context.Context, args *contact.UpdateContactArgs,
) (*contact.Contact, error) {
	contact, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetContact()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(args, contact); err != nil {
		return nil, err
	}

	contactDB := new(model.Contact)
	if err := scheme.Convert(contact, contactDB); err != nil {
		return nil, err
	}
	if args.Phone.Valid {
		if _, ok := validate.NormalizePhone(args.Phone.String); !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
		}
		contactDB.PhoneNorm = validate.NormalizeSearchPhone(args.Phone.String)
	}

	err = a.store(ctx).UpdateContactDB(contactDB)
	return contact, err
}

func (a *ContactAggregate) DeleteContact(
	ctx context.Context, args *contact.DeleteContactArgs,
) (deleted int, _ error) {
	if _, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetContactDB(); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy contact.").
			Throw()
	}

	return a.store(ctx).ID(args.ID).ShopID(args.ShopID).SoftDelete()
}
