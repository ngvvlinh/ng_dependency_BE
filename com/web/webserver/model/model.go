package model

import (
	"time"

	"etop.vn/api/top/types/etc/ws_banner_show_style"
	"etop.vn/api/top/types/etc/ws_list_product_show_style"
	"etop.vn/backend/com/main/address/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenWsCategory(&WsCategory{})

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

var _ = sqlgenWsProduct(&WsProduct{})

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

var _ = sqlgenWsPage(&WsPage{})

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

var _ = sqlgenWsWebsite(&WsWebsite{})

type WsWebsite struct {
	ShopID             dot.ID
	ID                 dot.ID
	MainColor          string
	Banner             []*Banner
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
	CreatedAt          time.Time `sq:"create"`
	UpdatedAt          time.Time `sq:"update"`
}

type ComparePrice struct {
	VariantID    dot.ID `json:"variant_id"`
	ComparePrice int    `json:"compare_price"`
}

type ShopInfo struct {
	Name            string         `json:"name"`
	Phone           string         `json:"phone"`
	Address         *model.Address `json:"address"`
	FacebookFanpage string         `json:"facebook_fanpage"`
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
