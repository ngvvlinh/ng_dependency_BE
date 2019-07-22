package identity

import (
	"context"

	"etop.vn/api/external/haravan/identity"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/external/haravan/identity/sqlstore"
	"etop.vn/common/bus"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.XAccountHaravanStoreFactory
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewXAccountHaravanStore(db),
	}
}

func (q *QueryService) MessageBus() identity.QueryBus {
	b := bus.New()
	return identity.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetExternalAccountHaravanByShopID(ctx context.Context, query *identity.GetExternalAccountHaravanByShopIDQueryArgs) (*identity.ExternalAccountHaravan, error) {
	s := q.store(ctx).ShopID(query.ShopID)
	return s.GetXAccountHaravan()
}

func (q *QueryService) GetExternalAccountHaravanByXShopID(ctx context.Context, query *identity.GetExternalAccountHaravanByXShopIDQueryArgs) (*identity.ExternalAccountHaravan, error) {
	s := q.store(ctx).ExternalShopID(query.ExternalShopID)
	return s.GetXAccountHaravan()
}
