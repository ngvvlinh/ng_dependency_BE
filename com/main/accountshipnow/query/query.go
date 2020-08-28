package query

import (
	"context"

	"o.o/api/main/accountshipnow"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/accountshipnow/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	"o.o/capi"
)

var _ accountshipnow.QueryService = &QueryService{}

type QueryService struct {
	eventBus             capi.EventBus
	xAccountAhamoveStore sqlstore.XAccountAhamoveStoreFactory
}

func NewQueryService(db com.MainDB, eventBus capi.EventBus) *QueryService {
	return &QueryService{
		eventBus:             eventBus,
		xAccountAhamoveStore: sqlstore.NewXAccountAhamoveStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) accountshipnow.QueryBus {
	b := bus.New()
	h := accountshipnow.NewQueryServiceHandler(q)
	return h.RegisterHandlers(b)
}

func (q *QueryService) GetExternalAccountAhamove(ctx context.Context, args *accountshipnow.GetExternalAccountAhamoveArgs) (*accountshipnow.ExternalAccountAhamove, error) {
	phone := args.Phone
	return q.xAccountAhamoveStore(ctx).Phone(phone).OwnerID(args.OwnerID).GetXAccountAhamove()
}

func (q *QueryService) GetExternalAccountAhamoveByExternalID(ctx context.Context, args *accountshipnow.GetExternalAccountAhamoveByExternalIDQueryArgs) (*accountshipnow.ExternalAccountAhamove, error) {
	return q.xAccountAhamoveStore(ctx).ExternalID(args.ExternalID).GetXAccountAhamove()
}

func (q *QueryService) GetAccountShipnow(ctx context.Context, args *accountshipnow.GetAccountShipnowArgs) (*accountshipnow.ExternalAccountAhamove, error) {
	phoneNorm, ok := validate.NormalizePhone(args.Phone)
	if !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin connection ID")
	}

	account, err := q.xAccountAhamoveStore(ctx).Phone(phoneNorm.String()).OwnerID(args.OwnerID).ConnectionID(args.ConnectionID).GetXAccountAhamove()
	if err != nil {
		return nil, err
	}

	if !account.ExternalVerified {
		// Tài khoản chưa được verify
		// Gọi qua NVC để cập nhật
		event := &accountshipnow.ExternalAccountShipnowUpdateVerificationInfoEvent{
			ID:           account.ID,
			OwnerID:      account.OwnerID,
			ConnectionID: account.ConnectionID,
		}
		if err := q.eventBus.Publish(ctx, event); err != nil {
			return nil, err
		}
		return q.xAccountAhamoveStore(ctx).ID(account.ID).GetXAccountAhamove()
	}
	return account, nil
}
