package sqlstore

import (
	"context"
	"time"

	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type AccountUserStoreFactory func(context.Context) *AccountUserStore

func NewAccoutnUserStore(db *cmsql.Database) AccountUserStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AccountUserStore {
		return &AccountUserStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AccountUserStore struct {
	query cmsql.QueryFactory
	preds []interface{}
	ft    AccountUserFilters
}

func (s *AccountUserStore) ByAccountID(id dot.ID) *AccountUserStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AccountUserStore) ByUserID(id dot.ID) *AccountUserStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *AccountUserStore) GetAccountUserDB() (*identitymodel.AccountUser, error) {
	var acc identitymodel.AccountUser
	err := s.query().Where(s.preds).ShouldGet(&acc)
	return &acc, err
}

type DeleteAccountUserArgs struct {
	AccountID dot.ID
	UserID    dot.ID
}

func (s *AccountUserStore) DeleteAccountUser(args DeleteAccountUserArgs) error {
	if args.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	if args.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing UserID")
	}

	update := &identitymodel.AccountUser{
		DeletedAt: time.Now(),
	}

	if err := s.query().Where(s.ft.ByUserID(args.UserID)).
		Where(s.ft.ByAccountID(args.AccountID)).
		ShouldDelete(update); err != nil {
		return err
	}
	return nil
}
