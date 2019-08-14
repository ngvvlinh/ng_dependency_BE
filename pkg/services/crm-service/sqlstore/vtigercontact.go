package sqlstore

import (
	"context"

	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/services/crm-service/model"
)

type VtigerContactStoreFactory func(context.Context) *VtigerContactStore

type VtigerContactStore struct {
	query  func() cmsql.QueryInterface
	preds  []interface{}
	ft     VtigerContactFilters
	paging meta.Paging
}

func NewVtigerStore(db cmsql.Database) VtigerContactStoreFactory {
	return func(ctx context.Context) *VtigerContactStore {
		return &VtigerContactStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *VtigerContactStore) Paging(paging meta.Paging) *VtigerContactStore {
	s.paging = paging
	return s
}

func (s *VtigerContactStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (v *VtigerContactStore) ByEmail(email string) *VtigerContactStore {
	v.preds = append(v.preds, v.ft.ByEmail(email))
	return v
}

func (v *VtigerContactStore) ByPhone(phone string) *VtigerContactStore {
	v.preds = append(v.preds, v.ft.ByPhone(phone))
	return v
}

func (v *VtigerContactStore) ByEtopID(id int64) *VtigerContactStore {
	v.preds = append(v.preds, v.ft.ByEtopID(id))
	return v
}

func (v *VtigerContactStore) GetVtigerContact() (*model.VtigerContact, error) {
	query := v.query().Where(v.preds)
	var contact model.VtigerContact
	err := query.ShouldGet(&contact)
	return &contact, err
}

func (v *VtigerContactStore) CreateVtigerContact(contact *model.VtigerContact) error {
	if contact.EtopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing EtopID")
	}
	err := v.query().ShouldInsert(contact)
	return err
}

func (v *VtigerContactStore) GetContact() (*model.VtigerContact, error) {
	var item model.VtigerContact
	err := v.query().Where(v.preds).ShouldGet(&item)
	return &item, err
}

var SortVtigerContact = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (v *VtigerContactStore) SearchContact(value string) ([]*model.VtigerContact, error) {
	query := v.query().Where(`search_norm @@ ?::tsquery`, validate.NormalizeSearchQueryAnd(value))
	query, err := sqlstore.LimitSort(query, &v.paging, SortVtigerContact)
	if err != nil {
		return nil, err
	}
	var contacts []*model.VtigerContact
	err = query.Find((*model.VtigerContacts)(&contacts))
	return contacts, err
}

func (v *VtigerContactStore) UpdateVtigerContact(contact *model.VtigerContact) error {
	query := v.query()
	query = query.Where(v.ft.ByEtopID(contact.EtopID))
	err := v.query().Where(v.ft.ByEtopID(contact.EtopID)).ShouldUpdate(contact)
	return err
}

func (v *VtigerContactStore) GetContacts() ([]*model.VtigerContact, error) {
	query := v.query().Where(v.preds)
	query, err := sqlstore.LimitSort(query, &v.paging, SortVtigerContact)
	if err != nil {
		return nil, err
	}
	var contacts []*model.VtigerContact
	err = v.query().Find((*model.VtigerContacts)(&contacts))
	return contacts, err
}
