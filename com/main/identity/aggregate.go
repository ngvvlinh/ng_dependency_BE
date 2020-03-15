package identity

import (
	"context"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/shipnow/carrier"
	carriertypes "etop.vn/api/main/shipnow/carrier/types"
	"etop.vn/backend/com/main/identity/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

var _ identity.Aggregate = &Aggregate{}

type Aggregate struct {
	db                    cmsql.Transactioner
	userStore             sqlstore.UserStoreFactory
	accountStore          sqlstore.AccountStoreFactory
	accountUserStore      sqlstore.AccountUserStoreFactory
	xAccountAhamove       sqlstore.XAccountAhamoveStoreFactory
	shipnowCarrierManager carrier.Manager
}

func NewAggregate(db *cmsql.Database, carrierManager carrier.Manager) *Aggregate {
	return &Aggregate{
		db:                    db,
		xAccountAhamove:       sqlstore.NewXAccountAhamoveStore(db),
		userStore:             sqlstore.NewUserStore(db),
		accountStore:          sqlstore.NewAccountStore(db),
		accountUserStore:      sqlstore.NewAccoutnUserStore(db),
		shipnowCarrierManager: carrierManager,
	}
}

func (a *Aggregate) MessageBus() identity.CommandBus {
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

func (a *Aggregate) CreateExternalAccountAhamove(ctx context.Context, args *identity.CreateExternalAccountAhamoveArgs) (_result *identity.ExternalAccountAhamove, _err error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		account, err := a.xAccountAhamove(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}

		if account != nil && account.ExternalToken != "" {
			_result = account
			return nil
		}

		var id dot.ID
		phone := args.Phone
		if account == nil {
			// create new account
			id = cm.NewID()
			args1 := &sqlstore.CreateXAccountAhamoveArgs{
				ID:      id,
				OwnerID: args.OwnerID,
				Phone:   phone,
				Name:    args.Name,
			}
			if _, err := a.xAccountAhamove(ctx).CreateXAccountAhamove(args1); err != nil {
				return err
			}
		} else {
			id = account.ID
		}

		// Ahamove register account

		args2 := &carrier.RegisterExternalAccountCommand{
			Phone:   phone,
			Name:    args.Name,
			Address: args.Address,
			Carrier: carriertypes.Ahamove,
		}
		regisResult, err := a.shipnowCarrierManager.RegisterExternalAccount(ctx, args2)
		if err != nil {
			return err
		}

		// Re-get external account and Update external info
		// First, update token
		_, err = a.xAccountAhamove(ctx).UpdateXAccountAhamove(&sqlstore.UpdateXAccountAhamoveInfoArgs{
			ID:            id,
			ExternalToken: regisResult.Token,
		})
		if err != nil {
			return err
		}

		xAccount, err := a.shipnowCarrierManager.GetExternalAccount(ctx, &carrier.GetExternalAccountCommand{
			OwnerID: args.OwnerID,
			Carrier: carriertypes.Ahamove,
		})
		if err != nil {
			return err
		}
		// update another info
		args3 := &sqlstore.UpdateXAccountAhamoveInfoArgs{
			ID:                id,
			ExternalID:        xAccount.ID,
			ExternalCreatedAt: xAccount.CreatedAt,
		}
		_result, err = a.xAccountAhamove(ctx).UpdateXAccountAhamove(args3)
		return err
	})

	return _result, err
}

func (a *Aggregate) RequestVerifyExternalAccountAhamove(ctx context.Context, args *identity.RequestVerifyExternalAccountAhamoveArgs) (*identity.RequestVerifyExternalAccountAhamoveResult, error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		account, err := a.xAccountAhamove(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
		if err != nil {
			return err
		}

		if account.ExternalVerified {
			return cm.Errorf(cm.FailedPrecondition, nil, "Tài khoản đã được xác thực.")
		}

		now := time.Now()
		lastSendVerify := account.LastSendVerifiedAt
		if now.Sub(lastSendVerify) < 2*24*time.Hour {
			return cm.Errorf(cm.FailedPrecondition, nil, "Yêu cầu xác thực tài khoản đã được gửi và sẽ được xử lý trong vòng 2 ngày làm việc kể từ ngày đăng ký xác thực. Vui lòng chờ.")
		}

		// check external account ahamove verification images
		if account.IDCardFrontImg == "" || account.IDCardBackImg == "" || account.PortraitImg == "" {
			return cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin người dùng. Vui lòng cung cấp ảnh 2 mặt CMND và 1 ảnh chân dung.")
		}

		args1 := &carrier.GetExternalAccountCommand{
			OwnerID: args.OwnerID,
			Carrier: carriertypes.Ahamove,
		}
		xAccount, err := a.shipnowCarrierManager.GetExternalAccount(ctx, args1)
		if err != nil {
			return err
		}
		if xAccount.Verified {
			update := &sqlstore.UpdateXAccountAhamoveInfoArgs{
				ID:               account.ID,
				ExternalID:       xAccount.ID,
				ExternalVerified: xAccount.Verified,
			}
			if _, err := a.xAccountAhamove(ctx).UpdateXAccountAhamove(update); err != nil {
				return err
			}
			return nil
		}

		// send verify request to Ahamove
		args2 := &carrier.VerifyExternalAccountCommand{
			OwnerID: args.OwnerID,
			Carrier: carriertypes.Ahamove,
		}
		res, err := a.shipnowCarrierManager.VerifyExternalAccount(ctx, args2)
		if err != nil {
			return err
		}
		// update external_ticket_id
		externalData, _ := jsonx.Marshal(res)
		update := &sqlstore.UpdateXAccountAhamoveVerifiedInfoArgs{
			ID:                   account.ID,
			ExternalTickerID:     res.TicketID,
			LastSendVerifiedAt:   time.Now(),
			ExternalDataVerified: externalData,
		}
		_, err = a.xAccountAhamove(ctx).UpdateXAccountAhamoveVerifiedInfo(update)
		return err
	})
	return nil, err
}

func (a *Aggregate) UpdateVerifiedExternalAccountAhamove(ctx context.Context, args *identity.UpdateVerifiedExternalAccountAhamoveArgs) (*identity.ExternalAccountAhamove, error) {
	account, err := a.xAccountAhamove(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
	if err != nil {
		return nil, err
	}
	if account.ExternalVerified {
		return account, nil
	}

	xAccount, err := a.shipnowCarrierManager.GetExternalAccount(ctx, &carrier.GetExternalAccountCommand{
		OwnerID: args.OwnerID,
		Carrier: carriertypes.Ahamove,
	})
	if err != nil {
		return nil, err
	}

	if !xAccount.Verified {
		return account, nil
	}

	update := &sqlstore.UpdateXAccountAhamoveVerifiedInfoArgs{
		ID:               account.ID,
		ExternalVerified: xAccount.Verified,
	}
	return a.xAccountAhamove(ctx).UpdateXAccountAhamoveVerifiedInfo(update)
}

func (a *Aggregate) UpdateExternalAccountAhamoveVerification(ctx context.Context, args *identity.UpdateExternalAccountAhamoveVerificationArgs) (*identity.ExternalAccountAhamove, error) {
	account, err := a.xAccountAhamove(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
	if err != nil {
		return nil, err
	}
	if account.ExternalVerified {
		return account, nil
	}

	update := &sqlstore.UpdateXAccountAhamoveVerificationImageArgs{
		ID:                  account.ID,
		IDCardFrontImg:      args.IDCardFrontImg,
		IDCardBackImg:       args.IDCardBackImg,
		PortraitImg:         args.PortraitImg,
		WebsiteURL:          args.WebsiteURL,
		FanpageURL:          args.FanpageURL,
		CompanyImgs:         args.CompanyImgs,
		BusinessLicenseImgs: args.BusinessLicenseImgs,
	}

	return a.xAccountAhamove(ctx).UpdateVerificationImages(update)
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
	return a.accountStore(ctx).CreateAffiliate(_args)
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
	return a.accountStore(ctx).UpdateAffiliate(_args)
}

func (a *Aggregate) UpdateAffiliateBankAccount(ctx context.Context, args *identity.UpdateAffiliateBankAccountArgs) (*identity.Affiliate, error) {
	_args := sqlstore.UpdateAffiliateArgs{
		ID:          args.ID,
		OwnerID:     args.OwnerID,
		BankAccount: args.BankAccount,
	}
	return a.accountStore(ctx).UpdateAffiliate(_args)
}

func (a *Aggregate) DeleteAffiliate(ctx context.Context, args *identity.DeleteAffiliateArgs) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		args1 := sqlstore.DeleteAffiliateArgs{
			ID:      args.ID,
			OwnerID: args.OwnerID,
		}
		if err := a.accountStore(ctx).DeleteAffiliate(args1); err != nil {
			return err
		}
		args2 := sqlstore.DeleteAccountUserArgs{
			AccountID: args.ID,
			UserID:    args.OwnerID,
		}
		return a.accountUserStore(ctx).DeleteAccountUser(args2)
	})
}
