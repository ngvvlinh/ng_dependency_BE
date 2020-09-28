package identity

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/com/main/identity/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

var _ identity.Aggregate = &Aggregate{}

type Aggregate struct {
	db               *cmsql.Database
	userStore        sqlstore.UserStoreFactory
	affiliateStore   sqlstore.AffiliateStoreFactory
	shopStore        sqlstore.ShopStoreFactory
	accountStore     sqlstore.AccountStoreFactory
	accountUserStore sqlstore.AccountUserStoreFactory
	userRefSaffStore sqlstore.UserRefSaffStoreFactory
}

func NewAggregate(db com.MainDB, eventBus capi.EventBus) *Aggregate {
	return &Aggregate{
		db:               db,
		userStore:        sqlstore.NewUserStore(db),
		shopStore:        sqlstore.NewShopStore(db),
		affiliateStore:   sqlstore.NewAffiliateStore(db),
		accountStore:     sqlstore.NewAccountStore(db),
		accountUserStore: sqlstore.NewAccountUserStore(db),
		userRefSaffStore: sqlstore.NewUserRefSaffStore(db),
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
	if currentUser, err := a.userStore(ctx).ByID(args.UserID).GetUserDB(ctx); err != nil {
		return err
	} else if currentUser.RefSaleID != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "RefUserID đã tồn tại. Không thể cập nhật.")
	}
	refUser, err := a.userStore(ctx).ByPhone(args.RefUserPhone).GetUserDB(ctx)
	if err != nil {
		return cm.Errorf(cm.NotFound, nil, "Số điện thoại người dùng không tồn tại.")
	}

	updateCmd := &sqlstore.UpdateRefferenceIDArgs{
		UserID:    args.UserID,
		RefUserID: refUser.ID,
	}
	return a.userStore(ctx).UpdateUserRefferenceID(updateCmd)
}

func (a *Aggregate) UpdateUserReferenceSaleID(ctx context.Context, args *identity.UpdateUserReferenceSaleIDArgs) error {
	if currentUser, err := a.userStore(ctx).ByID(args.UserID).GetUserDB(ctx); err != nil {
		return err
	} else if currentUser.RefSaleID != 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "RefSaleID đã tồn tại. Không thể cập nhật.")
	}
	refUser, err := a.userStore(ctx).ByPhone(args.RefSalePhone).GetUserDB(ctx)
	if err != nil {
		return cm.Errorf(cm.NotFound, nil, "Số điện thoại người dùng không tồn tại.")
	}

	updateCmd := &sqlstore.UpdateRefferenceIDArgs{
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
	_args := sqlstore.CreateAffiliateArgs{
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
	_args := sqlstore.UpdateAffiliateArgs{
		ID:      args.ID,
		OwnerID: args.OwnerID,
		Phone:   phoneNorm.String(),
		Email:   emailNorm.String(),
		Name:    args.Name,
	}
	return a.affiliateStore(ctx).UpdateAffiliate(_args)
}

func (a *Aggregate) UpdateAffiliateBankAccount(ctx context.Context, args *identity.UpdateAffiliateBankAccountArgs) (*identity.Affiliate, error) {
	_args := sqlstore.UpdateAffiliateArgs{
		ID:          args.ID,
		OwnerID:     args.OwnerID,
		BankAccount: args.BankAccount,
	}
	return a.affiliateStore(ctx).UpdateAffiliate(_args)
}

func (a *Aggregate) DeleteAffiliate(ctx context.Context, args *identity.DeleteAffiliateArgs) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		args1 := sqlstore.DeleteAffiliateArgs{
			ID:      args.ID,
			OwnerID: args.OwnerID,
		}
		if err := a.affiliateStore(ctx).DeleteAffiliate(args1); err != nil {
			return err
		}
		args2 := sqlstore.DeleteAccountUserArgs{
			AccountID: args.ID,
			UserID:    args.OwnerID,
		}
		return a.accountUserStore(ctx).DeleteAccountUser(args2)
	})
}

func (a *Aggregate) BlockUser(ctx context.Context, args *identity.BlockUserArgs) (*identity.User, error) {
	user, err := a.userStore(ctx).ByID(args.UserID).GetUserDB(ctx)
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
	user, err = a.userStore(ctx).ByID(args.UserID).GetUserDB(ctx)
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
	user, err := a.userStore(ctx).ByID(userID).GetUserDB(ctx)
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
	user, err = a.userStore(ctx).ByID(userID).GetUserDB(ctx)
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
	_, err := a.userStore(ctx).ByID(args.UserID).GetUserDB(ctx)
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
