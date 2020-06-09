package fbusering

import (
	"context"

	"o.o/api/meta"
	"o.o/api/shopping/customering"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

// +gen:api

type Aggregate interface {
	CreateFbExternalUser(context.Context, *CreateFbExternalUserArgs) (*FbExternalUser, error)
	CreateFbExternalUsers(context.Context, *CreateFbExternalUsersArgs) ([]*FbExternalUser, error)

	CreateFbExternalUserInternal(context.Context, *CreateFbExternalUserInternalArgs) (*FbExternalUserInternal, error)

	CreateFbExternalUserCombined(context.Context, *CreateFbExternalUserCombinedArgs) (*FbExternalUserCombined, error)

	CreateFbExternalUserShopCustomer(ctx context.Context, shopID dot.ID, externalID string, customerID dot.ID) (*FbExternalUserWithCustomer, error)

	DeleteFbExternalUserShopCustomer(context.Context, *DeleteFbExternalUserShopCustomerArgs) error
}

type QueryService interface {
	GetFbExternalUserByExternalID(_ context.Context, externalID string) (*FbExternalUser, error)
	ListFbExternalUsersByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalUser, error)

	GetFbExternalUserInternalByExternalID(_ context.Context, externalID string) (*FbExternalUserInternal, error)

	// -- FbExternalUser with ShopCustomer --
	GetFbExternalUserWithCustomerByExternalID(_ context.Context, shopID dot.ID, externalID string) (*FbExternalUserWithCustomer, error)
	ListFbExternalUserWithCustomer(_ context.Context, args ListFbExternalUserWithCustomerRequest) ([]*FbExternalUserWithCustomer, error)
	ListFbExternalUserWithCustomerByExternalIDs(_ context.Context, shopID dot.ID, externalIDs []string) ([]*FbExternalUserWithCustomer, error)
	ListFbExternalUsers(context.Context, *ListFbExternalUsersArgs) ([]*FbExternalUserWithCustomer, error)
	ListShopCustomerWithFbExternalUser(context.Context, *ListCustomerWithFbAvatarsArgs) ([]*ShopCustomerWithFbExternalUser, error)
}

type ListCustomerWithFbAvatarsArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type DeleteFbExternalUserShopCustomerArgs struct {
	ShopID     dot.ID
	ExternalID dot.NullString
	CustomerID dot.NullID
}

type ListFbExternalUserWithCustomerRequest struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

// +convert:create=FbExternalUser
type CreateFbExternalUserArgs struct {
	ExternalID     string
	ExternalInfo   *FbExternalUserInfo
	ExternalPageID string
	Status         status3.Status
}

type CreateFbExternalUsersArgs struct {
	FbExternalUsers []*CreateFbExternalUserArgs
}

// +convert:create=FbExternalUserInternal
type CreateFbExternalUserInternalArgs struct {
	ExternalID string
	Token      string
	ExpiresIn  int
}

type CreateFbExternalUserCombinedArgs struct {
	FbUser         *CreateFbExternalUserArgs
	FbUserInternal *CreateFbExternalUserInternalArgs
}

type FbExternalUserWithCustomer struct {
	*FbExternalUser
	*customering.ShopCustomer
}

type ListFbExternalUsersArgs struct {
	CustomerID dot.NullID
	ShopID     dot.ID
}
