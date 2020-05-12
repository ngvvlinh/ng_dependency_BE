package fbuser

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type FbUserAggregate struct {
	db                              *cmsql.Database
	fbUserStore                     sqlstore.FbExternalUserStoreFactory
	fbUserInternalStore             sqlstore.FbExternalUserInternalFactory
	fbPageAggr                      fbpaging.CommandBus
	customerQuery                   customering.QueryBus
	fbExternalUserShopCustomerStore sqlstore.FbExternalUserShopCustomerStoreFactory
}

func NewFbUserAggregate(
	db *cmsql.Database, fbPageA fbpaging.CommandBus, customerQ customering.QueryBus,
) *FbUserAggregate {
	return &FbUserAggregate{
		db:                              db,
		fbUserStore:                     sqlstore.NewFbExternalUserStore(db),
		fbUserInternalStore:             sqlstore.NewFbExternalUserInternalStore(db),
		fbExternalUserShopCustomerStore: sqlstore.NewFbExternalUserShopCustomerStore(db),
		fbPageAggr:                      fbPageA,

		customerQuery: customerQ,
	}
}

func FbUserAggregateMessageBus(a *FbUserAggregate) fbusering.CommandBus {
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
	fbUserResult := new(fbusering.FbExternalUser)
	fbUserInternalResult := new(fbusering.FbExternalUserInternal)
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
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

func (a *FbUserAggregate) CreateFbExternalUserShopCustomer(ctx context.Context, shopID dot.ID, externalID string, customerID dot.ID) (*fbusering.FbExternalUserWithCustomer, error) {
	var result = &fbusering.FbExternalUserWithCustomer{}
	fbUser, err := a.fbUserStore(ctx).ExternalID(externalID).GetFbExternalUser()
	if err != nil {
		return nil, err
	}
	query := &customering.GetCustomerByIDQuery{
		ID:     customerID,
		ShopID: shopID,
	}
	err = a.customerQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	fbExternalUserWithCustomer := &fbusering.FbExternalUserShopCustomer{
		ShopID:           shopID,
		FbExternalUserID: externalID,
		CustomerID:       customerID,
		Status:           fbUser.Status,
	}
	err = a.fbExternalUserShopCustomerStore(ctx).CreateFbExternalUserShopCustomer(fbExternalUserWithCustomer)
	if err != nil {
		return nil, err
	}
	result.ShopCustomer = query.Result
	result.FbExternalUser = fbUser
	return result, nil
}

func (a *FbUserAggregate) DeleteFbExternalUserShopCustomer(ctx context.Context, args *fbusering.DeleteFbExternalUserShopCustomerArgs) error {
	if args.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}
	if !args.CustomerID.Valid && !args.ExternalID.Valid {
		return cm.Errorf(cm.FailedPrecondition, nil, "Must have one in customer_id or external_id")
	}
	query := a.fbExternalUserShopCustomerStore(ctx).ShopID(args.ShopID)
	if args.CustomerID.Valid {
		query = query.ShopCustomerID(args.CustomerID.ID)
	}
	if args.ExternalID.Valid {
		query = query.FbExternalUserID(args.ExternalID.String)
	}
	return query.DeleteFbExternalUserShopCustomer()
}
