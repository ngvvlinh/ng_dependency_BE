package identity

import (
	"context"
	"time"

	"etop.vn/backend/pkg/common/bus"

	"etop.vn/api/main/shipnow/carrier"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/main/identity"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/identity/sqlstore"
)

var _ identity.Aggregate = &Aggregate{}

type Aggregate struct {
	db                    cmsql.Transactioner
	store                 sqlstore.IdentityStoreFactory
	xAccountAhamove       sqlstore.XAccountAhamoveStoreFactory
	shipnowCarrierManager carrier.Manager
}

func NewAggregate(db cmsql.Database, carrierManager carrier.Manager) *Aggregate {
	return &Aggregate{
		db:                    db,
		store:                 sqlstore.NewIdentityStore(db),
		xAccountAhamove:       sqlstore.NewXAccountAhamoveStore(db),
		shipnowCarrierManager: carrierManager,
	}
}

func (a *Aggregate) MessageBus() identity.CommandBus {
	b := bus.New()
	return identity.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateExternalAccountAhamove(ctx context.Context, args *identity.CreateExternalAccountAhamoveArgs) (_result *identity.ExternalAccountAhamove, _err error) {
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var id int64
		account, err := a.xAccountAhamove(ctx).Phone(args.Phone).OwnerID(args.OwnerID).GetXAccountAhamove()
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}

		if account != nil && account.ExternalToken != "" {
			_result = account
			return nil
		}

		if account == nil {
			// create new account
			id = cm.NewID()
			args1 := &sqlstore.CreateXAccountAhamoveArgs{
				ID:      id,
				OwnerID: args.OwnerID,
				Phone:   args.Phone,
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
			Phone:   args.Phone,
			Name:    args.Name,
			Carrier: carrier.Ahamove,
		}
		xAccount, err := a.shipnowCarrierManager.RegisterExternalAccount(ctx, args2)
		if err != nil {
			return err
		}
		if xAccount.Token == "" {
			return cm.Errorf(cm.ExternalServiceError, nil, "Tài khoản Ahamove không hợp lệ.")
		}

		// Update token
		args3 := &sqlstore.UpdateXAccountAhamoveInfoArgs{
			ID:                id,
			ExternalToken:     xAccount.Token,
			ExternalCreatedAt: time.Now(),
		}
		_result, err = a.xAccountAhamove(ctx).UpdateXAccountAhamove(args3)
		return err
	})

	return _result, err
}
