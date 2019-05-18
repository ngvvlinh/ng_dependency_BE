package identity

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"

	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/services/identity/sqlstore"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	s *sqlstore.IdentityStore
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		s: sqlstore.NewIdentityStore(db),
	}
}

func (q *QueryService) GetShopByID(ctx context.Context, query *identity.GetShopByIDQueryArgs) (*identity.Shop, error) {
	shop, err := q.s.GetByID(query.ID)
	return shop, err
}
