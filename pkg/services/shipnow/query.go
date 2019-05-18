package shipnow

import (
	"context"

	"etop.vn/api/main/shipnow"
	"etop.vn/backend/pkg/common/cmsql"
	shipnowconvert "etop.vn/backend/pkg/services/shipnow/convert"
	"etop.vn/backend/pkg/services/shipnow/model"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

var _ shipnow.QueryService = &QueryService{}

type QueryService struct {
	s *sqlstore.ShipnowStore
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		s: sqlstore.NewShipnowStore(db),
	}
}

func (q *QueryService) GetShipnowFulfillment(ctx context.Context, query *shipnow.GetShipnowFulfillmentQueryArgs) (*shipnow.GetShipnowFulfillmentQueryResult, error) {
	ffm, err := q.s.WithContext(ctx).GetByID(model.GetByIDArgs{
		ID:     query.Id,
		ShopID: query.ShopId,
	})
	if err != nil {
		return nil, err
	}
	return &shipnow.GetShipnowFulfillmentQueryResult{
		ShipnowFulfillment: ffm,
	}, nil
}

func (q *QueryService) GetShipnowFulfillments(ctx context.Context, query *shipnow.GetShipnowFulfillmentsQueryArgs) (*shipnow.GetShipnowFulfillmentsQueryResult, error) {
	args := &model.GetShipnowFulfillmentsArgs{
		ShopID: query.ShopId,
	}
	ffms, err := q.s.WithContext(ctx).GetShipnowFulfillments(args)
	if err != nil {
		return nil, err
	}
	return &shipnow.GetShipnowFulfillmentsQueryResult{
		ShipnowFulfillments: shipnowconvert.Shipnows(ffms),
	}, nil
}
