package ordering

import (
	"context"

	"etop.vn/api/main/ordering"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/ordering/convert"
	"etop.vn/backend/pkg/services/ordering/sqlstore"
)

var _ ordering.Aggregate = &Aggregate{}

type Aggregate struct {
	s *sqlstore.OrderStore
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{s: sqlstore.NewOrderStore(db)}
}

func (a *Aggregate) GetOrderByID(ctx context.Context, args ordering.GetOrderByIDArgs) (*ordering.Order, error) {
	ord, err := a.s.WithContext(ctx).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return convert.Order(ord), nil
}
