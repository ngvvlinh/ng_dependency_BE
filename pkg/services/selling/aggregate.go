package selling

import (
	"context"

	"etop.vn/api/main/order"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/selling/convert"
	"etop.vn/backend/pkg/services/selling/sqlstore"
)

var _ order.Aggregate = &Aggregate{}

type Aggregate struct {
	s *sqlstore.OrderStore
}

func NewAggregate(db cmsql.Database) *Aggregate {
	return &Aggregate{s: sqlstore.New(db)}
}

func (a *Aggregate) GetOrderByID(ctx context.Context, args order.GetOrderByIDArgs) (*order.Order, error) {
	ord, err := a.s.WithContext(ctx).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return convert.Order(ord), nil
}
