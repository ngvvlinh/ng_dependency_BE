package fbpaging

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

// +gen:api

type Aggregate interface {
	CreateFbPage(context.Context, *CreateFbPageArgs) (*FbPage, error)
	DisableFbPagesByIDs(context.Context, *DisableFbPagesByIDsArgs) (int, error)
	DisableAllFbPages(context.Context, *DisableAllFbPagesArgs) (int, error)

	CreateFbPageInternal(context.Context, *CreateFbPageInternalArgs) (*FbPageInternal, error)

	CreateFbPageCombined(context.Context, *CreateFbPageCombinedArgs) (*FbPageCombined, error)
	CreateFbPageCombineds(context.Context, *CreateFbPageCombinedsArgs) ([]*FbPageCombined, error)
}

type QueryService interface {
	GetFbPageByID(_ context.Context, ID dot.ID) (*FbPage, error)
	GetFbPageByExternalID(_ context.Context, externalID string) (*FbPage, error)
	ListFbPagesByIDs(_ context.Context, IDs filter.IDs) ([]*FbPage, error)
	ListFbPages(context.Context, *ListFbPagesArgs) (*FbPagesResponse, error)
	ListFbPagesActiveByExternalIDs(_ context.Context, externalIDs []string) ([]*FbPage, error)

	GetFbPageInternalByID(_ context.Context, ID dot.ID) (*FbPageInternal, error)
}

// +convert:create=FbPage
type CreateFbPageArgs struct {
	ID                   dot.ID
	ExternalID           string
	FbUserID             dot.ID
	ShopID               dot.ID
	UserID               dot.ID
	ExternalName         string
	ExternalCategory     string
	ExternalCategoryList []*ExternalCategory
	ExternalTasks        []string
	ExternalPermissions  []string
	ExternalImageURL     string
	Status               status3.Status
	ConnectionStatus     status3.Status
}

// +convert:create=FbPageInternal
type CreateFbPageInternalArgs struct {
	ID    dot.ID
	Token string
}

type CreateFbPageCombinedArgs struct {
	FbPage         *CreateFbPageArgs
	FbPageInternal *CreateFbPageInternalArgs
}

type CreateFbPageCombinedsArgs struct {
	ShopID          dot.ID
	UserID          dot.ID
	FbPageCombineds []*CreateFbPageCombinedArgs
}

type DisableFbPagesByIDsArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
	UserID dot.ID
}

type DisableAllFbPagesArgs struct {
	ShopID dot.ID
	UserID dot.ID
}

type ListFbPagesArgs struct {
	ShopID   dot.ID
	UserID   dot.ID
	FbUserID dot.NullID

	Paging  meta.Paging
	Filters meta.Filters
}

type FbPagesResponse struct {
	FbPages []*FbPage
	Paging  meta.PageInfo
}
