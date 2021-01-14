package fbuser

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var _ fbusering.Aggregate = &FbUserAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type FbUserAggregate struct {
	db                              *cmsql.Database
	fbShopUserTagStore              sqlstore.FbShopTagStoreFactory
	fbUserStore                     sqlstore.FbExternalUserStoreFactory
	fbUserConnectedStore            sqlstore.FbExternalUserConnectedStoreFactory
	fbUserInternalStore             sqlstore.FbExternalUserInternalFactory
	fbPageAggr                      fbpaging.CommandBus
	customerQuery                   customering.QueryBus
	fbExternalUserShopCustomerStore sqlstore.FbExternalUserShopCustomerStoreFactory
}

func NewFbUserAggregate(
	db com.MainDB, fbPageA fbpaging.CommandBus, customerQ customering.QueryBus,
) *FbUserAggregate {
	return &FbUserAggregate{
		db:                              db,
		fbShopUserTagStore:              sqlstore.NewFbShopTagStore(db),
		fbUserStore:                     sqlstore.NewFbExternalUserStore(db),
		fbUserConnectedStore:            sqlstore.NewFbExternalUserConnectedStore(db),
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

func (a *FbUserAggregate) CreateOrUpdateFbExternalUserCombined(
	ctx context.Context, args *fbusering.CreateOrUpdateFbExternalUserCombinedArgs,
) (*fbusering.FbExternalUserCombined, error) {

	if err := args.Validate(); err != nil {
		return nil, err
	}

	fbUserConnectedResult := new(fbusering.FbExternalUserConnected)
	if err := scheme.Convert(args.FbUserConnected, fbUserConnectedResult); err != nil {
		return nil, err
	}

	fbUserInternalResult := new(fbusering.FbExternalUserInternal)
	if err := scheme.Convert(args.FbUserInternal, fbUserInternalResult); err != nil {
		return nil, err
	}

	shopID := fbUserConnectedResult.ShopID

	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// disable all users are same shop_id
		if _, err := a.fbUserConnectedStore(ctx).ShopID(shopID).UpdateStatus(int(status3.N)); err != nil {
			return err
		}

		// create or update FbExternalUserConnected
		if err := a.fbUserConnectedStore(ctx).CreateOrUpdateFbExternalUserConnected(fbUserConnectedResult); err != nil {
			return err
		}

		// create or update FbExternalUserInternal
		if err := a.fbUserInternalStore(ctx).CreateOrUpdateFbExternalUserInternal(fbUserInternalResult); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &fbusering.FbExternalUserCombined{
		FbExternalUserConnected: fbUserConnectedResult,
		FbExternalUserInternal:  fbUserInternalResult,
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

func (a *FbUserAggregate) CreateShopUserTag(ctx context.Context, args *fbusering.CreateShopUserTagArgs) (*fbusering.FbShopUserTag, error) {
	tag := new(fbusering.FbShopUserTag)
	if err := scheme.Convert(args, tag); err != nil {
		return nil, err
	}
	if err := a.fbShopUserTagStore(ctx).CreateShopUserTag(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (a *FbUserAggregate) UpdateShopUserTag(ctx context.Context, args *fbusering.UpdateShopUserTagArgs) (*fbusering.FbShopUserTag, error) {
	tag := &fbusering.FbShopUserTag{}
	tag = convert.Apply_fbusering_UpdateShopUserTagArgs_fbusering_FbShopUserTag(args, tag)
	if err := a.fbShopUserTagStore(ctx).ByID(args.ID).UpdateShopUserTag(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (a *FbUserAggregate) DeleteShopUserTag(ctx context.Context, args *fbusering.DeleteShopUserTagArgs) (int, error) {
	err := a.fbShopUserTagStore(ctx).ByID(args.ID).ByShopID(args.ShopID).DeleteShopUserTag()
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (a *FbUserAggregate) UpdateShopUserTags(ctx context.Context, args *fbusering.UpdateShopUserTagsArgs) (*fbusering.FbExternalUser, error) {
	args.TagIDs = removeDuplicateTagID(args.TagIDs)
	tags, err := a.fbShopUserTagStore(ctx).ByIDs(args.TagIDs).GetShopUserTags()
	if err != nil {
		return nil, err
	}

	var newUserTagIds []dot.ID
	for _, tag := range tags {
		newUserTagIds = append(newUserTagIds, tag.ID)
	}

	extUser, err := a.fbUserStore(ctx).GetFbExternalUserByShopID(args.FbExternalUserID, args.ShopID)
	if err != nil {
		return nil, err
	}

	err = a.fbUserStore(ctx).ExternalID(extUser.ExternalID).UpdateUserTags(newUserTagIds)
	if err != nil {
		return nil, err
	}

	extUser.TagIDs = newUserTagIds
	return extUser, nil
}

func removeDuplicateTagID(userTags []dot.ID) []dot.ID {
	_m := map[dot.ID]struct{}{}
	var result []dot.ID
	for _, id := range userTags {
		if _, ok := _m[id]; ok {
			continue
		}
		result = append(result, id)
		_m[id] = struct{}{}
	}
	return result
}
