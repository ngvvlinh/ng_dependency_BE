package sqlstore

import (
	"context"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
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
	ft    sqlstore.AccountUserFilters
}

func (s *AccountUserStore) ByAccountID(id int64) *AccountUserStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AccountUserStore) ByUserID(id int64) *AccountUserStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *AccountUserStore) GetAccountUserDB() (*model.AccountUser, error) {
	var acc model.AccountUser
	err := s.query().Where(s.preds).ShouldGet(&acc)
	return &acc, err
}

type DeleteAccountUserArgs struct {
	AccountID int64
	UserID    int64
}

func (s *AccountUserStore) DeleteAccountUser(args DeleteAccountUserArgs) error {
	if args.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	if args.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing UserID")
	}

	update := &model.AccountUser{
		DeletedAt: time.Now(),
	}

	if err := s.query().Where(s.ft.ByUserID(args.UserID)).
		Where(s.ft.ByAccountID(args.AccountID)).
		ShouldDelete(update); err != nil {
		return err
	}
	return nil
}
