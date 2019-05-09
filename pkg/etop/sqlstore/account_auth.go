package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/etop/model"
)

type AccountAuthStore struct {
	ctx   context.Context
	ft    AccountAuthFilters
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

func (s *AccountAuthStore) AccountID(id int64) *AccountAuthStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AccountAuthStore) IncludeDeleted() *AccountAuthStore {
	s.includeDeleted = true
	return s
}

func (s *AccountAuthStore) Get() (*model.AccountAuth, error) {
	var item model.AccountAuth
	err := x.Where(s.preds...).Where(s.filterDeleted(s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *AccountAuthStore) Create(m *model.AccountAuth) error {
	if err := m.BeforeInsert(); err != nil {
		return err
	}
	return x.ShouldInsert(m)
}
