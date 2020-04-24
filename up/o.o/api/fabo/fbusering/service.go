package fbusering

import (
	"context"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateFbUser(context.Context, *CreateFbUserArgs) (*FbUser, error)

	CreateFbUserInternal(context.Context, *CreateFbUserInternalArgs) (*FbUserInternal, error)

	CreateFbUserCombined(context.Context, *CreateFbUserCombinedArgs) (*FbUserCombined, error)
}

type QueryService interface {
	GetFbUserByID(_ context.Context, ID dot.ID) (*FbUser, error)
	GetFbUserByExternalID(_ context.Context, externalID string) (*FbUser, error)
	GetFbUserByUserID(_ context.Context, userID dot.ID) (*FbUser, error)

	GetFbUserInternalByID(_ context.Context, ID dot.ID) (*FbUserInternal, error)
}

// +convert:create=FbUser
type CreateFbUserArgs struct {
	ID           dot.ID
	ExternalID   string
	UserID       dot.ID
	ExternalInfo *ExternalFBUserInfo
	Token        string
	Status       status3.Status
}

// +convert:create=FbUserInternal
type CreateFbUserInternalArgs struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
}

type CreateFbUserCombinedArgs struct {
	UserID         dot.ID
	ShopID         dot.ID
	FbUser         *CreateFbUserArgs
	FbUserInternal *CreateFbUserInternalArgs
}
