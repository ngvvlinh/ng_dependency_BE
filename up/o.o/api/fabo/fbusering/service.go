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
	GetFbUserByID(context.Context, *GetFbUserByIDArgs) (*FbUser, error)
	GetFbUserByExternalID(context.Context, *GetFbUserByExternalIDArgs) (*FbUser, error)
	GetFbUserByUserID(context.Context, *GetFbUserByUserIDArgs) (*FbUser, error)

	GetFbUserInternalByID(context.Context, *GetFbUserInternalByIDArgs) (*FbUserInternal, error)
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

type GetFbUserByIDArgs struct {
	ID dot.ID
}

type GetFbUserByExternalIDArgs struct {
	ExternalID string
}

type GetFbUserByUserIDArgs struct {
	UserID dot.ID
}

type GetFbUserInternalByIDArgs struct {
	ID dot.ID
}
