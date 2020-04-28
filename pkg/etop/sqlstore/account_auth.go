package sqlstore

import (
	"context"

	identitymodel "o.o/backend/com/main/identity/model"
	identitysqlstore "o.o/backend/com/main/identity/sqlstore"
	"o.o/capi/dot"
)

type AccountAuthStore struct {
	ctx   context.Context
	ft    identitysqlstore.AccountAuthFilters
	preds []interface{}

	includeDeleted
}

func AccountAuth(ctx context.Context) *AccountAuthStore {
	return &AccountAuthStore{ctx: ctx}
}

func (s *AccountAuthStore) Key(key string) *AccountAuthStore {
	s.preds = append(s.preds, s.ft.ByAuthKey(key))
	return s
}

func (s *AccountAuthStore) AccountID(id dot.ID) *AccountAuthStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AccountAuthStore) IncludeDeleted() *AccountAuthStore {
	s.includeDeleted = true
	return s
}

func (s *AccountAuthStore) Get() (*identitymodel.AccountAuth, error) {
	var item identitymodel.AccountAuth
	err := x.Where(s.preds...).Where(s.filterDeleted(&s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *AccountAuthStore) Create(m *identitymodel.AccountAuth) error {
	if err := m.BeforeInsert(); err != nil {
		return err
	}
	return x.ShouldInsert(m)
}
