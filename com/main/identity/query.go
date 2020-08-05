package identity

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/meta"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/com/main/identity/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	userStore        sqlstore.UserStoreFactory
	accountStore     sqlstore.AccountStoreFactory
	partnerStore     sqlstore.PartnerStoreFactory
	affiliateStore   sqlstore.AffiliateStoreFactory
	shopStore        sqlstore.ShopStoreFactory
	accountUserStore sqlstore.AccountUserStoreFactory
	xAccountAhamove  sqlstore.XAccountAhamoveStoreFactory
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		userStore:        sqlstore.NewUserStore(db),
		accountStore:     sqlstore.NewAccountStore(db),
		partnerStore:     sqlstore.NewPartnerStore(db),
		shopStore:        sqlstore.NewShopStore(db),
		affiliateStore:   sqlstore.NewAffiliateStore(db),
		accountUserStore: sqlstore.NewAccountUserStore(db),
		xAccountAhamove:  sqlstore.NewXAccountAhamoveStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) identity.QueryBus {
	b := bus.New()
	h := identity.NewQueryServiceHandler(q)

	// TODO: refactor pkg/etop/sqlstore.CreateShop
	bus.AddHandler("sql", h.HandleGetUserByID)
	return h.RegisterHandlers(b)
}

func (q *QueryService) GetShopByID(ctx context.Context, id dot.ID) (*identity.Shop, error) {
	return q.shopStore(ctx).ByID(id).GetShop()
}

func (q *QueryService) ListShopsByIDs(ctx context.Context, ids []dot.ID) ([]*identity.Shop, error) {
	return q.shopStore(ctx).ByIDs(ids...).ListShops()
}

func (q *QueryService) ListShopExtendeds(ctx context.Context, args *identity.ListShopQuery) (*identity.ListShopExtendedsResponse, error) {
	query := q.shopStore(ctx).Filters(args.Filters).WithPaging(args.Paging)
	if args.Name != "" {
		query.FullTextSearchName(args.Name)
	}
	shops, err := query.ListShopExtendeds()
	if err != nil {
		return nil, err
	}
	return &identity.ListShopExtendedsResponse{
		Shops:  shops,
		Paging: query.GetPaging(),
	}, nil
}

func (q *QueryService) GetUserByID(ctx context.Context, args *identity.GetUserByIDQueryArgs) (*identity.User, error) {
	return q.userStore(ctx).ByID(args.UserID).GetUser(ctx)
}
func (q *QueryService) GetUsersByIDs(ctx context.Context, ids []dot.ID) ([]*identity.User, error) {
	return q.userStore(ctx).ByIDs(ids).ListUsers()
}

func (q *QueryService) GetUserByPhoneOrEmail(ctx context.Context, args *identity.GetUserByPhoneOrEmailArgs) (*identity.User, error) {
	count := 0
	query := q.userStore(ctx)

	if args.Phone != "" {
		count += 1
		query = query.ByPhone(args.Phone)
	}
	if args.Email != "" {
		count += 1
		query = query.ByEmail(args.Email)
	}
	if count != 1 {
		return nil, cm.Error(cm.InvalidArgument, "", nil)
	}
	return query.GetUser(ctx)
}

func (q *QueryService) GetUserByPhone(ctx context.Context, phone string) (*identity.User, error) {
	return q.userStore(ctx).ByPhone(phone).GetUser(ctx)
}

func (q *QueryService) GetUserByEmail(ctx context.Context, email string) (*identity.User, error) {
	return q.userStore(ctx).ByEmail(email).GetUser(ctx)
}

func (q *QueryService) ListUsersByWLPartnerID(ctx context.Context, args *identity.ListUsersByWLPartnerID) ([]*identity.User, error) {
	return q.userStore(ctx).ByWLPartnerID(args.ID).ListUsers()
}

func (q *QueryService) GetExternalAccountAhamove(ctx context.Context, args *identity.GetExternalAccountAhamoveArgs) (*identity.ExternalAccountAhamove, error) {
	phone := args.Phone
	return q.xAccountAhamove(ctx).Phone(phone).OwnerID(args.OwnerID).GetXAccountAhamove()
}

func (q *QueryService) GetExternalAccountAhamoveByExternalID(ctx context.Context, args *identity.GetExternalAccountAhamoveByExternalIDQueryArgs) (*identity.ExternalAccountAhamove, error) {
	return q.xAccountAhamove(ctx).ExternalID(args.ExternalID).GetXAccountAhamove()
}

func (q *QueryService) GetAffiliateByID(ctx context.Context, id dot.ID) (*identity.Affiliate, error) {
	return q.affiliateStore(ctx).ByID(id).GetAffiliate()
}

func (q *QueryService) GetAffiliateWithPermission(ctx context.Context, affID dot.ID, userID dot.ID) (*identity.GetAffiliateWithPermissionResult, error) {
	if affID == 0 || userID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}
	res := &identity.GetAffiliateWithPermissionResult{}
	aff, err := q.GetAffiliateByID(ctx, affID)
	if err != nil {
		return nil, err
	}
	res.Affiliate = aff

	var accUser *identitymodel.AccountUser
	accUser, err = q.accountUserStore(ctx).GetAccountUserDB()
	if err != nil {
		return nil, err
	}
	res.Permission = convert.Permission(accUser.Permission)
	return res, nil
}

func (q *QueryService) GetAffiliatesByIDs(ctx context.Context, args *identity.GetAffiliatesByIDsArgs) ([]*identity.Affiliate, error) {
	return q.affiliateStore(ctx).ByIDs(args.AffiliateIDs...).ListAffiliates()
}

func (q *QueryService) GetAffiliatesByOwnerID(ctx context.Context, args *identity.GetAffiliatesByOwnerIDArgs) ([]*identity.Affiliate, error) {
	return q.affiliateStore(ctx).ByOwnerID(args.ID).ListAffiliates()
}

func (q *QueryService) ListPartnersForWhiteLabel(ctx context.Context, _ *meta.Empty) ([]*identity.Partner, error) {
	return q.partnerStore(ctx).WhiteLabel().ListPartners()
}

func (q *QueryService) GetPartnerByID(ctx context.Context, args *identity.GetPartnerByIDArgs) (*identity.Partner, error) {
	return q.partnerStore(ctx).ByID(args.ID).GetPartner()
}

func (q *QueryService) GetUsers(ctx context.Context, args *identity.ListUsersArgs) (*identity.UsersResponse, error) {
	query := q.userStore(ctx)

	if args.Name != "" {
		query = query.ByNameNorm(validate.NormalizeSearch(args.Name))
	}
	if args.Phone != "" {
		_, phone, _ := validate.NormalizeEmailOrPhone(args.Phone)
		query = query.ByPhone(phone)
	}
	if args.Email != "" {
		email, _, _ := validate.NormalizeEmailOrPhone(args.Email)
		query = query.ByEmail(email)
	}
	if !args.CreatedAt.IsZero() {
		if !args.CreatedAt.From.IsZero() {
			query = query.ByCreatedAtFrom(args.CreatedAt.From.ToTime())
		}
		if !args.CreatedAt.To.IsZero() {
			query = query.ByCreatedAtTo(args.CreatedAt.To.ToTime())
		}
	}
	users, err := query.WithPaging(args.Paging).ListUsers()
	if err != nil {
		return nil, err
	}
	return &identity.UsersResponse{
		ListUsers: users,
		Paging:    query.GetPaging(),
	}, nil
}

func (q *QueryService) GetAllAccountsByUsers(ctx context.Context, args *identity.GetAllAccountUsersArg) ([]*identity.AccountUser, error) {
	if len(args.UserIDs) == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing UserIDs", nil)
	}

	accUser, err := q.accountUserStore(ctx).ByUserIDs(args.UserIDs).ListAccountUserDBs()
	if err != nil {
		return nil, err
	}
	var accountIDs []dot.ID
	for _, accountUser := range accUser {
		accountIDs = append(accountIDs, accountUser.AccountID)
	}
	accounts, err := q.accountStore(ctx).ByType(args.Type.Enum).ByIDs(accountIDs...).ListAccountDBs()

	mapAccount := make(map[dot.ID]bool)
	for _, account := range accounts {
		mapAccount[account.ID] = true
	}
	var result []*identitymodel.AccountUser
	for _, accountUser := range accUser {
		if mapAccount[accountUser.AccountID] {
			result = append(result, accountUser)
		}
	}
	return convert.Convert_identitymodel_AccountUsers_identity_AccountUsers(result), err
}

func (q *QueryService) GetUsersByAccount(ctx context.Context, accountID dot.ID) ([]*identity.AccountUser, error) {
	accountUsers, err := q.accountUserStore(ctx).ByAccountID(accountID).ListAccountUserDBs()
	return convert.Convert_identitymodel_AccountUsers_identity_AccountUsers(accountUsers), err
}
