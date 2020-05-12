package fbuser

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

var _ fbusering.QueryService = &FbUserQuery{}

type FbUserQuery struct {
	db                              *cmsql.Database
	fbUserStore                     sqlstore.FbExternalUserStoreFactory
	fbUserInternalStore             sqlstore.FbExternalUserInternalFactory
	fbExternalUserShopCustomerStore sqlstore.FbExternalUserShopCustomerStoreFactory
	customerQuery                   customering.QueryBus
}

func NewFbUserQuery(database *cmsql.Database, customerQ customering.QueryBus) *FbUserQuery {
	return &FbUserQuery{
		db:                              database,
		fbUserStore:                     sqlstore.NewFbExternalUserStore(database),
		fbExternalUserShopCustomerStore: sqlstore.NewFbExternalUserShopCustomerStore(database),
		customerQuery:                   customerQ,
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
	return fbUser, nil
}

func (q *FbUserQuery) ListFbExternalUsersByExternalIDs(
	ctx context.Context, externalIDs filter.Strings,
) ([]*fbusering.FbExternalUser, error) {
	return q.fbUserStore(ctx).ExternalIDs(externalIDs).ListFbExternalUsers()
}

func (q *FbUserQuery) ListFbExternalUsers(ctx context.Context, args *fbusering.ListFbExternalUsersArgs) ([]*fbusering.FbExternalUserWithCustomer, error) {
	query := q.fbExternalUserShopCustomerStore(ctx).ShopID(args.ShopID)
	if args.CustomerID.Valid {
		query = query.ShopCustomerID(args.CustomerID.ID)
	}
	fbUserCustomer, err := query.ListFbExternalUser()
	if err != nil {
		return nil, err
	}
	var fbUserIDs []string
	for _, v := range fbUserCustomer {
		fbUserIDs = append(fbUserIDs, v.FbExternalUserID)
	}
	result, err := q.fbUserStore(ctx).ExternalIDs(fbUserIDs).ListFbExternalUsers()
	if err != nil {
		return nil, err
	}
	return q.populateFbExternalUsersWithCustomerInfo(ctx, args.ShopID, result)
}

func (q *FbUserQuery) GetFbExternalUserWithCustomerByExternalID(ctx context.Context, shopID dot.ID, externalID string) (*fbusering.FbExternalUserWithCustomer, error) {
	fbExternalUser, err := q.fbUserStore(ctx).ExternalID(externalID).GetFbExternalUser()
	if err != nil {
		return nil, err
	}
	var result = &fbusering.FbExternalUserWithCustomer{
		FbExternalUser: fbExternalUser,
	}

	fbExternalUserShopCustomer, err := q.fbExternalUserShopCustomerStore(ctx).ShopID(shopID).FbExternalUserID(externalID).ListFbExternalUser()
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
	fbUsersCustomers, err := q.fbExternalUserShopCustomerStore(ctx).WithPaging(args.Paging).Filters(args.Filters).ListFbExternalUser()
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
	return q.populateFbExternalUsersWithCustomerInfo(ctx, args.ShopID, fbUsers)
}

func (q *FbUserQuery) ListFbExternalUserWithCustomerByExternalIDs(ctx context.Context, shopID dot.ID, externalID []string) ([]*fbusering.FbExternalUserWithCustomer, error) {
	fbUsersCustomers, err := q.fbExternalUserShopCustomerStore(ctx).FbExternalUserIDs(externalID).ListFbExternalUser()
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

	fbUserWithCustomer, err := q.fbExternalUserShopCustomerStore(ctx).FbExternalUserIDs(fbUserIDS).ListFbExternalUser()
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
