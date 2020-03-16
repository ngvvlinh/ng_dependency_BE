package identity

import (
	"context"

	"etop.vn/api/main/identity"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/identity/convert"
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/backend/com/main/identity/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	userStore        sqlstore.UserStoreFactory
	accountStore     sqlstore.AccountStoreFactory
	partnerStore     sqlstore.PartnerStoreFactory
	accountUserStore sqlstore.AccountUserStoreFactory
	xAccountAhamove  sqlstore.XAccountAhamoveStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		userStore:        sqlstore.NewUserStore(db),
		accountStore:     sqlstore.NewAccountStore(db),
		partnerStore:     sqlstore.NewPartnerStore(db),
		accountUserStore: sqlstore.NewAccoutnUserStore(db),
		xAccountAhamove:  sqlstore.NewXAccountAhamoveStore(db),
	}
}

func (a *QueryService) MessageBus() identity.QueryBus {
	b := bus.New()
	return identity.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (a *QueryService) GetShopByID(ctx context.Context, id dot.ID) (*identity.Shop, error) {
	return a.accountStore(ctx).ShopByID(id).GetShop()
}

func (a *QueryService) GetUserByID(ctx context.Context, args *identity.GetUserByIDQueryArgs) (*identity.User, error) {
	return a.userStore(ctx).ByID(args.UserID).GetUser(ctx)
}

func (a *QueryService) GetUserByPhone(ctx context.Context, phone string) (*identity.User, error) {
	return a.userStore(ctx).ByPhone(phone).GetUser(ctx)
}

func (a *QueryService) GetUserByEmail(ctx context.Context, email string) (*identity.User, error) {
	return a.userStore(ctx).ByEmail(email).GetUser(ctx)
}

func (a *QueryService) GetExternalAccountAhamove(ctx context.Context, args *identity.GetExternalAccountAhamoveArgs) (*identity.ExternalAccountAhamove, error) {
	phone := args.Phone
	return a.xAccountAhamove(ctx).Phone(phone).OwnerID(args.OwnerID).GetXAccountAhamove()
}

func (a *QueryService) GetExternalAccountAhamoveByExternalID(ctx context.Context, args *identity.GetExternalAccountAhamoveByExternalIDQueryArgs) (*identity.ExternalAccountAhamove, error) {
	return a.xAccountAhamove(ctx).ExternalID(args.ExternalID).GetXAccountAhamove()
}

func (a *QueryService) GetAffiliateByID(ctx context.Context, id dot.ID) (*identity.Affiliate, error) {
	return a.accountStore(ctx).AffiliateByID(id).GetAffiliate()
}

func (a *QueryService) GetAffiliateWithPermission(ctx context.Context, affID dot.ID, userID dot.ID) (*identity.GetAffiliateWithPermissionResult, error) {
	if affID == 0 || userID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}
	res := &identity.GetAffiliateWithPermissionResult{}
	aff, err := a.GetAffiliateByID(ctx, affID)
	if err != nil {
		return nil, err
	}
	res.Affiliate = aff

	var accUser *identitymodel.AccountUser
	accUser, err = a.accountUserStore(ctx).GetAccountUserDB()
	if err != nil {
		return nil, err
	}
	res.Permission = convert.Permission(accUser.Permission)
	return res, nil
}

func (a *QueryService) GetAffiliatesByIDs(ctx context.Context, args *identity.GetAffiliatesByIDsArgs) ([]*identity.Affiliate, error) {
	return a.accountStore(ctx).AffiliatesByIDs(args.AffiliateIDs...).GetAffiliates()
}

func (a *QueryService) GetAffiliatesByOwnerID(ctx context.Context, args *identity.GetAffiliatesByOwnerIDArgs) ([]*identity.Affiliate, error) {
	return a.accountStore(ctx).AffiliatesByOwnerID(args.ID).GetAffiliates()
}

func (a *QueryService) ListPartnersForWhiteLabel(ctx context.Context, _ *meta.Empty) ([]*identity.Partner, error) {
	return a.partnerStore(ctx).WhiteLabel().ListPartners()
}
