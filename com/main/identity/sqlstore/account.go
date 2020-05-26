package sqlstore

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type AccountStoreFactory func(context.Context) *AccountStore

func NewAccountStore(db *cmsql.Database) AccountStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AccountStore {
		return &AccountStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type AccountStore struct {
	query     cmsql.QueryFactory
	preds     []interface{}
	accountFt AccountFilters
	sqlstore.Paging
	filter         meta.Filters
	ctx            context.Context
	includeDeleted sqlstore.IncludeDeleted
}

func (s *AccountStore) WithPaging(paging meta.Paging) *AccountStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *AccountStore) Filters(filters meta.Filters) *AccountStore {
	if s.filter == nil {
		s.filter = filters
	} else {
		s.filter = append(s.filter, filters...)
	}
	return s
}

func (s *AccountStore) ByType(ty account_type.AccountType) *AccountStore {
	s.preds = append(s.preds, s.accountFt.ByType(ty))
	return s
}

func (s *AccountStore) ByIDs(ids ...dot.ID) *AccountStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *AccountStore) ListAccountDBs() ([]*identitymodel.Account, error) {
	var accounts identitymodel.Accounts
	err := s.query().Where(s.preds).Find(&accounts)
	return accounts, err
}
