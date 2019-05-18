package sqlstore

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/services/identity/convert"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"

	"etop.vn/backend/pkg/common/cmsql"
)

type IdentityStore struct {
	ctx context.Context
	db  cmsql.Database
}

func NewIdentityStore(db cmsql.Database) *IdentityStore {
	return &IdentityStore{
		ctx: context.Background(),
		db:  db,
	}
}

func (s *IdentityStore) WithContext(ctx context.Context) *IdentityStore {
	return &IdentityStore{
		ctx: ctx,
		db:  s.db,
	}
}

func (s *IdentityStore) GetByID(ID int64) (*identity.Shop, error) {
	if ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing ID")
	}

	q := s.db.WithContext(s.ctx).Where("id = ?", ID)
	result := &model.Shop{}
	if err := q.ShouldGet(result); err != nil {
		return nil, err
	}
	return convert.Shop(result), nil
}
