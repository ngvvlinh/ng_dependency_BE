package shipnow

import (
	"context"

	"o.o/api/main/shipnow"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipnow/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

var _ shipnow.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.ShipnowStoreFactory
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		store: sqlstore.NewShipnowStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) shipnow.QueryBus {
	b := bus.New()
	return shipnow.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetShipnowFulfillment(ctx context.Context, query *shipnow.GetShipnowFulfillmentQueryArgs) (*shipnow.GetShipnowFulfillmentQueryResult, error) {
	if query.ID == 0 && query.ShippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng cung cấp id hoặc shipping_code")
	}
	s := q.store(ctx).OptionalID(query.ID).OptionalShippingCode(query.ShippingCode)
	if query.ShopID != 0 {
		s = s.ShopID(query.ShopID)
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
