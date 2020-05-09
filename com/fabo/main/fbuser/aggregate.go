package fbuser

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type FbUserAggregate struct {
	db                  *cmsql.Database
	fbUserStore         sqlstore.FbExternalUserStoreFactory
	fbUserInternalStore sqlstore.FbExternalUserInternalFactory
	fbPageAggr          fbpaging.CommandBus
}

func NewFbUserAggregate(
	db *cmsql.Database, fbPageA fbpaging.CommandBus,
) *FbUserAggregate {
	return &FbUserAggregate{
		db:                  db,
		fbUserStore:         sqlstore.NewFbExternalUserStore(db),
		fbUserInternalStore: sqlstore.NewFbExternalUserInternalStore(db),
		fbPageAggr:          fbPageA,
	}
}

func (a *FbUserAggregate) MessageBus() fbusering.CommandBus {
	b := bus.New()
	return fbusering.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbUserAggregate) CreateFbExternalUsers(
	ctx context.Context, args *fbusering.CreateFbExternalUsersArgs,
) ([]*fbusering.FbExternalUser, error) {
	newFbExternalUsers := make([]*fbusering.FbExternalUser, 0, len(args.FbExternalUsers))
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalUser := range args.FbExternalUsers {
			newFbExternalUser := new(fbusering.FbExternalUser)
			if err := scheme.Convert(fbExternalUser, newFbExternalUser); err != nil {
				return err
			}
			newFbExternalUsers = append(newFbExternalUsers, newFbExternalUser)
		}
		if err := a.fbUserStore(ctx).CreateFbExternalUsers(newFbExternalUsers); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return newFbExternalUsers, nil
}

func (a *FbUserAggregate) CreateFbExternalUser(
	ctx context.Context, args *fbusering.CreateFbExternalUserArgs,
) (*fbusering.FbExternalUser, error) {
	fbUserResult := new(fbusering.FbExternalUser)
	if err := scheme.Convert(args, fbUserResult); err != nil {
		return nil, err
	}
	if err := a.fbUserStore(ctx).CreateFbExternalUser(fbUserResult); err != nil {
		return nil, err
	}
	return fbUserResult, nil
}

func (a *FbUserAggregate) CreateFbExternalUserInternal(
	ctx context.Context, args *fbusering.CreateFbExternalUserInternalArgs,
) (*fbusering.FbExternalUserInternal, error) {
	fbUserInternalResult := new(fbusering.FbExternalUserInternal)
	if err := scheme.Convert(args, fbUserInternalResult); err != nil {
		return nil, err
	}
	if err := a.fbUserInternalStore(ctx).CreateFbExternalUserInternal(fbUserInternalResult); err != nil {
		return nil, err
	}
	return fbUserInternalResult, nil
}

func (a *FbUserAggregate) CreateFbExternalUserCombined(
	ctx context.Context, args *fbusering.CreateFbExternalUserCombinedArgs,
) (*fbusering.FbExternalUserCombined, error) {
	oldFbUser, err := a.fbUserStore(ctx).UserID(args.UserID).ExternalID(args.FbUser.ExternalID).GetFbExternalUser()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}

	fbUserResult := new(fbusering.FbExternalUser)
	fbUserInternalResult := new(fbusering.FbExternalUserInternal)
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if oldFbUser != nil {
			args.FbUser.ID = oldFbUser.ID
			args.FbUserInternal.ID = oldFbUser.ID
		}

		// create FbExternalUser
		if err := scheme.Convert(args.FbUser, fbUserResult); err != nil {
			return err
		}
		if err := a.fbUserStore(ctx).CreateFbExternalUser(fbUserResult); err != nil {
			return err
		}

		// create FbExternalUserInternal
		if err := scheme.Convert(args.FbUserInternal, fbUserInternalResult); err != nil {
			return err
		}
		if err := a.fbUserInternalStore(ctx).CreateFbExternalUserInternal(fbUserInternalResult); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &fbusering.FbExternalUserCombined{
		FbExternalUser:         fbUserResult,
		FbExternalUserInternal: fbUserInternalResult,
	}, nil
}
