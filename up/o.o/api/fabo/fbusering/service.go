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

	/* -- ShopTag -- */
	CreateShopUserTag(context.Context, *CreateShopUserTagArgs) (*FbShopUserTag, error)
	UpdateShopUserTag(context.Context, *UpdateShopUserTagArgs) (*FbShopUserTag, error)
	DeleteShopUserTag(context.Context, *DeleteShopUserTagArgs) (int, error)
	UpdateShopUserTags(context.Context, *UpdateShopUserTagsArgs) (*FbExternalUser, error)
}

type QueryService interface {
	GetShopUserTag(context.Context, *GetShopUserTagArgs) (*FbShopUserTag, error)
	ListShopUserTags(context.Context, *ListShopUserTagsArgs) ([]*FbShopUserTag, error)

	GetFbExternalUserByExternalID(_ context.Context, externalID string) (*FbExternalUser, error)
	ListFbExternalUsersByExternalIDs(_ context.Context, externalIDs filter.Strings, externalPageID dot.NullString) ([]*FbExternalUser, error)
	GetFbExternalUserInternalByExternalID(_ context.Context, externalID string) (*FbExternalUserInternal, error)

	// -- FbExternalUser with ShopCustomer --
	GetFbExternalUserWithCustomerByExternalID(_ context.Context, shopID dot.ID, externalID string) (*FbExternalUserWithCustomer, error)
	ListFbExternalUserWithCustomer(_ context.Context, args ListFbExternalUserWithCustomerRequest) ([]*FbExternalUserWithCustomer, error)
	ListFbExternalUserWithCustomerByExternalIDs(_ context.Context, shopID dot.ID, externalIDs []string) ([]*FbExternalUserWithCustomer, error)
	ListFbExternalUsers(context.Context, *ListFbExternalUsersArgs) ([]*FbExternalUserWithCustomer, error)
	ListShopCustomerWithFbExternalUser(context.Context, *ListCustomerWithFbAvatarsArgs) (*ListShopCustomerWithFbExternalUserResponse, error)
	ListFbExternalUserByIDs(ctx context.Context, extFbUserIDs []string) ([]*FbExternalUser, error)
	ListShopCustomerIDWithPhoneNorm(_ context.Context, shopID dot.ID, phone string) ([]dot.ID, error)
	ListFbExtUserShopCustomersByShopCustomerIDs(_ context.Context, customerIDs []dot.ID) ([]*FbExternalUserShopCustomer, error)
	ListFbExternalUserIDsByShopCustomerIDs(_ context.Context, customerIDs []dot.ID) ([]string, error)
}

type ListShopCustomerWithFbExternalUserResponse struct {
	Customers []*ShopCustomerWithFbExternalUser
	Paging    meta.PageInfo
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

// +convert:create=FbShopUserTag
type CreateShopUserTagArgs struct {
	Name   string
	Color  string
	ShopID dot.ID
}

// +convert:update=FbShopUserTag
type UpdateShopUserTagArgs struct {
	Name  string
	Color string
	ID    dot.ID
}

type DeleteShopUserTagArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetShopUserTagArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type ListShopUserTagsArgs struct {
	ShopID dot.ID
}

type FbExternalUserAddTagArgs struct {
	ShopID           dot.ID
	TagID            dot.ID
	FbExternalUserID dot.ID
}

type FbExternalUserRemoveTagArgs struct {
	ShopID           dot.ID
	TagID            dot.ID
	FbExternalUserID dot.ID
}

type UpdateShopUserTagsArgs struct {
	ShopID           dot.ID
	TagIDs           []dot.ID
	FbExternalUserID dot.ID
}
