package shipnow

import (
	"context"

	"etop.vn/api/main/shipnow"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

var _ shipnow.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.ShipnowStoreFactory
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewShipnowStore(db),
	}
}

func (a *QueryService) MessageBus() shipnow.QueryBus {
	b := bus.New()
	return shipnow.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (q *QueryService) GetShipnowFulfillment(ctx context.Context, query *shipnow.GetShipnowFulfillmentQueryArgs) (*shipnow.GetShipnowFulfillmentQueryResult, error) {
	s := q.store(ctx).ID(query.Id)
	if query.ShopId != 0 {
		s = s.ShopID(query.ShopId)
	}
	ffm, err := s.GetShipnow()
	if err != nil {
		return nil, err
	}
	return &shipnow.GetShipnowFulfillmentQueryResult{
		ShipnowFulfillment: ffm,
	}, nil
}

func (q *QueryService) GetShipnowFulfillments(ctx context.Context, query *shipnow.GetShipnowFulfillmentsQueryArgs) (*shipnow.GetShipnowFulfillmentsQueryResult, error) {
	s := q.store(ctx).ShopID(query.ShopId).Filters(nil)
	ffms, err := s.ListShipnows(nil)
	if err != nil {
		return nil, err
	}
	// count, err := s.Count()
	// if err != nil {
	// 	return nil, err
	// }

	return &shipnow.GetShipnowFulfillmentsQueryResult{
		ShipnowFulfillments: ffms,
	}, nil
}

func (q *QueryService) GetShipnowFulfillmentByShippingCode(ctx context.Context, query *shipnow.GetShipnowFulfillmentByShippingCodeQueryArgs) (*shipnow.GetShipnowFulfillmentQueryResult, error) {
	ffm, err := q.store(ctx).ShippingCode(query.ShippingCode).GetShipnow()
	if err != nil {
		return nil, err
	}
	return &shipnow.GetShipnowFulfillmentQueryResult{
		ShipnowFulfillment: ffm,
	}, nil
}
