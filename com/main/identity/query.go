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
	"o.o/capi/filter"
)

var _ identity.QueryService = &QueryService{}

type QueryService struct {
	userStore            sqlstore.UserStoreFactory
	accountStore         sqlstore.AccountStoreFactory
	partnerStore         sqlstore.PartnerStoreFactory
	partnerRelationStore sqlstore.PartnerRelationStoreFactory
	affiliateStore       sqlstore.AffiliateStoreFactory
	shopStore            sqlstore.ShopStoreFactory
	accountUserStore     sqlstore.AccountUserStoreFactory
	userRefSaffStore     sqlstore.UserRefSaffStoreFactory
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		userStore:            sqlstore.NewUserStore(db),
		accountStore:         sqlstore.NewAccountStore(db),
		partnerStore:         sqlstore.NewPartnerStore(db),
		partnerRelationStore: sqlstore.NewPartnerRelationStore(db),
		shopStore:            sqlstore.NewShopStore(db),
		affiliateStore:       sqlstore.NewAffiliateStore(db),
		accountUserStore:     sqlstore.NewAccountUserStore(db),
		userRefSaffStore:     sqlstore.NewUserRefSaffStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) identity.QueryBus {
	b := bus.New()
	h := identity.NewQueryServiceHandler(q)
	return h.RegisterHandlers(b)
}

func (q *QueryService) GetShopByID(ctx context.Context, id dot.ID) (*identity.Shop, error) {
	return q.shopStore(ctx).ByID(id).GetShop()
}

func (q *QueryService) GetShopByCode(ctx context.Context, code string) (*identity.Shop, error) {
	return q.shopStore(ctx).ByCode(code).GetShop()
}

func (q *QueryService) ListShopsByIDs(ctx context.Context, args *identity.ListShopsByIDsArgs) ([]*identity.Shop, error) {
	query := q.shopStore(ctx).ByIDs(args.IDs...)
	if args.IncludeWLPartnerShop {
		query.IncludeWLPartnerShop()
	}
	if args.IsPriorMoneyTransaction {
		query.IsPriorMoneyTransaction()
	}
	return query.ListShops()
}

func (q *QueryService) ListShopExtendeds(ctx context.Context, args *identity.ListShopQuery) (*identity.ListShopExtendedsResponse, error) {
	query := q.shopStore(ctx).Filters(args.Filters).WithPaging(args.Paging)
	if args.Name != "" {
		query = query.FullTextSearchName(args.Name)
	}
	if args.ShopIDs != nil && len(args.ShopIDs) > 0 {
		query = query.ByShopIDs(args.ShopIDs...)
	}
	if args.OwnerID != 0 {
		query = query.ByOwnerID(args.OwnerID)
	}
	if args.IncludeWLPartnerShop {
		query = query.IncludeWLPartnerShop()
	}
	if args.DateTo.Before(args.DateFrom) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "date_to must be after date_from")
	}
	if args.DateFrom.IsZero() != args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide both DateFrom and DateTo")
	}
	if !args.DateFrom.IsZero() {
		query = query.BetweenDateFromAndDateTo(args.DateFrom, args.DateTo)
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
	return q.userStore(ctx).ByID(args.UserID).GetUser()
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
		return nil, cm.Error(cm.InvalidArgument, "Vui lòng chỉ cung cấp email hoặc số điện thoại.", nil)
	}
	return query.GetUser()
}

func (q *QueryService) GetUserByPhone(ctx context.Context, phone string) (*identity.User, error) {
	return q.userStore(ctx).ByPhone(phone).GetUser()
}

func (q *QueryService) GetUserByEmail(ctx context.Context, email string) (*identity.User, error) {
	return q.userStore(ctx).ByEmail(email).GetUser()
}

func (q *QueryService) ListUsersByWLPartnerID(ctx context.Context, args *identity.ListUsersByWLPartnerID) ([]*identity.User, error) {
	return q.userStore(ctx).ByWLPartnerID(args.ID).ListUsers()
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

func (q *QueryService) GetAllAccountsByUsers(ctx context.Context, args *identity.GetAllAccountUsersArg) ([]*identity.AccountUser, error) {
	if len(args.UserIDs) == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing UserIDs", nil)
	}

	queryAccountUser := q.accountUserStore(ctx).ByUserIDs(args.UserIDs)
	if len(args.Roles) > 0 {
		queryAccountUser = queryAccountUser.ByRoles(args.Roles...)
	}
	accUser, err := queryAccountUser.ListAccountUserDBs()
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

func (q *QueryService) GetUserFtRefSaffByID(ctx context.Context, args *identity.GetUserByIDQueryArgs) (*identity.UserFtRefSaff, error) {
	return q.userStore(ctx).ByID(args.UserID).GetUserFtRefSaff(ctx)
}

func (q *QueryService) GetUsers(ctx context.Context, args *identity.ListUsersArgs) (*identity.UsersResponse, error) {
	query, err := q.buildCommonGetUserQuery(ctx, args.Name, args.Phone, args.Email, args.CreatedAt)
	if err != nil {
		return nil, err
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

func (q *QueryService) GetUserFtRefSaffs(ctx context.Context, args *identity.ListUserFtRefSaffsArgs) (*identity.UserFtRefSaffsResponse, error) {
	query, err := q.buildCommonGetUserQuery(ctx, args.Name, args.Phone, args.Email, args.CreatedAt)
	if err != nil {
		return nil, err
	}

	if args.RefAff != "" {
		refAffs, err := q.userRefSaffStore(ctx).ByRefAff(args.RefAff).ListUserRefSaff()
		if err != nil {
			return nil, err
		}
		var userIDs []dot.ID
		for _, user := range refAffs {
			userIDs = append(userIDs, user.UserID)
		}
		query = query.ByIDs(userIDs)
	}

	if args.RefSale != "" {
		refSales, err := q.userRefSaffStore(ctx).ByRefSale(args.RefSale).ListUserRefSaff()
		if err != nil {
			return nil, err
		}
		var userIDs []dot.ID
		for _, user := range refSales {
			userIDs = append(userIDs, user.UserID)
		}
		query = query.ByIDs(userIDs)
	}

	users, err := query.WithPaging(args.Paging).ListUserFtRefSaffs()
	if err != nil {
		return nil, err
	}
	return &identity.UserFtRefSaffsResponse{
		ListUsers: users,
		Paging:    query.GetPaging(),
	}, nil
}

func (q *QueryService) buildCommonGetUserQuery(ctx context.Context, name filter.FullTextSearch, phone, email string, createdAt filter.Date) (*sqlstore.UserStore, error) {
	query := q.userStore(ctx)
	if name != "" {
		query = query.ByFullNameNorm(name)
	}
	if phone != "" {
		phone, ok := validate.NormalizePhone(phone)
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
		}
		query = query.ByPhone(phone.String())
	}
	if email != "" {
		email, ok := validate.NormalizeEmail(email)
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Email không hợp lệ")
		}
		query = query.ByEmail(email.String())
	}
	if !createdAt.IsZero() {
		if !createdAt.From.IsZero() {
			query = query.ByCreatedAtFrom(createdAt.From.ToTime())
		}
		if !createdAt.To.IsZero() {
			query = query.ByCreatedAtTo(createdAt.To.ToTime())
		}
	}
	return query, nil
}

func (q *QueryService) GetAccountByID(ctx context.Context, ID dot.ID) (*identity.Account, error) {
	if ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Account ID")
	}
	return q.accountStore(ctx).ByID(ID).GetAccount()
}

func (q *QueryService) ListUsersByIDsAndNameNorm(ctx context.Context, args *identity.ListUsersByIDsAndNameNormArgs) ([]*identity.User, error) {
	query := q.userStore(ctx)
	if len(args.IDs) > 0 {
		query = query.ByIDs(args.IDs)
	}
	if len(args.NameNorm) > 0 {
		query = query.ByFullNameNorm(args.NameNorm)
	}
	return query.ListUsers()
}

func (q *QueryService) GetAccountUser(ctx context.Context, userID, accountID dot.ID) (*identity.AccountUser, error) {
	return q.accountUserStore(ctx).ByUserID(userID).ByAccountID(accountID).GetAccountUser()
}

func (q *QueryService) ListAccountUsers(ctx context.Context, args *identity.ListAccountUsersArgs) (*identity.ListAccountUsersResponse, error) {
	query := q.accountUserStore(ctx)

	if len(args.UserIDs) > 0 {
		query = query.ByUserIDs(args.UserIDs)
	}
	if args.AccountID != 0 {
		query = query.ByAccountID(args.AccountID)
	}
	if args.FullNameNorm != "" {
		query = query.ByFullNameNorm(args.FullNameNorm)
	}
	if args.PhoneNorm != "" {
		query = query.ByPhoneNorm(args.PhoneNorm)
	}
	if args.ExtensionNumberNorm != "" {
		query = query.ByExtensionNumberNorm(args.ExtensionNumberNorm)
	}
	if args.HasExtension.Valid {
		IsAssignedToExtension := args.HasExtension.Bool
		query = query.HasExtension(IsAssignedToExtension)
	}

	if args.HasDepartment.Valid {
		query = query.HasDepartment(args.HasDepartment.Bool)
	}

	if len(args.ExactRoles) > 0 {
		var roles []string
		for _, role := range args.ExactRoles {
			roles = append(roles, role.String())
		}
		query = query.ByExactRoles(roles...)
	}

	if len(args.Roles) > 0 {
		var roles []string
		for _, role := range args.Roles {
			roles = append(roles, role.String())
		}
		query = query.ByRoles(roles...)
	}

	if args.DepartmentID != 0 {
		query = query.ByDepartmentID(args.DepartmentID)
	}

	accountUsers, err := query.WithPaging(args.Paging).ListAccountUsers()
	if err != nil {
		return nil, err
	}
	return &identity.ListAccountUsersResponse{
		Paging:       query.GetPaging(),
		AccountUsers: accountUsers,
	}, nil
}

func (q *QueryService) ListExtendedAccountUsers(ctx context.Context, args *identity.ListExtendedAccountUsersArgs) (*identity.ListExtendedAccountUsersResponse, error) {
	// Get account users
	res, err := q.ListAccountUsers(ctx, &identity.ListAccountUsersArgs{
		Paging:              args.Paging,
		AccountID:           args.AccountID,
		FullNameNorm:        args.FullNameNorm,
		PhoneNorm:           args.PhoneNorm,
		ExtensionNumberNorm: args.ExtensionNumberNorm,
		Roles:               args.Roles,
		ExactRoles:          args.ExactRoles,
		UserIDs:             args.UserIDs,
		HasExtension:        args.HasExtension,
		HasDepartment:       args.HasDepartment,
		DepartmentID:        args.DepartmentID,
	})
	if err != nil {
		return nil, err
	}
	accountUsers := res.AccountUsers

	// Get UserIDs
	userIDs := make([]dot.ID, 0, len(accountUsers))
	for _, accUser := range accountUsers {
		userIDs = append(userIDs, accUser.UserID)
	}

	// Get users by userIDs and map user
	users, err := q.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	mapUser := map[dot.ID]identity.User{}
	for _, user := range users {
		mapUser[user.ID] = *user
	}

	var extendedAccountUsers []*identity.AccountUserExtended
	for _, accountUser := range accountUsers {
		var isDeleted bool
		if !accountUser.DeletedAt.IsZero() {
			isDeleted = true
		}
		user := mapUser[accountUser.UserID]
		extendedAccountUser := &identity.AccountUserExtended{
			UserID:       accountUser.UserID,
			AccountID:    accountUser.AccountID,
			DepartmentID: accountUser.DepartmentID,
			Roles:        accountUser.Roles,
			Permissions:  accountUser.Permissions,
			FullName:     user.FullName,
			ShortName:    user.ShortName,
			Email:        user.Email,
			Phone:        user.Phone,
			Position:     accountUser.Position,
			Deleted:      isDeleted,
		}
		extendedAccountUsers = append(extendedAccountUsers, extendedAccountUser)
	}

	return &identity.ListExtendedAccountUsersResponse{
		Paging:       res.Paging,
		AccountUsers: extendedAccountUsers,
	}, nil
}

func (q *QueryService) ListPartnerRelationsBySubjectIDs(ctx context.Context, args *identity.ListPartnerRelationsBySubjectIDsArgs) ([]*identity.PartnerRelation, error) {
	return q.partnerRelationStore(ctx).BySubjectType(args.SubjectType).BySubjectIDs(args.SubjectIDs...).ListPartnerRelations()
}

func (q *QueryService) ListAccountUsersByDepartmentIDs(ctx context.Context, ID dot.ID) ([]*identity.AccountUserWithGroupByDepartment, error) {
	accounts, err := q.accountUserStore(ctx).ByAccountID(ID).ListAccountUsersWithGroupByDepartment(ID)
	if err != nil {
		return nil, err
	}
	return convert.Convert_identitymodel_AccountUserWithGroupByDepartments_identity_AccountUserWithGroupByDepartments(accounts), nil
}
