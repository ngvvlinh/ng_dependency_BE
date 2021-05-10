package fbuser

import (
	"context"
	"fmt"

	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	"o.o/backend/com/fabo/pkg/fbclient"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

var _ fbusering.QueryService = &FbUserQuery{}

type FbUserQuery struct {
	db *cmsql.Database

	fbShopUserTagStore              sqlstore.FbShopTagStoreFactory
	fbUserStore                     sqlstore.FbExternalUserStoreFactory
	fbUserConnectedStore            sqlstore.FbExternalUserConnectedStoreFactory
	fbUserInternalStore             sqlstore.FbExternalUserInternalFactory
	fbExternalUserShopCustomerStore sqlstore.FbExternalUserShopCustomerStoreFactory
	customerQuery                   customering.QueryBus
	defaultAvatarLink               string
}

func NewFbUserQuery(database com.MainDB, customerQ customering.QueryBus) *FbUserQuery {
	_defaultAvatarLink := fmt.Sprintf("%s/%s", cm.MainSiteBaseURL(), "dl/fabo/default_avatar.png")
	return &FbUserQuery{
		db:                              database,
		fbShopUserTagStore:              sqlstore.NewFbShopTagStore(database),
		fbUserStore:                     sqlstore.NewFbExternalUserStore(database),
		fbUserInternalStore:             sqlstore.NewFbExternalUserInternalStore(database),
		fbUserConnectedStore:            sqlstore.NewFbExternalUserConnectedStore(database),
		fbExternalUserShopCustomerStore: sqlstore.NewFbExternalUserShopCustomerStore(database),
		customerQuery:                   customerQ,
		defaultAvatarLink:               _defaultAvatarLink,
	}
}

func FbUserQueryMessageBus(q *FbUserQuery) fbusering.QueryBus {
	b := bus.New()
	return fbusering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbUserQuery) GetFbExternalUserInternalByExternalID(ctx context.Context, externalID string) (*fbusering.FbExternalUserInternal, error) {
	fbUser, err := q.fbUserInternalStore(ctx).ExternalID(externalID).GetFbExternalUserInternal()
	if err != nil {
		return nil, err
	}
	return fbUser, nil
}

func (q *FbUserQuery) GetFbExternalUserByExternalID(
	ctx context.Context, externalID string,
) (*fbusering.FbExternalUser, error) {
	fbUser, err := q.fbUserStore(ctx).ExternalID(externalID).Status(status3.P).GetFbExternalUser()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatar(fbUser, q.defaultAvatarLink)

	return fbUser, nil
}

func (q *FbUserQuery) GetFbExternalUserConnectedByShopID(
	ctx context.Context, shopID dot.ID,
) (*fbusering.FbExternalUserConnected, error) {
	return q.fbUserConnectedStore(ctx).ShopID(shopID).Status(status3.P).GetFbExternalUserConnected()
}

func (q *FbUserQuery) GetFbExternalUserConnectedByExternalID(
	ctx context.Context, externalID string,
) (*fbusering.FbExternalUserConnected, error) {
	return q.fbUserConnectedStore(ctx).ExternalID(externalID).GetFbExternalUserConnected()
}

func (q *FbUserQuery) ListFbExternalUsersByExternalIDs(
	ctx context.Context, externalIDs filter.Strings, externalPageID dot.NullString,
) ([]*fbusering.FbExternalUser, error) {
	query := q.fbUserStore(ctx).ExternalIDs(externalIDs)
	if externalPageID.Valid {
		query = query.ExternalPageID(externalPageID.String)
	}
	fbUsers, err := query.ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatars(fbUsers, q.defaultAvatarLink)

	return fbUsers, nil
}

func (q *FbUserQuery) ListFbExternalUsers(ctx context.Context, args *fbusering.ListFbExternalUsersArgs) ([]*fbusering.FbExternalUserWithCustomer, error) {
	query := q.fbExternalUserShopCustomerStore(ctx).ShopID(args.ShopID)
	if args.CustomerID.Valid {
		query = query.ShopCustomerID(args.CustomerID.ID)
	}
	fbUserCustomer, err := query.ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	var fbUserIDs []string
	for _, v := range fbUserCustomer {
		fbUserIDs = append(fbUserIDs, v.FbExternalUserID)
	}
	fbUsers, err := q.fbUserStore(ctx).ExternalIDs(fbUserIDs).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatars(fbUsers, q.defaultAvatarLink)

	return q.populateFbExternalUsersWithCustomerInfo(ctx, args.ShopID, fbUsers)
}

func (q *FbUserQuery) GetFbExternalUserWithCustomerByExternalID(ctx context.Context, shopID dot.ID, externalID string) (*fbusering.FbExternalUserWithCustomer, error) {
	fbExternalUser, err := q.fbUserStore(ctx).ExternalID(externalID).GetFbExternalUser()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatar(fbExternalUser, q.defaultAvatarLink)

	var result = &fbusering.FbExternalUserWithCustomer{
		FbExternalUser: fbExternalUser,
	}

	fbExternalUserShopCustomer, err := q.fbExternalUserShopCustomerStore(ctx).ShopID(shopID).FbExternalUserID(externalID).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}

	if len(fbExternalUserShopCustomer) > 0 {
		query := &customering.GetCustomerByIDQuery{
			ID:             fbExternalUserShopCustomer[0].CustomerID,
			ShopID:         shopID,
			IncludeDeleted: true,
		}
		err = q.customerQuery.Dispatch(ctx, query)
		if err != nil {
			return nil, err
		}
		result.ShopCustomer = query.Result
	}
	return result, nil
}

func (q *FbUserQuery) ListFbExternalUserWithCustomer(ctx context.Context, args fbusering.ListFbExternalUserWithCustomerRequest) ([]*fbusering.FbExternalUserWithCustomer, error) {
	fbUsersCustomers, err := q.fbExternalUserShopCustomerStore(ctx).WithPaging(args.Paging).Filters(args.Filters).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	var fbUsersIDs []string
	for _, v := range fbUsersCustomers {
		fbUsersIDs = append(fbUsersIDs, v.FbExternalUserID)
	}
	fbUsers, err := q.fbUserStore(ctx).ExternalIDs(fbUsersIDs).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatars(fbUsers, q.defaultAvatarLink)

	return q.populateFbExternalUsersWithCustomerInfo(ctx, args.ShopID, fbUsers)
}

func (q *FbUserQuery) ListFbExternalUserWithCustomerByExternalIDs(ctx context.Context, shopID dot.ID, externalID []string) ([]*fbusering.FbExternalUserWithCustomer, error) {
	fbUsersCustomers, err := q.fbExternalUserShopCustomerStore(ctx).FbExternalUserIDs(externalID).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	var fbUsersIDs []string
	for _, v := range fbUsersCustomers {
		fbUsersIDs = append(fbUsersIDs, v.FbExternalUserID)
	}
	fbUsers, err := q.fbUserStore(ctx).ExternalIDs(fbUsersIDs).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatars(fbUsers, q.defaultAvatarLink)

	return q.populateFbExternalUsersWithCustomerInfo(ctx, shopID, fbUsers)
}

func (q *FbUserQuery) populateFbExternalUsersWithCustomerInfo(ctx context.Context, shopID dot.ID, fbUsers []*fbusering.FbExternalUser) ([]*fbusering.FbExternalUserWithCustomer, error) {
	var result []*fbusering.FbExternalUserWithCustomer
	if len(fbUsers) == 0 {
		return result, nil
	}

	var fbUserIDS []string
	for _, v := range fbUsers {
		fbUserIDS = append(fbUserIDS, v.ExternalID)
	}

	fbUserWithCustomer, err := q.fbExternalUserShopCustomerStore(ctx).FbExternalUserIDs(fbUserIDS).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	var mapFbUserCustomer = make(map[string]dot.ID)
	var customerIDs []dot.ID
	for _, v := range fbUserWithCustomer {
		if v.CustomerID != 0 {
			customerIDs = append(customerIDs, v.CustomerID)
		}
		mapFbUserCustomer[v.FbExternalUserID] = v.CustomerID
	}
	query := &customering.ListCustomersByIDsQuery{
		IDs:    customerIDs,
		ShopID: shopID,
	}
	err = q.customerQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapCustomer = make(map[dot.ID]*customering.ShopCustomer)
	for _, v := range query.Result.Customers {
		mapCustomer[v.ID] = v
	}
	for k, v := range fbUsers {
		if mapFbUserCustomer[v.ExternalID] != 0 {
			result = append(result, &fbusering.FbExternalUserWithCustomer{
				FbExternalUser: fbUsers[k],
				ShopCustomer:   mapCustomer[mapFbUserCustomer[fbUsers[k].ExternalID]],
			})

		}
	}
	return result, nil
}

func (q *FbUserQuery) ListShopCustomerIDWithPhoneNorm(
	ctx context.Context,
	shopID dot.ID,
	phone string,
) ([]dot.ID, error) {
	query := &customering.ListCustomersByPhoneNormQuery{
		ShopID: shopID,
		Phone:  phone,
	}
	err := q.customerQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}

	customers := query.Result
	var result []dot.ID
	for _, customer := range customers {
		result = append(result, customer.ID)
	}

	return result, nil
}

func (q *FbUserQuery) ListShopCustomerWithFbExternalUser(ctx context.Context, args *fbusering.ListCustomerWithFbAvatarsArgs) (*fbusering.ListShopCustomerWithFbExternalUserResponse, error) {
	query := &customering.ListCustomersQuery{
		ShopID:  args.ShopID,
		Paging:  args.Paging,
		Filters: args.Filters,
	}
	err := q.customerQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var customerIDs []dot.ID
	var listCustomerFbUser []*fbusering.ShopCustomerWithFbExternalUser
	for _, v := range query.Result.Customers {
		customerIDs = append(customerIDs, v.ID)
		listCustomerFbUser = append(listCustomerFbUser, &fbusering.ShopCustomerWithFbExternalUser{
			ShopCustomer: v,
		})
	}
	FbUserCustomers, err := q.fbExternalUserShopCustomerStore(ctx).ShopCustomerIDs(customerIDs).ShopID(args.ShopID).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	var fbUserIDs []string
	var mapFbUserExternalIDCustomerID = make(map[string]dot.ID)
	for _, v := range FbUserCustomers {
		fbUserIDs = append(fbUserIDs, v.FbExternalUserID)
		mapFbUserExternalIDCustomerID[v.FbExternalUserID] = v.CustomerID
	}
	fbUsers, err := q.fbUserStore(ctx).ExternalIDs(fbUserIDs).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	replaceDefaultAvatars(fbUsers, q.defaultAvatarLink)

	var mapFbUsers = make(map[string]*fbusering.FbExternalUser)
	for _, v := range fbUsers {
		mapFbUsers[v.ExternalID] = v
	}
	var customerfbUsers = make(map[dot.ID][]*fbusering.FbExternalUser)
	for _, v := range fbUsers {
		customerfbUsers[mapFbUserExternalIDCustomerID[v.ExternalID]] = append(customerfbUsers[mapFbUserExternalIDCustomerID[v.ExternalID]], mapFbUsers[v.ExternalID])
	}
	for k, v := range listCustomerFbUser {
		listCustomerFbUser[k].FbUsers = customerfbUsers[v.ID]
	}
	return &fbusering.ListShopCustomerWithFbExternalUserResponse{
		Customers: listCustomerFbUser,
		Paging:    query.Result.Paging,
	}, nil
}

func (q *FbUserQuery) ListFbExternalUserByIDs(ctx context.Context, extFbUserIDs []string) ([]*fbusering.FbExternalUser, error) {
	return q.fbUserStore(ctx).ExternalIDs(extFbUserIDs).ListFbExternalUsers()
}

func (q *FbUserQuery) GetShopUserTag(ctx context.Context, args *fbusering.GetShopUserTagArgs) (*fbusering.FbShopUserTag, error) {
	return q.fbShopUserTagStore(ctx).ByID(args.ID).ByShopID(args.ShopID).GetShopUserTag()
}

func (q *FbUserQuery) ListShopUserTags(ctx context.Context, args *fbusering.ListShopUserTagsArgs) ([]*fbusering.FbShopUserTag, error) {
	return q.fbShopUserTagStore(ctx).ByShopID(args.ShopID).ListShopUserTags()
}

func (q *FbUserQuery) ListFbExtUserShopCustomersByShopCustomerIDs(ctx context.Context, customerIDs []dot.ID) ([]*fbusering.FbExternalUserShopCustomer, error) {
	return q.fbExternalUserShopCustomerStore(ctx).ShopCustomerIDs(customerIDs).ListFbExternalUsers()
}

func (q *FbUserQuery) ListFbExternalUserIDsByShopCustomerIDs(ctx context.Context, customerIDs []dot.ID) ([]string, error) {
	extUsers, err := q.ListFbExtUserShopCustomersByShopCustomerIDs(ctx, customerIDs)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, user := range extUsers {
		result = append(result, user.FbExternalUserID)
	}
	return result, nil
}

func replaceDefaultAvatars(fbExternalUsers []*fbusering.FbExternalUser, linkAvatar string) {
	for _, fbExternalUser := range fbExternalUsers {
		replaceDefaultAvatar(fbExternalUser, linkAvatar)
	}
}

func replaceDefaultAvatar(fbExternalUser *fbusering.FbExternalUser, linkAvatar string) {
	if fbExternalUser != nil && fbExternalUser.ExternalInfo != nil &&
		fbExternalUser.ExternalInfo.ImageURL == fbclient.DefaultFaboImage {
		fbExternalUser.ExternalInfo.ImageURL = linkAvatar
	}
}
