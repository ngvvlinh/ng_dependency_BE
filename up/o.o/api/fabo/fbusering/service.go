package fbusering

import (
	"context"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateFbExternalUser(context.Context, *CreateFbExternalUserArgs) (*FbExternalUser, error)

	CreateFbExternalUserInternal(context.Context, *CreateFbExternalUserInternalArgs) (*FbExternalUserInternal, error)

	CreateFbExternalUserCombined(context.Context, *CreateFbExternalUserCombinedArgs) (*FbExternalUserCombined, error)
}

type QueryService interface {
	GetFbExternalUserByID(_ context.Context, ID dot.ID) (*FbExternalUser, error)
	GetFbExternalUserByExternalID(_ context.Context, externalID string) (*FbExternalUser, error)
	GetFbExternalUserByUserID(_ context.Context, userID dot.ID) (*FbExternalUser, error)

	GetFbExternalUserInternalByID(_ context.Context, ID dot.ID) (*FbExternalUserInternal, error)
}

// +convert:create=FbExternalUser
type CreateFbExternalUserArgs struct {
	ID           dot.ID
	UserID       dot.ID
	ExternalID   string
	ExternalInfo *FbExternalUserInfo
	Token        string
	Status       status3.Status
}

// +convert:create=FbExternalUserInternal
type CreateFbExternalUserInternalArgs struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
}

type CreateFbExternalUserCombinedArgs struct {
	UserID         dot.ID
	ShopID         dot.ID
	FbUser         *CreateFbExternalUserArgs
	FbUserInternal *CreateFbExternalUserInternalArgs
}
