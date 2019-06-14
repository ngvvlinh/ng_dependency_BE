package sqlstore

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/services/identity/convert"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"

	"etop.vn/backend/pkg/common/cmsql"
	identitymodelx "etop.vn/backend/pkg/services/identity/modelx"
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

func (s *ShopStore) GetByID(args identitymodelx.GetByIDArgs) (*identity.Shop, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing ID")
	}

	q := s.query().Where("id = ?", args.ID)
	result := &model.Shop{}
	if err := q.ShouldGet(result); err != nil {
		return nil, err
	}
	return convert.Shop(result), nil
}
