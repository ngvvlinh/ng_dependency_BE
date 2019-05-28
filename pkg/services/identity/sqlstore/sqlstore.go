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

type IdentityStoreFactory func(context.Context) *IdentityStore

type IdentityStore struct {
	query cmsql.Query
}

func NewIdentityStore(db cmsql.Database) IdentityStoreFactory {
	return func(ctx context.Context) *IdentityStore {
		return &IdentityStore{query: db.WithContext(ctx)}
	}
}

func (s *IdentityStore) GetByID(args identitymodelx.GetByIDArgs) (*identity.Shop, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing ID")
	}

	q := s.query.Where("id = ?", args.ID)
	result := &model.Shop{}
	if err := q.ShouldGet(result); err != nil {
		return nil, err
	}
	return convert.Shop(result), nil
}
