package fbpaging

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateFbExternalPage(context.Context, *CreateFbExternalPageArgs) (*FbExternalPage, error)
	DisableFbExternalPagesByExternalIDs(context.Context, *DisableFbExternalPagesByIDsArgs) (int, error)
	DisableAllFbExternalPages(context.Context, *DisableAllFbExternalPagesArgs) (int, error)

	CreateFbExternalPageInternal(context.Context, *CreateFbExternalPageInternalArgs) (*FbExternalPageInternal, error)

	CreateFbExternalPageCombined(context.Context, *CreateFbExternalPageCombinedArgs) (*FbExternalPageCombined, error)
	CreateFbExternalPageCombineds(context.Context, *CreateFbExternalPageCombinedsArgs) ([]*FbExternalPageCombined, error)
}

type QueryService interface {
	GetFbExternalPageByID(_ context.Context, ID dot.ID) (*FbExternalPage, error)
	GetFbExternalPageByExternalID(_ context.Context, externalID string) (*FbExternalPage, error)
	GetFbExternalPageActiveByExternalID(_ context.Context, externalID string) (*FbExternalPage, error)

	GetFbExternalPageInternalByID(_ context.Context, ID dot.ID) (*FbExternalPageInternal, error)
	GetFbExternalPageInternalByExternalID(_ context.Context, externalID string) (*FbExternalPageInternal, error)

	ListFbExternalPages(context.Context, *ListFbExternalPagesArgs) (*FbPagesResponse, error)
	ListFbExternalPagesByIDs(_ context.Context, IDs []dot.ID) ([]*FbExternalPage, error)
	ListFbExternalPagesByExternalIDs(_ context.Context, externalIDs []string) ([]*FbExternalPage, error)
	ListFbExternalPagesActiveByExternalIDs(_ context.Context, externalIDs []string) ([]*FbExternalPage, error)
}

// +convert:create=FbExternalPage
type CreateFbExternalPageArgs struct {
	ID                   dot.ID
	ExternalID           string
	ShopID               dot.ID
	ExternalName         string
	ExternalCategory     string
	ExternalCategoryList []*ExternalCategory
	ExternalTasks        []string
	ExternalPermissions  []string
	ExternalImageURL     string
	Status               status3.Status
	ConnectionStatus     status3.Status
}

// +convert:create=FbExternalPageInternal
type CreateFbExternalPageInternalArgs struct {
	ID         dot.ID
	ExternalID string
	Token      string
}

type CreateFbExternalPageCombinedArgs struct {
	FbPage         *CreateFbExternalPageArgs
	FbPageInternal *CreateFbExternalPageInternalArgs
}

type CreateFbExternalPageCombinedsArgs struct {
	FbPageCombineds []*CreateFbExternalPageCombinedArgs
}

type DisableFbExternalPagesByIDsArgs struct {
	ExternalIDs []string
	ShopID      dot.ID
}

type DisableAllFbExternalPagesArgs struct {
	ShopID dot.ID
}

type ListFbExternalPagesArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type FbPagesResponse struct {
	FbPages []*FbExternalPage
	Paging  meta.PageInfo
}
