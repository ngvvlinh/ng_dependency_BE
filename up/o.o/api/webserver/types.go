package webserver

import (
	"time"

	"o.o/api/main/catalog"
	"o.o/api/top/types/etc/ws_banner_show_style"
	"o.o/api/top/types/etc/ws_list_product_show_style"
	"o.o/capi/dot"
)

const MaxSlide = 5

type WsCategory struct {
	ID        dot.ID
	ShopID    dot.ID
	Slug      string
	SEOConfig *WsSEOConfig
	Image     string
	Appear    bool
	Category  *catalog.ShopCategory
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WsSEOConfig struct {
	Content     string
	Keyword     string
	Description string
}

type WsProduct struct {
	ID           dot.ID
	ShopID       dot.ID
	SEOConfig    *WsSEOConfig
	Slug         string
	Appear       bool
	ComparePrice []*ComparePrice
	DescHTML     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Product      *catalog.ShopProductWithVariants
	IsSale       bool
}
type ComparePrice struct {
	VariantID    dot.ID
	ComparePrice int
}

type WsPage struct {
	ShopID    dot.ID
	ID        dot.ID
	SEOConfig *WsSEOConfig
	Name      string
	Slug      string
	DescHTML  string
	Image     string
	Appear    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type WsWebsite struct {
	ShopID             dot.ID
	ID                 dot.ID
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
	CreatedAt          time.Time
	UpdatedAt          time.Time
	SiteSubdomain      string
}

type ShopInfo struct {
	Name            string
	Email           string
	Phone           string
	Address         *AddressShopInfo
	FacebookFanpage string
}

type AddressShopInfo struct {
	ProvinceCode string
	DistrictCode string
	WardCode     string

	Province string
	District string
	Ward     string

	Address string
}

type Facebook struct {
	FacebookID     string
	WelcomeMessage string
}

type WsGeneralSEO struct {
	Title               string
	SiteContent         string
	SiteMetaKeyword     string
	SiteMetaDescription string
}

type Banner struct {
	BannerItems []*BannerItem
	Style       ws_banner_show_style.WsBannerShowStyle
}

type BannerItem struct {
	Alt   string
	Url   string
	Image string
}

type SpecialProduct struct {
	ProductIDs []dot.ID
	Style      ws_list_product_show_style.WsListProductShowStyle
	Products   []*WsProduct
}
