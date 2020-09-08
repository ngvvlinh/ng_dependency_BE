package sqlstore

import (
	"context"

	com "o.o/backend/com/main"
	identitymodel "o.o/backend/com/main/identity/model"
	identitysqlstore "o.o/backend/com/main/identity/sqlstore"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type AccountAuthStore struct {
	query cmsql.QueryFactory
	ft    identitysqlstore.AccountAuthFilters
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type AccountAuthStoreFactory func(ctx context.Context) *AccountAuthStore

func NewAccountAuthStore(db com.MainDB) AccountAuthStoreFactory {
	return func(ctx context.Context) *AccountAuthStore {
		s := &AccountAuthStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
		return s
	}
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
	err := s.query().Where(s.preds...).Where(s.includeDeleted.FilterDeleted(&s.ft)).ShouldGet(&item)
	return &item, err
}

func (s *AccountAuthStore) Create(m *identitymodel.AccountAuth) error {
	if err := m.BeforeInsert(); err != nil {
		return err
	}
	return s.query().ShouldInsert(m)
}
