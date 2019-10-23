package identity

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/com/main/identity/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	userStore        sqlstore.UserStoreFactory
	accountStore     sqlstore.AccountStoreFactory
	accountUserStore sqlstore.AccountUserStoreFactory
	xAccountAhamove  sqlstore.XAccountAhamoveStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		userStore:        sqlstore.NewUserStore(db),
		accountStore:     sqlstore.NewAccountStore(db),
		accountUserStore: sqlstore.NewAccoutnUserStore(db),
		xAccountAhamove:  sqlstore.NewXAccountAhamoveStore(db),
	}
}

func (a *QueryService) MessageBus() identity.QueryBus {
	b := bus.New()
	return identity.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (q *QueryService) GetShopByID(ctx context.Context, id int64) (*identity.Shop, error) {
	return q.accountStore(ctx).ShopByID(id).GetShop()
}

func (q *QueryService) GetUserByID(ctx context.Context, args *identity.GetUserByIDQueryArgs) (*identity.User, error) {
	return q.userStore(ctx).ByID(args.UserID).GetUser()
}

func (q *QueryService) GetUserByPhone(ctx context.Context, phone string) (*identity.User, error) {
	return q.userStore(ctx).ByPhone(phone).GetUser()
}

func (q *QueryService) GetExternalAccountAhamove(ctx context.Context, args *identity.GetExternalAccountAhamoveArgs) (*identity.ExternalAccountAhamove, error) {
	phone := args.Phone
	return q.xAccountAhamove(ctx).Phone(phone).OwnerID(args.OwnerID).GetXAccountAhamove()
}

func (q *QueryService) GetExternalAccountAhamoveByExternalID(ctx context.Context, args *identity.GetExternalAccountAhamoveByExternalIDQueryArgs) (*identity.ExternalAccountAhamove, error) {
	return q.xAccountAhamove(ctx).ExternalID(args.ExternalID).GetXAccountAhamove()
}

func (q *QueryService) GetAffiliateByID(ctx context.Context, id int64) (*identity.Affiliate, error) {
	return q.accountStore(ctx).AffiliateByID(id).GetAffiliate()
}

func (q *QueryService) GetAffiliateWithPermission(ctx context.Context, affID int64, userID int64) (*identity.GetAffiliateWithPermissionResult, error) {
	if affID == 0 || userID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}
	res := &identity.GetAffiliateWithPermissionResult{}
	aff, err := q.GetAffiliateByID(ctx, affID)
	if err != nil {
		return nil, err
	}
	res.Affiliate = aff

	var accUser *model.AccountUser
	accUser, err = q.accountUserStore(ctx).GetAccountUserDB()
	if err != nil {
		return nil, err
	}
	res.Permission = convert.Permission(accUser.Permission)
	return res, nil
}

func (q *QueryService) GetAffiliatesByIDs(ctx context.Context, args *identity.GetAffiliatesByIDsArgs) ([]*identity.Affiliate, error) {
	return q.accountStore(ctx).AffiliatesByIDs(args.AffiliateIDs...).GetAffiliates()
}

func (q *QueryService) GetAffiliatesByOwnerID(ctx context.Context, args *identity.GetAffiliatesByOwnerIDArgs) ([]*identity.Affiliate, error) {
	return q.accountStore(ctx).AffiliatesByOwnerID(args.ID).GetAffiliates()
}
