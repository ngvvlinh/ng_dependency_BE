package identity

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"

	"etop.vn/backend/pkg/common/bus"

	"etop.vn/backend/pkg/common/cmsql"

	"etop.vn/api/main/identity"
	identitymodelx "etop.vn/backend/pkg/services/identity/modelx"
	"etop.vn/backend/pkg/services/identity/sqlstore"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	store           sqlstore.ShopStoreFactory
	userStore       sqlstore.UserStoreFactory
	xAccountAhamove sqlstore.XAccountAhamoveStoreFactory
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		store:           sqlstore.NewIdentityStore(db),
		userStore:       sqlstore.NewUserStore(db),
		xAccountAhamove: sqlstore.NewXAccountAhamoveStore(db),
	}
}

func (a *QueryService) MessageBus() identity.QueryBus {
	b := bus.New()
	return identity.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (q *QueryService) GetShopByID(ctx context.Context, args *identity.GetShopByIDQueryArgs) (*identity.GetShopByIDQueryResult, error) {
	shop, err := q.store(ctx).GetByID(identitymodelx.GetByIDArgs{
		ID: args.ID,
	})
	if err != nil {
		return nil, err
	}
	return &identity.GetShopByIDQueryResult{
		Shop: shop,
	}, nil
}

func (q *QueryService) GetUserByID(ctx context.Context, args *identity.GetUserByIDQueryArgs) (*identity.User, error) {
	return q.userStore(ctx).GetUserByID(sqlstore.GetUserByIDArgs{
		ID: args.UserID,
	})
}

func (q *QueryService) GetExternalAccountAhamove(ctx context.Context, args *identity.GetExternalAccountAhamoveArgs) (*identity.ExternalAccountAhamove, error) {
	phone := args.Phone
	if cm.IsDev() {
		phone, _, _ = validate.TrimTest(phone)
	}
	return q.xAccountAhamove(ctx).Phone(phone).OwnerID(args.OwnerID).GetXAccountAhamove()
}

func (q *QueryService) GetExternalAccountAhamoveByExternalID(ctx context.Context, args *identity.GetExternalAccountAhamoveByExternalIDQueryArgs) (*identity.ExternalAccountAhamove, error) {
	return q.xAccountAhamove(ctx).ExternalID(args.ExternalID).GetXAccountAhamove()
}
