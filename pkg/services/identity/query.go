package identity

import (
	"context"

	"etop.vn/backend/pkg/common/bus"

	"etop.vn/backend/pkg/common/cmsql"

	"etop.vn/api/main/identity"
	identitymodelx "etop.vn/backend/pkg/services/identity/modelx"
	"etop.vn/backend/pkg/services/identity/sqlstore"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	store           sqlstore.IdentityStoreFactory
	xAccountAhamove sqlstore.XAccountAhamoveStoreFactory
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		store:           sqlstore.NewIdentityStore(db),
		xAccountAhamove: sqlstore.NewXAccountAhamoveStore(db),
	}
}

func (a *QueryService) MessageBus() identity.QueryBus {
	b := bus.New()
	return identity.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (q *QueryService) GetShopByID(ctx context.Context, query *identity.GetShopByIDQueryArgs) (*identity.GetShopByIDQueryResult, error) {
	shop, err := q.store(ctx).GetByID(identitymodelx.GetByIDArgs{
		ID: query.ID,
	})
	if err != nil {
		return nil, err
	}
	return &identity.GetShopByIDQueryResult{
		Shop: shop,
	}, nil
}

func (q *QueryService) GetExternalAccountAhamoveByPhone(ctx context.Context, args *identity.GetExternalAccountAhamoveByPhoneArgs) (*identity.ExternalAccountAhamove, error) {
	return q.xAccountAhamove(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
}
