package address

import (
	"context"

	"etop.vn/api/main/address"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/address/sqlstore"
)

var _ address.QueryService = &QueryService{}

type QueryService struct {
	s *sqlstore.AddressStore
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		s: sqlstore.NewAddressStore(db),
	}
}

func (q *QueryService) GetAddressByID(ctx context.Context, query *address.GetAddressByIDQueryArgs) (*address.Address, error) {
	address, err := q.s.GetByID(query.ID)
	return address, err
}
