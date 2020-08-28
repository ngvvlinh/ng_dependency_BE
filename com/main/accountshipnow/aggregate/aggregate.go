package aggregate

import (
	"context"
	"time"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/connectioning"
	shipnowcarrier "o.o/api/main/shipnow/carrier"
	shipnowcarriertypes "o.o/api/main/shipnow/carrier/types"
	"o.o/api/top/types/etc/connection_type"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/accountshipnow/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ accountshipnow.Aggregate = &Aggregate{}

type Aggregate struct {
	db                   *cmsql.Database
	xAccountAhamoveStore sqlstore.XAccountAhamoveStoreFactory
	eventBus             capi.EventBus
}

func NewAggregate(db com.MainDB, eventBus capi.EventBus) *Aggregate {
	return &Aggregate{
		db:                   db,
		eventBus:             eventBus,
		xAccountAhamoveStore: sqlstore.NewXAccountAhamoveStore(db),
	}
}

func AggregateMessageBus(a *Aggregate) accountshipnow.CommandBus {
	b := bus.New()
	return accountshipnow.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateExternalAccountAhamove(ctx context.Context, args *accountshipnow.CreateExternalAccountAhamoveArgs) (_result *accountshipnow.ExternalAccountAhamove, _err error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		connectionID := shipnowcarrier.GetConnectionID(args.ConnectionID, shipnowcarriertypes.Ahamove, connection_type.ConnectionMethodDirect)
		account, err := a.xAccountAhamoveStore(ctx).Phone(args.Phone).OwnerID(args.OwnerID).ConnectionID(connectionID).GetXAccountAhamove()
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}

		var id dot.ID
		phone := args.Phone
		if account == nil {
			// create new account
			id = cm.NewID()
			args1 := &sqlstore.CreateXAccountAhamoveArgs{
				ID:           id,
				OwnerID:      args.OwnerID,
				Phone:        phone,
				Name:         args.Name,
				ConnectionID: connectionID,
			}
			if _, err := a.xAccountAhamoveStore(ctx).CreateXAccountAhamove(args1); err != nil {
				return err
			}
		} else {
			id = account.ID
		}

		event := &accountshipnow.ExternalAccountAhamoveCreatedEvent{
			ID:           id,
			Phone:        phone,
			OwnerID:      args.OwnerID,
			ShopID:       args.ShopID,
			ConnectionID: connectionID,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		_result, err = a.xAccountAhamoveStore(ctx).ID(id).GetXAccountAhamove()
		return err
	})
	return _result, err
}

func (a *Aggregate) RequestVerifyExternalAccountAhamove(ctx context.Context, args *accountshipnow.RequestVerifyExternalAccountAhamoveArgs) (*accountshipnow.RequestVerifyExternalAccountAhamoveResult, error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		account, err := a.xAccountAhamoveStore(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
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

		event := &accountshipnow.ExternalAccountShipnowVerifyRequestedEvent{
			ID:           account.ID,
			ConnectionID: connectioning.DefaultTopShipAhamoveConnectionID,
			OwnerID:      args.OwnerID,
		}
		return a.eventBus.Publish(ctx, event)
	})
	return nil, err
}

func (a *Aggregate) UpdateVerifiedExternalAccountAhamove(ctx context.Context, args *accountshipnow.UpdateVerifiedExternalAccountAhamoveArgs) (*accountshipnow.ExternalAccountAhamove, error) {
	phoneNorm, ok := validate.NormalizePhone(args.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	account, err := a.xAccountAhamoveStore(ctx).Phone(phoneNorm.String()).OwnerID(args.OwnerID).GetXAccountAhamove()
	if err != nil {
		return nil, err
	}
	if account.ExternalVerified {
		return account, nil
	}

	event := &accountshipnow.ExternalAccountShipnowUpdateVerificationInfoEvent{
		ID:           account.ID,
		OwnerID:      args.OwnerID,
		ConnectionID: args.ConnectionID,
	}
	if err := a.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}
	return a.xAccountAhamoveStore(ctx).ID(account.ID).GetXAccountAhamove()
}

func (a *Aggregate) UpdateExternalAccountAhamoveVerification(ctx context.Context, args *accountshipnow.UpdateExternalAccountAhamoveVerificationArgs) (*accountshipnow.ExternalAccountAhamove, error) {
	account, err := a.xAccountAhamoveStore(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
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

	return a.xAccountAhamoveStore(ctx).UpdateVerificationImages(update)
}

func (a *Aggregate) UpdateExternalAccountAhamoveExternalInfo(ctx context.Context, args *accountshipnow.UpdateXAccountAhamoveExternalInfoArgs) error {
	return a.xAccountAhamoveStore(ctx).UpdateXAccountAhamove(args)
}

func (a *Aggregate) DeleteExternalAccountAhamove(ctx context.Context, args *accountshipnow.DeleteXAccountAhamoveArgs) error {
	return a.xAccountAhamoveStore(ctx).OwnerID(args.OwnerID).DeleteXAccountAhamove()
}
