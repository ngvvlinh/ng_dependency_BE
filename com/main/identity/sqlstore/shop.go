package sqlstore

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
)

type ShopStoreFactory func(context.Context) *ShopStore

func NewIdentityStore(db cmsql.Database) ShopStoreFactory {
	return func(ctx context.Context) *ShopStore {
		return &ShopStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type ShopStore struct {
	query func() cmsql.QueryInterface
	preds []interface{}
}

func (s *ShopStore) GetByID(args modelx.GetByIDArgs) (*identity.Shop, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	q := s.query().Where("id = ?", args.ID)
	result := &model.Shop{}
	if err := q.ShouldGet(result); err != nil {
		return nil, err
	}
	return convert.Shop(result), nil
}
