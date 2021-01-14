package fbpaging

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	DisableFbExternalPagesByExternalIDs(context.Context, *DisableFbExternalPagesByIDsArgs) (int, error)

	CreateFbExternalPageCombineds(context.Context, *CreateFbExternalPageCombinedsArgs) ([]*FbExternalPageCombined, error)
}

type QueryService interface {
	GetFbExternalPageByExternalID(_ context.Context, externalID string) (*FbExternalPage, error)
	ListFbExternalPages(context.Context, *ListFbExternalPagesArgs) (*FbPagesResponse, error)
	GetFbExternalPageActiveByExternalID(_ context.Context, externalID string) (*FbExternalPage, error)

	GetFbExternalPageInternalByExternalID(_ context.Context, externalID string) (*FbExternalPageInternal, error)
	GetPageAccessToken(_ context.Context, externalID string) (string, error)
	GetFbExternalPageInternalActiveByExternalID(_ context.Context, externalID string) (*FbExternalPageInternal, error)

	ListFbExternalPagesByExternalIDs(_ context.Context, externalIDs []string) ([]*FbExternalPage, error)
	ListFbExternalPagesActiveByExternalIDs(_ context.Context, externalIDs []string) ([]*FbExternalPage, error)
	ListFbPagesByShop(_ context.Context, shopIDs []dot.ID) ([]*FbExternalPage, error)
	ListActiveFbPagesByShopIDs(_ context.Context, shopIDs []dot.ID) ([]*FbExternalPage, error)
}

// +convert:create=FbExternalPage
type CreateFbExternalPageArgs struct {
	ID                   dot.ID
	ExternalID           string
	ExternalUserID       string
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
