package fbusering

import (
	"context"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/filter"
)

// +gen:api

type Aggregate interface {
	CreateFbExternalUser(context.Context, *CreateFbExternalUserArgs) (*FbExternalUser, error)
	CreateFbExternalUsers(context.Context, *CreateFbExternalUsersArgs) ([]*FbExternalUser, error)

	CreateFbExternalUserInternal(context.Context, *CreateFbExternalUserInternalArgs) (*FbExternalUserInternal, error)

	CreateFbExternalUserCombined(context.Context, *CreateFbExternalUserCombinedArgs) (*FbExternalUserCombined, error)
}

type QueryService interface {
	GetFbExternalUserByExternalID(_ context.Context, externalID string) (*FbExternalUser, error)
	ListFbExternalUsersByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalUser, error)

	GetFbExternalUserInternalByExternalID(_ context.Context, externalID string) (*FbExternalUserInternal, error)
}

// +convert:create=FbExternalUser
type CreateFbExternalUserArgs struct {
	ExternalID   string
	ExternalInfo *FbExternalUserInfo
	Token        string
	Status       status3.Status
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
