package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/contact"
	"o.o/api/meta"
	"o.o/backend/com/main/contact/convert"
	"o.o/backend/com/main/contact/model"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ContactStoreFactory func(ctx context.Context) *ContactStore

func NewContactStore(db *cmsql.Database) ContactStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ContactStore {
		return &ContactStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type ContactStore struct {
	ft ContactFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters

	includeDeleted sqlstore.IncludeDeleted

	ctx context.Context
}

func (s *ContactStore) Filters(filters meta.Filters) *ContactStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *ContactStore) ID(id dot.ID) *ContactStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ContactStore) ShopID(id dot.ID) *ContactStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ContactStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}

func (s *ContactStore) GetContactDB() (*model.Contact, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var contact model.Contact
	err := query.ShouldGet(&contact)
	return &contact, err
}

func (s *ContactStore) GetContact() (contactResult *contact.Contact, _ error) {
	contact, err := s.GetContactDB()
	if err != nil {
		return nil, err
	}
	contactResult = convert.Convert_contactmodel_Contact_contact_Contact(contact, contactResult)
	return contactResult, nil
}

func (s *ContactStore) CreateContact(contact *contact.Contact) error {
	sqlstore.MustNoPreds(s.preds)

	contactDB := new(model.Contact)
	convert.Convert_contact_Contact_contactmodel_Contact(contact, contactDB)
	contactDB.WLPartnerID = wl.GetWLPartnerID(s.ctx)

	if _, err := s.query().Insert(contactDB); err != nil {
		return err
	}

	var tempContact model.Contact
	if err := s.query().Where(s.ft.ByID(contact.ID), s.ft.ByShopID(contact.ShopID)).ShouldGet(&tempContact); err != nil {
		return err
	}

	contact.CreatedAt = tempContact.CreatedAt
	contact.UpdatedAt = tempContact.UpdatedAt
	return nil
}

func (s *ContactStore) UpdateContactDB(contact *model.Contact) error {
	sqlstore.MustNoPreds(s.preds)
	query := s.query().Where(s.ft.ByID(contact.ID), s.ft.ByShopID(contact.ShopID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(contact)
}

func (s *ContactStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("contact").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}
