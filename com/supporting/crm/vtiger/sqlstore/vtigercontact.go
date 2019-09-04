package sqlstore

import (
	"context"

	model2 "etop.vn/backend/com/supporting/crm/vtiger/model"

	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/backend/pkg/common/validate"
)

type VtigerContactStoreFactory func(context.Context) *VtigerContactStore

type VtigerContactStore struct {
	query   func() cmsql.QueryInterface
	preds   []interface{}
	ft      VtigerContactFilters
	paging  meta.Paging
	OrderBy string
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

func (v *VtigerContactStore) SortBy(order string) *VtigerContactStore {
	v.OrderBy = order
	return v
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

func (v *VtigerContactStore) ByEtopUserID(id int64) *VtigerContactStore {
	v.preds = append(v.preds, v.ft.ByEtopUserID(id))
	return v
}

func (v *VtigerContactStore) GetVtigerContact() (*model2.VtigerContact, error) {
	query := v.query().Where(v.preds)
	var contact model2.VtigerContact
	err := query.ShouldGet(&contact)
	return &contact, err
}

func (v *VtigerContactStore) CreateVtigerContact(contact *model2.VtigerContact) error {
	if contact.EtopUserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "missing EtopUserID")
	}
	err := v.query().ShouldInsert(contact)
	return err
}

func (v *VtigerContactStore) GetContact() (*model2.VtigerContact, error) {
	var item model2.VtigerContact
	err := v.query().Where(v.preds).ShouldGet(&item)
	return &item, err
}

var SortVtigerContact = map[string]string{
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (v *VtigerContactStore) SearchContact(value string) ([]*model2.VtigerContact, error) {
	query := v.query().Where(`search_norm @@ ?::tsquery`, validate.NormalizeSearchQueryAnd(value))
	if v.OrderBy != "" {
		query = query.OrderBy(v.OrderBy)
	}
	query, err := sqlstore.LimitSort(query, &v.paging, SortVtigerContact)
	if err != nil {
		return nil, err
	}
	var contacts []*model2.VtigerContact
	err = query.Find((*model2.VtigerContacts)(&contacts))
	return contacts, err
}

func (v *VtigerContactStore) UpdateVtigerContact(contact *model2.VtigerContact) error {
	query := v.query()
	query = query.Where(v.ft.ByEtopUserID(contact.EtopUserID))
	err := v.query().Where(v.ft.ByEtopUserID(contact.EtopUserID)).ShouldUpdate(contact)
	return err
}

func (v *VtigerContactStore) GetContacts() ([]*model2.VtigerContact, error) {
	query := v.query().Where(v.preds)
	if v.OrderBy != "" {
		query = query.OrderBy(v.OrderBy)
	}
	query, err := sqlstore.LimitSort(query, &v.paging, nil)
	if err != nil {
		return nil, err
	}
	var contacts []*model2.VtigerContact
	err = query.Find((*model2.VtigerContacts)(&contacts))
	return contacts, err
}
