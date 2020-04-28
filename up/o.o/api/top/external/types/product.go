package types

import (
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
)

type Tag struct {
	Id    dot.ID `json:"id"`
	Label string `json:"label"`
}

func (m *Tag) String() string { return jsonx.MustMarshalToString(m) }

type ShopProduct struct {
	ExternalId   dot.NullString `json:"external_id"`
	ExternalCode dot.NullString `json:"external_code"`

	// @required
	Id dot.ID `json:"id"`

	Name        dot.NullString     `json:"name"`
	Description dot.NullString     `json:"description"`
	ShortDesc   dot.NullString     `json:"short_desc"`
	ImageUrls   []string           `json:"image_urls"`
	CategoryId  dot.NullID         `json:"category_id"`
	Note        dot.NullString     `json:"note"`
	Status      status3.NullStatus `json:"status"`
	ListPrice   dot.NullInt        `json:"list_price"`
	RetailPrice dot.NullInt        `json:"retail_price"`
	Variants    []*ShopVariant     `json:"variants"`

	CreatedAt dot.Time   `json:"created_at"`
	UpdatedAt dot.Time   `json:"updated_at"`
	BrandId   dot.NullID `json:"brand_id"`

	Deleted bool `json:"deleted"`
}

func (m *ShopProduct) String() string { return jsonx.MustMarshalToString(m) }

type ShopProductsResponse struct {
	Products []*ShopProduct         `json:"products"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *ShopProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetProductRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *GetProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *ProductFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListProductsRequest struct {
	Filter         ProductFilter        `json:"filter"`
	Paging         *common.CursorPaging `json:"paging"`
	IncludeDeleted bool                 `json:"include_deleted"`
}

func (m *ListProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateProductRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Unit        string   `json:"unit"`
	Note        string   `json:"note"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	ImageUrls   []string `json:"image_urls"`
	CostPrice   int      `json:"cost_price"`
	ListPrice   int      `json:"list_price"`
	RetailPrice int      `json:"retail_price"`
	BrandId     dot.ID   `json:"brand_id"`
}

func (m *CreateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateProductRequest struct {
	// @required
	Id dot.ID `json:"id"`

	Name dot.NullString `json:"name"`

	Note dot.NullString `json:"note"`

	Unit dot.NullString `json:"unit"`

	Description dot.NullString `json:"description"`

	ShortDesc dot.NullString `json:"short_desc"`

	// ImageURLs *[]string `json:"image_urls"`

	CostPrice dot.NullInt `json:"cost_price"`

	ListPrice dot.NullInt `json:"list_price"`

	RetailPrice dot.NullInt `json:"retail_price"`

	BrandId dot.NullID `json:"brand_id"`
}

func (m *UpdateProductRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *ShopProduct) HasChanged() bool {
	return m.Name.Valid ||
		m.Description.Valid ||
		m.ShortDesc.Valid ||
		m.Note.Valid ||
		m.Status.Valid ||
		m.ListPrice.Valid ||
		m.RetailPrice.Valid ||
		m.BrandId.Valid ||
		m.CategoryId.Valid
}
