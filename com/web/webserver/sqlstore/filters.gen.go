// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type WsCategoryFilters struct{ prefix string }

func NewWsCategoryFilters(prefix string) WsCategoryFilters {
	return WsCategoryFilters{prefix}
}

func (ft *WsCategoryFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft WsCategoryFilters) Prefix() string {
	return ft.prefix
}

func (ft *WsCategoryFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *WsCategoryFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *WsCategoryFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *WsCategoryFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *WsCategoryFilters) BySlug(Slug string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "slug",
		Value:  Slug,
		IsNil:  Slug == "",
	}
}

func (ft *WsCategoryFilters) BySlugPtr(Slug *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "slug",
		Value:  Slug,
		IsNil:  Slug == nil,
		IsZero: Slug != nil && (*Slug) == "",
	}
}

func (ft *WsCategoryFilters) ByImage(Image string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "image",
		Value:  Image,
		IsNil:  Image == "",
	}
}

func (ft *WsCategoryFilters) ByImagePtr(Image *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "image",
		Value:  Image,
		IsNil:  Image == nil,
		IsZero: Image != nil && (*Image) == "",
	}
}

func (ft *WsCategoryFilters) ByAppear(Appear bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "appear",
		Value:  Appear,
		IsNil:  bool(!Appear),
	}
}

func (ft *WsCategoryFilters) ByAppearPtr(Appear *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "appear",
		Value:  Appear,
		IsNil:  Appear == nil,
		IsZero: Appear != nil && bool(!(*Appear)),
	}
}

func (ft *WsCategoryFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *WsCategoryFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *WsCategoryFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *WsCategoryFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type WsPageFilters struct{ prefix string }

func NewWsPageFilters(prefix string) WsPageFilters {
	return WsPageFilters{prefix}
}

func (ft *WsPageFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft WsPageFilters) Prefix() string {
	return ft.prefix
}

func (ft *WsPageFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *WsPageFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *WsPageFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *WsPageFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *WsPageFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *WsPageFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *WsPageFilters) BySlug(Slug string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "slug",
		Value:  Slug,
		IsNil:  Slug == "",
	}
}

func (ft *WsPageFilters) BySlugPtr(Slug *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "slug",
		Value:  Slug,
		IsNil:  Slug == nil,
		IsZero: Slug != nil && (*Slug) == "",
	}
}

func (ft *WsPageFilters) ByDescHTML(DescHTML string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == "",
	}
}

func (ft *WsPageFilters) ByDescHTMLPtr(DescHTML *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == nil,
		IsZero: DescHTML != nil && (*DescHTML) == "",
	}
}

func (ft *WsPageFilters) ByImage(Image string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "image",
		Value:  Image,
		IsNil:  Image == "",
	}
}

func (ft *WsPageFilters) ByImagePtr(Image *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "image",
		Value:  Image,
		IsNil:  Image == nil,
		IsZero: Image != nil && (*Image) == "",
	}
}

func (ft *WsPageFilters) ByAppear(Appear bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "appear",
		Value:  Appear,
		IsNil:  bool(!Appear),
	}
}

func (ft *WsPageFilters) ByAppearPtr(Appear *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "appear",
		Value:  Appear,
		IsNil:  Appear == nil,
		IsZero: Appear != nil && bool(!(*Appear)),
	}
}

func (ft *WsPageFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *WsPageFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *WsPageFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *WsPageFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *WsPageFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *WsPageFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

type WsProductFilters struct{ prefix string }

func NewWsProductFilters(prefix string) WsProductFilters {
	return WsProductFilters{prefix}
}

func (ft *WsProductFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft WsProductFilters) Prefix() string {
	return ft.prefix
}

func (ft *WsProductFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *WsProductFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *WsProductFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *WsProductFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *WsProductFilters) BySlug(Slug string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "slug",
		Value:  Slug,
		IsNil:  Slug == "",
	}
}

func (ft *WsProductFilters) BySlugPtr(Slug *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "slug",
		Value:  Slug,
		IsNil:  Slug == nil,
		IsZero: Slug != nil && (*Slug) == "",
	}
}

func (ft *WsProductFilters) ByAppear(Appear bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "appear",
		Value:  Appear,
		IsNil:  bool(!Appear),
	}
}

func (ft *WsProductFilters) ByAppearPtr(Appear *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "appear",
		Value:  Appear,
		IsNil:  Appear == nil,
		IsZero: Appear != nil && bool(!(*Appear)),
	}
}

func (ft *WsProductFilters) ByDescHTML(DescHTML string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == "",
	}
}

func (ft *WsProductFilters) ByDescHTMLPtr(DescHTML *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "desc_html",
		Value:  DescHTML,
		IsNil:  DescHTML == nil,
		IsZero: DescHTML != nil && (*DescHTML) == "",
	}
}

func (ft *WsProductFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *WsProductFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *WsProductFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *WsProductFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type WsWebsiteFilters struct{ prefix string }

func NewWsWebsiteFilters(prefix string) WsWebsiteFilters {
	return WsWebsiteFilters{prefix}
}

func (ft *WsWebsiteFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft WsWebsiteFilters) Prefix() string {
	return ft.prefix
}

func (ft *WsWebsiteFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *WsWebsiteFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *WsWebsiteFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *WsWebsiteFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *WsWebsiteFilters) ByMainColor(MainColor string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "main_color",
		Value:  MainColor,
		IsNil:  MainColor == "",
	}
}

func (ft *WsWebsiteFilters) ByMainColorPtr(MainColor *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "main_color",
		Value:  MainColor,
		IsNil:  MainColor == nil,
		IsZero: MainColor != nil && (*MainColor) == "",
	}
}

func (ft *WsWebsiteFilters) ByGoogleAnalyticsID(GoogleAnalyticsID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "google_analytics_id",
		Value:  GoogleAnalyticsID,
		IsNil:  GoogleAnalyticsID == "",
	}
}

func (ft *WsWebsiteFilters) ByGoogleAnalyticsIDPtr(GoogleAnalyticsID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "google_analytics_id",
		Value:  GoogleAnalyticsID,
		IsNil:  GoogleAnalyticsID == nil,
		IsZero: GoogleAnalyticsID != nil && (*GoogleAnalyticsID) == "",
	}
}

func (ft *WsWebsiteFilters) ByDomainName(DomainName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "domain_name",
		Value:  DomainName,
		IsNil:  DomainName == "",
	}
}

func (ft *WsWebsiteFilters) ByDomainNamePtr(DomainName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "domain_name",
		Value:  DomainName,
		IsNil:  DomainName == nil,
		IsZero: DomainName != nil && (*DomainName) == "",
	}
}

func (ft *WsWebsiteFilters) ByOverStock(OverStock int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "over_stock",
		Value:  OverStock,
		IsNil:  OverStock == 0,
	}
}

func (ft *WsWebsiteFilters) ByOverStockPtr(OverStock *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "over_stock",
		Value:  OverStock,
		IsNil:  OverStock == nil,
		IsZero: OverStock != nil && (*OverStock) == 0,
	}
}

func (ft *WsWebsiteFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *WsWebsiteFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *WsWebsiteFilters) ByLogoImage(LogoImage string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "logo_image",
		Value:  LogoImage,
		IsNil:  LogoImage == "",
	}
}

func (ft *WsWebsiteFilters) ByLogoImagePtr(LogoImage *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "logo_image",
		Value:  LogoImage,
		IsNil:  LogoImage == nil,
		IsZero: LogoImage != nil && (*LogoImage) == "",
	}
}

func (ft *WsWebsiteFilters) ByFaviconImage(FaviconImage string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "favicon_image",
		Value:  FaviconImage,
		IsNil:  FaviconImage == "",
	}
}

func (ft *WsWebsiteFilters) ByFaviconImagePtr(FaviconImage *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "favicon_image",
		Value:  FaviconImage,
		IsNil:  FaviconImage == nil,
		IsZero: FaviconImage != nil && (*FaviconImage) == "",
	}
}

func (ft *WsWebsiteFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *WsWebsiteFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *WsWebsiteFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *WsWebsiteFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}
