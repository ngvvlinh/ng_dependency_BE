package address

import (
	"context"

	"etop.vn/api/main/address"
	"etop.vn/backend/com/main/address/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ address.QueryService = &QueryService{}

type QueryService struct {
	s *sqlstore.AddressStore
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		s: sqlstore.NewAddressStore(db),
	}
}

func (q *QueryService) MessageBus() address.QueryBus {
	b := bus.New()
	return address.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetAddressByID(ctx context.Context, query *address.GetAddressByIDQueryArgs) (*address.Address, error) {
	addr, err := q.s.GetByID(query.ID)
	return addr, err
}
