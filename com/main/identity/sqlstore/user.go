package sqlstore

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/backend/com/main/identity/convert"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
)

type UserStoreFactory func(context.Context) *UserStore

func NewUserStore(db *cmsql.Database) UserStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *UserStore {
		return &UserStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type UserStore struct {
	query func() cmsql.QueryInterface
	ft    sqlstore.UserFilters
	preds []interface{}
}

func (s *UserStore) ByID(id int64) *UserStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *UserStore) ByPhone(phone string) *UserStore {
	s.preds = append(s.preds, s.ft.ByPhone(phone))
	return s
}

func (s *UserStore) GetUserDB() (*model.User, error) {
	var user model.User
	err := s.query().Where(s.preds).ShouldGet(&user)
	return &user, err
}

func (s *UserStore) GetUser() (*identity.User, error) {
	result, err := s.GetUserDB()
	if err != nil {
		return nil, err
	}
	return convert.User(result), nil
}

type UpdateRefferenceIDArgs struct {
	UserID    int64
	RefUserID int64
	RefSaleID int64
}

func (s *UserStore) UpdateUserRefferenceID(args *UpdateRefferenceIDArgs) error {
	if args.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing userID")
	}
	if args.RefUserID == 0 && args.RefSaleID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing reference userID")
	}
	user := &model.User{
		RefUserID: args.RefUserID,
		RefSaleID: args.RefSaleID,
	}

	if err := s.query().Where(s.ft.ByID(args.UserID)).ShouldUpdate(user); err != nil {
		return err
	}
	return nil
}
