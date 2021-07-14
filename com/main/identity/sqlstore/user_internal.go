package sqlstore

import (
	"context"

	"o.o/api/main/identity"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type UserInternalStoreFactory func(context.Context) *UserInternalStore

func NewUserInternalStore(db *cmsql.Database) UserInternalStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *UserInternalStore {
		return &UserInternalStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type UserInternalStore struct {
	query cmsql.QueryFactory
	ft    UserInternalFilters
	preds []interface{}
	ctx   context.Context
}

func (s *UserInternalStore) UserID(userID dot.ID) *UserInternalStore {
	s.preds = append(s.preds, s.ft.ByID(userID))
	return s
}

func (s *UserInternalStore) CreateUserInternal(user *identity.UserInternal) error {
	sqlstore.MustNoPreds(s.preds)
	if user.ID == 0 {
		user.ID = cm.NewID()
	}
	userDB := &identitymodel.UserInternal{}
	if err := scheme.Convert(user, userDB); err != nil {
		return err
	}
	return s.query().ShouldInsert(userDB)
}

func (s *UserInternalStore) UpdateUserInternal(user *identity.UserInternal) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "UpdateUserInternal: Must provide preds")
	}

	var userInternalDB identitymodel.UserInternal
	if err := scheme.Convert(user, &userInternalDB); err != nil {
		return err
	}
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(&userInternalDB)
}
