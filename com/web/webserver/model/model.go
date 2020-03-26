package model

import (
	"time"

	"o.o/api/top/types/etc/ws_banner_show_style"
	"o.o/api/top/types/etc/ws_list_product_show_style"
	"o.o/capi/dot"
)

// +sqlgen
type WsCategory struct {
	ID        dot.ID
	ShopID    dot.ID
	Slug      string
	SEOConfig *WsSEOConfig
	Image     string
	Appear    bool
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

type WsSEOConfig struct {
	Content     string `json:"content"`
	Keyword     string `json:"keyword"`
	Description string `json:"description"`
}

// +sqlgen
type WsProduct struct {
	ID           dot.ID
	ShopID       dot.ID
	SEOConfig    *WsSEOConfig
	Slug         string
	Appear       bool
	ComparePrice []*ComparePrice
	DescHTML     string
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

// +sqlgen
type WsPage struct {
	ShopID    dot.ID
	ID        dot.ID
	SEOConfig *WsSEOConfig
	Name      string
	Slug      string
	DescHTML  string
	Image     string
	Appear    bool
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

// +sqlgen
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
	OverStock          int
	ShopInfo           *ShopInfo
	Description        string
	LogoImage          string
	FaviconImage       string
	SiteSubdomain      string
	CreatedAt          time.Time `sq:"create"`
	UpdatedAt          time.Time `sq:"update"`
}

type ComparePrice struct {
	VariantID    dot.ID `json:"variant_id"`
	ComparePrice int    `json:"compare_price"`
}

type ShopInfo struct {
	Email           string           `json:"email"`
	Name            string           `json:"name"`
	Phone           string           `json:"phone"`
	Address         *AddressShopInfo `json:"address"`
	FacebookFanpage string           `json:"facebook_fanpage"`
}

type Facebook struct {
	FacebookID     string `json:"facebook_id"`
	WelcomeMessage string `json:"welcome_message"`
}

type WsGeneralSEO struct {
	Title               string `json:"title"`
	SiteContent         string `json:"site_content"`
	SiteMetaKeyword     string `json:"site_meta_keyword"`
	SiteMetaDescription string `json:"site_meta_description"`
}

type Banner struct {
	BannerItems []*BannerItem                          `json:"banner_items"`
	Style       ws_banner_show_style.WsBannerShowStyle `json:"style"`
}

type BannerItem struct {
	Alt   string `json:"alt"`
	Url   string `json:"url"`
	Image string `json:"image"`
}

type SpecialProduct struct {
	ProductIDs []dot.ID                                          `json:"product_ids"`
	Style      ws_list_product_show_style.WsListProductShowStyle `json:"style"`
}

type AddressShopInfo struct {
	ProvinceCode string `json:"province_code"`
	DistrictCode string `json:"district_code"`
	WardCode     string `json:"ward_code"`

	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`

	Address string `json:"address"`
}
