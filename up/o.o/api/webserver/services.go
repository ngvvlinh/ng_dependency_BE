package webserver

import (
	"context"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateOrUpdateWsCategory(context.Context, *CreateOrUpdateWsCategoryArgs) (*WsCategory, error)

	CreateOrUpdateWsProduct(context.Context, *CreateOrUpdateWsProductArgs) (*WsProduct, error)

	CreateWsPage(context.Context, *CreateWsPageArgs) (*WsPage, error)
	UpdateWsPage(context.Context, *UpdateWsPageArgs) (*WsPage, error)
	DeleteWsPage(ctx context.Context, shopID dot.ID, ID dot.ID) (int, error)

	CreateWsWebsite(context.Context, *CreateWsWebsiteArgs) (*WsWebsite, error)
	UpdateWsWebsite(context.Context, *UpdateWsWebsiteArgs) (*WsWebsite, error)
}

type QueryService interface {
	GetWsCategoryByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*WsCategory, error)
	ListWsCategoriesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*WsCategory, error)
	ListWsCategories(context.Context, ListWsCategoriesArgs) (*ListWsCategoriesResponse, error)

	GetWsProductByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*WsProduct, error)
	ListWsProductsByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*WsProduct, error)
	ListWsProducts(context.Context, ListWsProductsArgs) (*ListWsProductsResponse, error)

	GetWsPageByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*WsPage, error)
	ListWsPagesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*WsPage, error)
	ListWsPages(context.Context, ListWsPagesArgs) (*ListWsPagesResponse, error)

	GetWsWebsiteByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*WsWebsite, error)
	ListWsWebsitesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*WsWebsite, error)
	ListWsWebsites(context.Context, ListWsWebsitesArgs) (*ListWsWebsitesResponse, error)
}

// +convert:create=WsWebsite
type CreateWsWebsiteArgs struct {
	ShopID             dot.ID
	MainColor          string
	Banner             *Banner
	OutstandingProduct *SpecialProduct
	NewProduct         *SpecialProduct
	SEOConfig          *WsGeneralSEO
	Facebook           *Facebook
	GoogleAnalyticsID  string
	DomainName         string
	OverStock          bool
	ShopInfo           *ShopInfo
	Description        string
	LogoImage          string
	FaviconImage       string
}

// +convert:update=WsWebsite
type UpdateWsWebsiteArgs struct {
	ShopID             dot.ID
	ID                 dot.ID
	MainColor          dot.NullString
	Banner             *Banner
	OutstandingProduct *SpecialProduct
	NewProduct         *SpecialProduct
	SEOConfig          *WsGeneralSEO
	Facebook           *Facebook
	GoogleAnalyticsID  dot.NullString
	DomainName         dot.NullString
	OverStock          dot.NullBool
	ShopInfo           *ShopInfo
	Description        dot.NullString
	LogoImage          dot.NullString
	FaviconImage       dot.NullString
}

type ListWsWebsitesArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ListWsWebsitesResponse struct {
	PageInfo   meta.PageInfo
	WsWebsites []*WsWebsite
}

type ListWsPagesArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ListWsPagesResponse struct {
	ShopID   dot.ID
	PageInfo meta.PageInfo
	WsPages  []*WsPage
}

type ListWsProductsArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ListWsProductsResponse struct {
	ShopID     dot.ID
	PageInfo   meta.PageInfo
	WsProducts []*WsProduct
}

type ListWsCategoriesArgs struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters
}

type ListWsCategoriesResponse struct {
	PageInfo     meta.PageInfo
	WsCategories []*WsCategory
}

// +convert:update=WsCategory
type CreateOrUpdateWsCategoryArgs struct {
	ShopID    dot.ID
	ID        dot.ID
	Slug      dot.NullString
	SEOConfig *WsSEOConfig
	Image     dot.NullString
	Appear    dot.NullBool
}

// +convert:update=WsProduct
type CreateOrUpdateWsProductArgs struct {
	ID           dot.ID
	ShopID       dot.ID
	SEOConfig    *WsSEOConfig
	Slug         dot.NullString
	Appear       dot.NullBool
	ComparePrice []*ComparePrice
	DescHTML     dot.NullString
}

// +convert:create=WsPage
type CreateWsPageArgs struct {
	ShopID    dot.ID
	SEOConfig *WsSEOConfig
	Name      string
	Slug      string
	DescHTML  string
	Image     string
	Appear    bool
}

// +convert:update=WsPage
type UpdateWsPageArgs struct {
	ShopID    dot.ID
	ID        dot.ID
	SEOConfig *WsSEOConfig
	Name      dot.NullString
	Slug      dot.NullString
	DescHTML  dot.NullString
	Image     dot.NullString
	Appear    dot.NullBool
}
