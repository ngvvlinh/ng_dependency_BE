package fbpaging

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
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
	GetFbPageByID(context.Context, *GetFbPageByIDArgs) (*FbPage, error)
	GetFbPageByExternalID(context.Context, *GetFbPageByExternalIDArgs) (*FbPage, error)
	ListFbPagesByIDs(context.Context, *ListFbPagesByIDsArgs) ([]*FbPage, error)
	ListFbPages(context.Context, *ListFbPagesArgs) (*FbPagesResponse, error)

	GetFbPageInternalByID(context.Context, *GetFbPageInternalByIDArgs) (*FbPageInternal, error)
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
	Status               status3.Status
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

type GetFbPageByIDArgs struct {
	ID dot.ID
}

type GetFbPageByExternalIDArgs struct {
	ExternalID string
}

type GetFbPageInternalByIDArgs struct {
	ID dot.ID
}

type ListFbPagesByIDsArgs struct {
	IDs []dot.ID
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
