package identity

import (
	"context"
	"encoding/json"
	"time"

	"etop.vn/common/bus"

	"etop.vn/api/main/shipnow/carrier"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/identity/sqlstore"
)

var _ identity.Aggregate = &Aggregate{}

type Aggregate struct {
	db                    cmsql.Transactioner
	shopStore             sqlstore.ShopStoreFactory
	userStore             sqlstore.UserStoreFactory
	xAccountAhamove       sqlstore.XAccountAhamoveStoreFactory
	shipnowCarrierManager carrier.Manager
}

func NewAggregate(db cmsql.Database, carrierManager carrier.Manager) *Aggregate {
	return &Aggregate{
		db:                    db,
		shopStore:             sqlstore.NewIdentityStore(db),
		xAccountAhamove:       sqlstore.NewXAccountAhamoveStore(db),
		userStore:             sqlstore.NewUserStore(db),
		shipnowCarrierManager: carrierManager,
	}
}

func (a *Aggregate) MessageBus() identity.CommandBus {
	b := bus.New()
	return identity.NewAggregateHandler(a).RegisterHandlers(b)
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

		var id int64
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
			Carrier: carrier.Ahamove,
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
			Carrier: carrier.Ahamove,
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
			Carrier: carrier.Ahamove,
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
			Carrier: carrier.Ahamove,
		}
		res, err := a.shipnowCarrierManager.VerifyExternalAccount(ctx, args2)
		if err != nil {
			return err
		}
		// update external_ticket_id
		externalData, _ := json.Marshal(res)
		update := &sqlstore.UpdateXAccountAhamoveVerifiedInfoArgs{
			ID:                    account.ID,
			ExternalTickerID:      res.TicketID,
			LastSendVerifiedAt:    time.Now(),
			ExternaleDataVerified: externalData,
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
		Carrier: carrier.Ahamove,
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
