package identity

import (
	"context"
	"strings"
	"time"

	"o.o/api/main/authorization"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/try_on"
	"o.o/api/top/types/etc/user_source"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitystore "o.o/backend/com/main/identity/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var ll = l.New()

var _ identity.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

const (
	UserEmailKey            = "user_email_key"
	UserPhoneKey            = "user_phone_key"
	UserPhoneWLPartnerIDKey = "user_phone_wl_partner_id_idx"
	UserEmailWLPartnerIDKey = "user_email_wl_partner_id_idx"

	MsgCreateUserDuplicatedPhone = `Số điện thoại đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng số điện thoại khác. Nếu cần thêm thông tin, vui lòng liên hệ %v.`
	MsgCreateUserDuplicatedEmail = `Email đã được sử dụng. Vui lòng đăng nhập hoặc sử dụng email khác. Nếu cần thêm thông tin, vui lòng liên hệ %v.`
)

var mapUserError = map[string]string{
	UserEmailKey:            MsgCreateUserDuplicatedEmail,
	UserEmailWLPartnerIDKey: MsgCreateUserDuplicatedEmail,
	UserPhoneKey:            MsgCreateUserDuplicatedPhone,
	UserPhoneWLPartnerIDKey: MsgCreateUserDuplicatedPhone,
}

type Aggregate struct {
	db                *cmsql.Database
	userStore         identitystore.UserStoreFactory
	userInternalStore identitystore.UserInternalStoreFactory
	affiliateStore    identitystore.AffiliateStoreFactory
	shopStore         identitystore.ShopStoreFactory
	accountStore      identitystore.AccountStoreFactory
	accountUserStore  identitystore.AccountUserStoreFactory
	userRefSaffStore  identitystore.UserRefSaffStoreFactory
	eventBus          capi.EventBus
}

func NewAggregate(db com.MainDB, eventBus capi.EventBus) *Aggregate {
	return &Aggregate{
		db:                db,
		userStore:         identitystore.NewUserStore(db),
		userInternalStore: identitystore.NewUserInternalStore(db),
		shopStore:         identitystore.NewShopStore(db),
		affiliateStore:    identitystore.NewAffiliateStore(db),
		accountStore:      identitystore.NewAccountStore(db),
		accountUserStore:  identitystore.NewAccountUserStore(db),
		userRefSaffStore:  identitystore.NewUserRefSaffStore(db),
		eventBus:          eventBus,
	}
}

func AggregateMessageBus(a *Aggregate) identity.CommandBus {
	b := bus.New()
	return identity.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) UpdateUserEmail(ctx context.Context, userID dot.ID, email string) error {
	if email == "" || userID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	res, ok := validate.NormalizeEmail(email)
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "Email không hợp lệ")
	}
	_, err := a.userStore(ctx).ByID(userID).UpdateUserEmail(res.String())
	return err
}

func (a *Aggregate) UpdateUserPhone(ctx context.Context, userID dot.ID, phone string) error {
	if phone == "" || userID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	res, isPhone := validate.NormalizePhone(phone)
	if !isPhone {
		return cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	_, err := a.userStore(ctx).ByID(userID).UpdateUserPhone(res.String())
	return err
}

func (a *Aggregate) UpdateUserReferenceUserID(ctx context.Context, args *identity.UpdateUserReferenceUserIDArgs) error {
	if currentUser, err := a.userStore(ctx).ByID(args.UserID).GetUserDB(); err != nil {
		return err
	} else if currentUser.RefSaleID != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "RefUserID đã tồn tại. Không thể cập nhật.")
	}
	refUser, err := a.userStore(ctx).ByPhone(args.RefUserPhone).GetUserDB()
	if err != nil {
		return cm.Errorf(cm.NotFound, nil, "Số điện thoại người dùng không tồn tại.")
	}

	updateCmd := &identitystore.UpdateRefferenceIDArgs{
		UserID:    args.UserID,
		RefUserID: refUser.ID,
	}
	return a.userStore(ctx).UpdateUserRefferenceID(updateCmd)
}

func (a *Aggregate) UpdateUserReferenceSaleID(ctx context.Context, args *identity.UpdateUserReferenceSaleIDArgs) error {
	if currentUser, err := a.userStore(ctx).ByID(args.UserID).GetUserDB(); err != nil {
		return err
	} else if currentUser.RefSaleID != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "RefSaleID đã tồn tại. Không thể cập nhật.")
	}
	refUser, err := a.userStore(ctx).ByPhone(args.RefSalePhone).GetUserDB()
	if err != nil {
		return cm.Errorf(cm.NotFound, nil, "Số điện thoại người dùng không tồn tại.")
	}

	updateCmd := &identitystore.UpdateRefferenceIDArgs{
		UserID:    args.UserID,
		RefSaleID: refUser.ID,
	}
	return a.userStore(ctx).UpdateUserRefferenceID(updateCmd)
}

func (a *Aggregate) CreateAffiliate(ctx context.Context, args *identity.CreateAffiliateArgs) (*identity.Affiliate, error) {
	var ok bool
	var emailNorm model.NormalizedEmail
	var phoneNorm model.NormalizedPhone
	if args.Name, ok = validate.NormalizeName(args.Name); !ok {
		return nil, cm.Error(cm.InvalidArgument, "Tên người dùng không hợp lệ", nil)
	}
	if args.Email != "" {
		if emailNorm, ok = validate.NormalizeEmail(args.Email); !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Email không hợp lệ")
		}
	}
	if phoneNorm, ok = validate.NormalizePhone(args.Phone); !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	_args := identitystore.CreateAffiliateArgs{
		Name:        args.Name,
		OwnerID:     args.OwnerID,
		Phone:       phoneNorm.String(),
		Email:       emailNorm.String(),
		IsTest:      args.IsTest,
		BankAccount: args.BankAccount,
	}
	return a.affiliateStore(ctx).CreateAffiliate(_args)
}

func (a *Aggregate) UpdateAffiliateInfo(ctx context.Context, args *identity.UpdateAffiliateInfoArgs) (*identity.Affiliate, error) {
	var ok bool
	var emailNorm model.NormalizedEmail
	var phoneNorm model.NormalizedPhone

	if args.Name != "" {
		if args.Name, ok = validate.NormalizeName(args.Name); !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Tên người dùng không hợp lệ")
		}
	}
	if args.Email != "" {
		if emailNorm, ok = validate.NormalizeEmail(args.Email); !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Email không hợp lệ")
		}
	}
	if args.Phone != "" {
		if phoneNorm, ok = validate.NormalizePhone(args.Phone); !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
		}
	}
	_args := identitystore.UpdateAffiliateArgs{
		ID:      args.ID,
		OwnerID: args.OwnerID,
		Phone:   phoneNorm.String(),
		Email:   emailNorm.String(),
		Name:    args.Name,
	}
	return a.affiliateStore(ctx).UpdateAffiliate(_args)
}

func (a *Aggregate) UpdateAffiliateBankAccount(ctx context.Context, args *identity.UpdateAffiliateBankAccountArgs) (*identity.Affiliate, error) {
	_args := identitystore.UpdateAffiliateArgs{
		ID:          args.ID,
		OwnerID:     args.OwnerID,
		BankAccount: args.BankAccount,
	}
	return a.affiliateStore(ctx).UpdateAffiliate(_args)
}

func (a *Aggregate) DeleteAffiliate(ctx context.Context, args *identity.DeleteAffiliateArgs) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		args1 := identitystore.DeleteAffiliateArgs{
			ID:      args.ID,
			OwnerID: args.OwnerID,
		}
		if err := a.affiliateStore(ctx).DeleteAffiliate(args1); err != nil {
			return err
		}
		args2 := identitystore.DeleteAccountUserArgs{
			AccountID: args.ID,
			UserID:    args.OwnerID,
		}
		return a.accountUserStore(ctx).DeleteAccountUser(args2)
	})
}

func (a *Aggregate) BlockUser(ctx context.Context, args *identity.BlockUserArgs) (*identity.User, error) {
	user, err := a.userStore(ctx).ByID(args.UserID).GetUserDB()
	if err != nil {
		return nil, err
	}
	if user.Status == status3.N {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể khóa tài khoản đã bị khóa")
	}
	userUpdate := &identitymodel.User{
		BlockReason: args.BlockReason,
		Status:      status3.N,
		BlockedAt:   time.Now(),
		BlockedBy:   args.BlockBy,
	}
	err = a.userStore(ctx).ByID(user.ID).UpdateUserDB(userUpdate)
	if err != nil {
		return nil, err
	}
	user, err = a.userStore(ctx).ByID(args.UserID).GetUserDB()
	if err != nil {
		return nil, err
	}
	shopUpdate := &identitymodel.Shop{
		Status: status3.N,
	}
	err = a.shopStore(ctx).ByOwnerID(args.UserID).NotDeleted().UpdateShopDB(shopUpdate)
	if err != nil {
		return nil, err
	}
	return convert.User(user), nil
}

func (a *Aggregate) UnblockUser(ctx context.Context, userID dot.ID) (*identity.User, error) {
	user, err := a.userStore(ctx).ByID(userID).GetUserDB()
	if err != nil {
		return nil, err
	}
	if user.Status != status3.N {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tài khoản chưa bị khóa vui lòng kiểm tra lại")
	}
	_, err = a.userStore(ctx).ByID(userID).UnblockUser()
	if err != nil {
		return nil, err
	}
	user, err = a.userStore(ctx).ByID(userID).GetUserDB()
	if err != nil {
		return nil, err
	}
	shopUpdate := &identitymodel.Shop{
		Status: status3.P,
	}
	err = a.shopStore(ctx).ByOwnerID(userID).NotDeleted().UpdateShopDB(shopUpdate)
	if err != nil {
		return nil, err
	}
	return convert.User(user), nil
}

func (a *Aggregate) UpdateUserRef(ctx context.Context, args *identity.UpdateUserRefArgs) (*identity.UserRefSaff, error) {
	// Validate user
	_, err := a.userStore(ctx).ByID(args.UserID).IncludeWLPartnerUser().GetUserDB()
	if err != nil {
		return nil, err
	}

	userRef, err := a.userRefSaffStore(ctx).ByUserID(args.UserID).GetUserRefSaff()
	if err != nil {
		// if not exist, create new
		if cm.ErrorCode(err) != cm.NotFound {
			return nil, err
		}

		userRef = &identity.UserRefSaff{
			UserID:  args.UserID,
			RefAff:  args.RefAff,
			RefSale: args.RefSale,
		}
		if err = a.userRefSaffStore(ctx).CreateUserRefSaff(userRef); err != nil {
			return nil, err
		}
		return userRef, nil
	}

	// update shopUserRef
	updateUserRef := &identity.UserRefSaff{
		UserID:  args.UserID,
		RefAff:  args.RefAff,
		RefSale: args.RefSale,
	}
	if err = a.userRefSaffStore(ctx).ByUserID(args.UserID).Update(updateUserRef); err != nil {
		return nil, err
	}
	return updateUserRef, nil
}

func (a *Aggregate) UpdateShipFromAddressID(ctx context.Context, args *identity.UpdateShipFromAddressArgs) error {
	_, err := a.shopStore(ctx).ByID(args.ID).UpdateShipFromAddressID(args.ShipFromAddressID)
	return err
}

func (a *Aggregate) RegisterSimplify(ctx context.Context, args *identity.RegisterSimplifyArgs) error {
	normalizePhone, ok := validate.NormalizePhone(args.Phone)
	if !ok {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số điện thoại không hợp lệ")
	}

	_, err := a.userStore(ctx).ByPhone(normalizePhone.String()).GetUser()
	if err == nil {
		// phone đã tồn tại
		// exit
		return nil
	}
	if cm.ErrorCode(err) != cm.NotFound {
		return err
	}

	// create new user + a default shop for this user
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		now := time.Now()
		userArgs := &identity.CreateUserArgs{
			Phone:                   normalizePhone.String(),
			FullName:                args.FullName,
			Password:                args.Password,
			Status:                  status3.P,
			Source:                  user_source.Etop,
			PhoneVerifiedAt:         now,
			PhoneVerificationSentAt: now,
			Email:                   args.Email,
		}
		user, err := a.CreateUser(ctx, userArgs)
		if err != nil {
			return err
		}

		if args.IsCreateDefaultShop {
			name := normalizePhone.String()
			if args.CompanyName != "" {
				name = args.CompanyName
			}
			shopArgs := &identity.CreateShopArgs{
				Name:    name,
				OwnerID: user.ID,
				Phone:   normalizePhone.String(),
			}
			return a.createShop(ctx, shopArgs)
		}
		return nil
	})
}

func (a *Aggregate) CreateShop(ctx context.Context, args *identity.CreateShopArgs) (*identity.Shop, error) {
	var err error
	args.ID, err = checkShopID(args.ID)
	if err != nil {
		return nil, err
	}
	if err := a.createShop(ctx, args); err != nil {
		return nil, err
	}
	return a.shopStore(ctx).ByID(args.ID).GetShop()
}

func checkShopID(id dot.ID) (dot.ID, error) {
	if id == 0 {
		id = idutil.NewShopID()
	} else {
		if !idutil.IsShopID(id) {
			return 0, cm.Errorf(cm.InvalidArgument, nil, "invalid shop ID")
		}
	}
	return id, nil
}

func (a *Aggregate) createShop(ctx context.Context, args *identity.CreateShopArgs) error {
	_, err := a.userStore(ctx).ByID(args.OwnerID).GetUser()
	if err != nil {
		return cm.Error(cm.InvalidArgument, "invalid owner_id", nil)
	}

	id, err := checkShopID(args.ID)
	if err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		account := &identity.Account{
			ID:       id,
			Name:     args.Name,
			Type:     account_type.Shop,
			ImageURL: args.ImageURL,
			URLSlug:  args.URLSlug,
		}
		if err := a.accountStore(ctx).CreateAccount(account); err != nil {
			return err
		}

		au := &identity.AccountUser{
			AccountID: id,
			UserID:    args.OwnerID,
			Status:    status3.P,
			Permission: identity.Permission{
				Roles: []string{string(authorization.RoleShopOwner)},
			},
		}
		if err := a.accountUserStore(ctx).CreateAccountUser(au); err != nil {
			return err
		}
		if args.Address != nil {
			// TODO
			// create address
		}
		code, err := sqlstore.GenerateCode(ctx, tx, model.CodeTypeShop, "")
		if err != nil {
			return err
		}
		shop := &identity.Shop{
			ID:                            id,
			Name:                          args.Name,
			OwnerID:                       args.OwnerID,
			AddressID:                     args.AddressID,
			Phone:                         args.Phone,
			WebsiteURL:                    args.WebsiteURL.String,
			ImageURL:                      args.ImageURL,
			Email:                         args.Email,
			Code:                          code,
			MoneyTransactionRRule:         args.MoneyTransactionRRule,
			SurveyInfo:                    args.SurveyInfo,
			ShippingServiceSelectStrategy: args.ShippingServicePickStrategy,
			Status:                        status3.P,
			BankAccount:                   args.BankAccount,
			TryOn:                         try_on.Open,
			CompanyInfo:                   args.CompanyInfo,
		}
		if args.MoneyTransactionRRule == "" {
			// set shop MoneyTransactionRRule default value: FREQ=WEEKLY;BYDAY=TU,TH
			shop.MoneyTransactionRRule = "FREQ=WEEKLY;BYDAY=TU,TH"
		}
		if err := shop.CheckInfo(); err != nil {
			return err
		}
		if args.IsTest {
			shop.IsTest = 1
		}
		if err := a.shopStore(ctx).CreateShop(shop); err != nil {
			return err
		}

		// create for search
		// TODO: shop_search

		event := &identity.AccountCreatedEvent{
			ShopID: id,
			UserID: args.OwnerID,
		}
		return a.eventBus.Publish(ctx, event)
	})
}

func (a *Aggregate) CreateUser(ctx context.Context, args *identity.CreateUserArgs) (*identity.User, error) {
	switch args.Status {
	case status3.P:
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid status")
	}

	now := time.Now()
	userID := args.UserID
	if userID == 0 {
		userID = cm.NewIDWithTag(account_tag.TagUser)
	}
	user := &identity.User{
		ID:                      userID,
		FullName:                args.FullName,
		ShortName:               args.ShortName,
		Email:                   args.Email,
		Phone:                   args.Phone,
		Status:                  args.Status,
		WLPartnerID:             wl.GetWLPartnerID(ctx),
		Source:                  args.Source,
		PhoneVerificationSentAt: args.PhoneVerificationSentAt,
		PhoneVerifiedAt:         args.PhoneVerifiedAt,
	}
	if args.AgreeEmailInfo {
		user.AgreedEmailInfoAt = now
	}
	if args.AgreeTOS {
		user.AgreedTOSAt = now
	}
	if args.IsTest {
		user.IsTest = 1
	}

	userInternal := &identity.UserInternal{
		ID: userID,
	}
	if args.Password != "" {
		userInternal.Hashpwd = EncodePassword(args.Password)
	}
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		_, err := a.userStore(ctx).CreateUser(user)
		if xerr, ok := err.(*xerrors.APIError); ok && xerr.Err != nil {
			msg := xerr.Err.Error()
			for errKey, errMsg := range mapUserError {
				if strings.Contains(msg, errKey) {
					return cm.Errorf(cm.FailedPrecondition, nil, errMsg, wl.X(ctx).CSEmail)
				}
			}
			return err
		}

		if err := a.userInternalStore(ctx).CreateUserInternal(userInternal); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.userStore(ctx).ByID(userID).GetUser()
}

func (a *Aggregate) UpdateShopInfo(ctx context.Context, args *identity.UpdateShopInfoArgs) error {
	if args.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop ID")
	}
	update := &identity.Shop{
		MoneyTransactionRRule: args.MoneyTransactionRrule,
	}
	if err := update.CheckInfo(); err != nil {
		return err
	}
	if args.IsPriorMoneyTransaction.Valid {
		update.IsPriorMoneyTransaction = args.IsPriorMoneyTransaction
	}
	return a.shopStore(ctx).ByID(args.ShopID).UpdateShop(update)
}

func (a *Aggregate) CreateAccountUser(ctx context.Context, args *identity.CreateAccountUserArgs) (*identity.AccountUser, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	var accountUser identity.AccountUser
	if err := scheme.Convert(args, &accountUser); err != nil {
		return nil, err
	}
	if err := a.accountUserStore(ctx).CreateAccountUser(&accountUser); err != nil {
		return nil, err
	}
	return a.accountUserStore(ctx).ByUserID(args.UserID).ByAccountID(args.AccountID).GetAccountUser()
}

func (a *Aggregate) UpdateAccountUserPermission(ctx context.Context, args *identity.UpdateAccountUserPermissionArgs) error {
	if args.UserID == 0 || args.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing require params")
	}
	update := &identity.AccountUser{
		Permission: args.Permission,
	}
	if err := a.accountUserStore(ctx).ByUserID(args.UserID).ByAccountID(args.AccountID).UpdateAccountUser(update); err != nil {
		return err
	}
	return nil
}

func (a *Aggregate) DeleteAccount(ctx context.Context, args *identity.DeleteAccountArgs) error {
	if args.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing account_id")
	}
	if args.OwnerID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing owner_id")
	}

	account, err := a.accountStore(ctx).ByID(args.AccountID).ByOwnerID(args.OwnerID).GetAccount()
	if err != nil {
		return err
	}
	switch account.Type {
	case account_type.Shop:
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Does not support delete this account (type: %v)", account.Type.Name())
	}

	event := &identity.AccountDeletingEvent{
		AccountID:   args.AccountID,
		AccountType: account.Type,
	}
	if err = a.eventBus.Publish(ctx, event); err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		_, err = a.accountStore(ctx).ByID(args.AccountID).ByOwnerID(args.OwnerID).SoftDelete()
		if err != nil {
			return err
		}
		switch account.Type {
		case account_type.Shop:
			// delete shop
			if err = a.shopStore(ctx).ByID(args.AccountID).ByOwnerID(args.OwnerID).SoftDelete(); err != nil {
				return err
			}
		default:

		}

		_, err = a.accountUserStore(ctx).ByAccountID(args.AccountID).ByUserID(args.OwnerID).SoftDeleteAccountUsers()
		return err
	})
}

func (a *Aggregate) DeleteAccountUsers(ctx context.Context, args *identity.DeleteAccountUsersArgs) (int, error) {
	count := 0
	query := a.accountUserStore(ctx)
	if args.AccountID != 0 {
		query = query.ByAccountID(args.AccountID)
		count++
	}
	if args.UserID != 0 {
		query = query.ByUserID(args.UserID)
		count++
	}
	if count == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Please provide either user_id or account_id")
	}
	return query.SoftDeleteAccountUsers()
}
