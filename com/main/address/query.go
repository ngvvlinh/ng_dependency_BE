package address

import (
	"context"

	"o.o/api/main/address"
	"o.o/backend/com/main/address/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
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
