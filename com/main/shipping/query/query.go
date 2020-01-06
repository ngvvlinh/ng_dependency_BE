package shipping

import (
	"context"

	"etop.vn/api/main/shipping"
	"etop.vn/backend/com/main/shipping/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ shipping.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.FulfillmentStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewFulfillmentStore(db),
	}
}

func (q *QueryService) MessageBus() shipping.QueryBus {
	b := bus.New()
	return shipping.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetFulfillmentByIDOrShippingCode(ctx context.Context, id dot.ID, shippingCode string) (*shipping.Fulfillment, error) {
	if id == 0 && shippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	query := q.store(ctx)
	if id != 0 {
		query = query.ID(id)
	}
	if shippingCode != "" {
		query = query.ShippingCode(shippingCode)
	}
	return query.GetFulfillment()
}
