package sqlstore

import (
	"context"

	"etop.vn/backend/com/main/identity/convert"

	"etop.vn/api/main/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
)

type UserStoreFactory func(context.Context) *UserStore

func NewUserStore(db cmsql.Database) UserStoreFactory {
	return func(ctx context.Context) *UserStore {
		return &UserStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type UserStore struct {
	query func() cmsql.QueryInterface
	preds []interface{}
}

type GetUserByIDArgs struct {
	ID int64
}

func (s *UserStore) GetUserByID(args GetUserByIDArgs) (*identity.User, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing ID")
	}

	q := s.query().Table("user").Where("id = ?", args.ID)
	result := &model.User{}
	if err := q.ShouldGet(result); err != nil {
		return nil, err
	}
	return convert.User(result), nil
}
