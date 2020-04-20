package fbuser

import (
	"context"

	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/fabo/fbusering"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/fabo/main/fbuser/convert"
	"etop.vn/backend/com/fabo/main/fbuser/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/common/l"
)

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type FbUserAggregate struct {
	db                  *cmsql.Database
	fbUserStore         sqlstore.FbUserStoreFactory
	fbUserInternalStore sqlstore.FbUserInternalFactory
	fbPageAggr          fbpaging.CommandBus
}

func NewFbUserAggregate(
	db *cmsql.Database, fbPageA fbpaging.CommandBus,
) *FbUserAggregate {
	return &FbUserAggregate{
		db:                  db,
		fbUserStore:         sqlstore.NewFbUserStore(db),
		fbUserInternalStore: sqlstore.NewFbUserInternalStore(db),
		fbPageAggr:          fbPageA,
	}
}

func (a *FbUserAggregate) MessageBus() fbusering.CommandBus {
	b := bus.New()
	return fbusering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbUserAggregate) CreateFbUser(
	ctx context.Context, args *fbusering.CreateFbUserArgs,
) (*fbusering.FbUser, error) {
	fbUserResult := new(fbusering.FbUser)
	if err := scheme.Convert(args, fbUserResult); err != nil {
		return nil, err
	}
	if err := a.fbUserStore(ctx).CreateFbUser(fbUserResult); err != nil {
		return nil, err
	}
	return fbUserResult, nil
}

func (a *FbUserAggregate) CreateFbUserInternal(
	ctx context.Context, args *fbusering.CreateFbUserInternalArgs,
) (*fbusering.FbUserInternal, error) {
	fbUserInternalResult := new(fbusering.FbUserInternal)
	if err := scheme.Convert(args, fbUserInternalResult); err != nil {
		return nil, err
	}
	if err := a.fbUserInternalStore(ctx).CreateFbUserInternal(fbUserInternalResult); err != nil {
		return nil, err
	}
	return fbUserInternalResult, nil
}

func (a *FbUserAggregate) CreateFbUserCombined(
	ctx context.Context, args *fbusering.CreateFbUserCombinedArgs,
) (*fbusering.FbUserCombined, error) {
	oldFbUser, err := a.fbUserStore(ctx).GetFbUser()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}

	fbUserResult := new(fbusering.FbUser)
	fbUserInternalResult := new(fbusering.FbUserInternal)
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if oldFbUser != nil {
			if oldFbUser.ExternalID != args.FbUser.ExternalID {
				_, err := a.fbUserStore(ctx).ExternalID(oldFbUser.ExternalID).UpdateStatus(int(status3.N))
				if err != nil {
					return err
				}

				// disable all fbPages of old FbUser
				disableAllFbPagesCmd := &fbpaging.DisableAllFbPagesCommand{
					ShopID: args.ShopID,
					UserID: args.UserID,
				}
				if err := a.fbPageAggr.Dispatch(ctx, disableAllFbPagesCmd); err != nil {
					return err
				}
			} else {
				args.FbUser.ID = oldFbUser.ID
				args.FbUserInternal.ID = oldFbUser.ID
			}
		}

		// create FbUser
		if err := scheme.Convert(args.FbUser, fbUserResult); err != nil {
			return err
		}
		if err := a.fbUserStore(ctx).CreateFbUser(fbUserResult); err != nil {
			return err
		}

		// create FbUserInternal
		if err := scheme.Convert(args.FbUserInternal, fbUserInternalResult); err != nil {
			return err
		}
		if err := a.fbUserInternalStore(ctx).CreateFbUserInternal(fbUserInternalResult); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &fbusering.FbUserCombined{
		FbUser:         fbUserResult,
		FbUserInternal: fbUserInternalResult,
	}, nil
}
