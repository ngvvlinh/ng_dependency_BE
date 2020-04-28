package shipnow

import (
	"context"

	"o.o/api/main/shipnow"
	"o.o/backend/com/main/shipnow/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ shipnow.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.ShipnowStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewShipnowStore(db),
	}
}

func (q *QueryService) MessageBus() shipnow.QueryBus {
	b := bus.New()
	return shipnow.NewQueryServiceHandler(q).RegisterHandlers(b)
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
	s := q.store(ctx).ShopIDs(query.ShopIds...).Filters(query.Filters)
	if query.Paging != nil && len(query.Paging.Sort) == 0 {
		query.Paging.Sort = []string{"-created_at"}
	}
	ffms, err := s.ListShipnows(query.Paging)
	if err != nil {
		return nil, err
	}
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
