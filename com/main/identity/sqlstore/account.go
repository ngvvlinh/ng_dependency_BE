package sqlstore

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	"o.o/backend/com/main/identity/convert"
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

func (s *AccountStore) ByID(id dot.ID) *AccountStore {
	s.preds = append(s.preds, s.accountFt.ByID(id))
	return s
}

func (s *AccountStore) ListAccountDBs() ([]*identitymodel.Account, error) {
	var accounts identitymodel.Accounts
	err := s.query().Where(s.preds).Find(&accounts)
	return accounts, err
}

func (s *AccountStore) GetAccountDB() (*identitymodel.Account, error) {
	var shop identitymodel.Account
	query := s.query().Where(s.preds)

	err := query.ShouldGet(&shop)
	return &shop, err
}

func (s *AccountStore) GetAccount() (*identity.Account, error) {
	shop, err := s.GetAccountDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_identitymodel_Account_identity_Account(shop, nil), nil
}
