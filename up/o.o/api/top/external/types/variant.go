package types

import (
	catalogtypes "o.o/api/main/catalog/types"
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
)

type ShopVariant struct {
	ExternalId   dot.NullString `json:"external_id"`
	ExternalCode dot.NullString `json:"external_code"`

	// @required
	Id dot.ID `json:"id"`

	Code dot.NullString `json:"code"`

	Name        dot.NullString     `json:"name"`
	Description dot.NullString     `json:"description"`
	ShortDesc   dot.NullString     `json:"short_desc"`
	ImageUrls   []string           `json:"image_urls"`
	ListPrice   dot.NullInt        `json:"list_price"`
	RetailPrice dot.NullInt        `json:"retail_price"`
	Note        dot.NullString     `json:"note"`
	Status      status3.NullStatus `json:"status"`

	CostPrice dot.NullInt `json:"cost_price"`

	Attributes []*catalogtypes.Attribute `json:"attributes"`

	Deleted bool `json:"deleted"`
}

func (m *ShopVariant) String() string { return jsonx.MustMarshalToString(m) }

type ShopVariantsResponse struct {
	ShopVariants []*ShopVariant         `json:"shop_variants"`
	Paging       *common.CursorPageInfo `json:"paging"`
}

func (m *ShopVariantsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetVariantRequest struct {
	Id         dot.ID `json:"id"`
	Code       string `json:"code"`
	ExternalId string `json:"external_id"`
}

func (m *GetVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type VariantFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *VariantFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListVariantsRequest struct {
	Filter         VariantFilter        `json:"filter"`
	Paging         *common.CursorPaging `json:"paging"`
	IncludeDeleted bool                 `json:"include_deleted"`
}

func (m *ListVariantsRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateVariantRequest struct {
	ExternalId   string `json:"external_id"`
	ExternalCode string `json:"external_code"`

	Code string `json:"code"`

	Name string `json:"name"`

	ProductId dot.ID `json:"product_id"`

	Note string `json:"note"`

	Description string `json:"description"`

	ShortDesc string `json:"short_desc"`

	CostPrice int `json:"cost_price"`

	ListPrice int `json:"list_price"`

	RetailPrice int `json:"retail_price"`

	ImageUrls []string `json:"image_urls"`

	Attributes []*catalogtypes.Attribute `json:"attributes"`
}

func (m *CreateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateVariantRequest struct {
	// @required
	Id dot.ID `json:"id"`

	Name dot.NullString `json:"name"`

	Note dot.NullString `json:"note"`

	Code dot.NullString `json:"code"`

	CostPrice dot.NullInt `json:"cost_price"`

	ListPrice dot.NullInt `json:"list_price"`

	RetailPrice dot.NullInt `json:"retail_price"`

	Description dot.NullString `json:"description"`

	ShortDesc dot.NullString `json:"short_desc"`

	// ImageURLs *[]string `json:"image_urls"`

	// Attributes *[]*catalogtypes.Attribute `json:"attributes"`
}

func (m *UpdateVariantRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *ShopVariant) HasChanged() bool {
	return m.Code.Valid ||
		m.Name.Valid ||
		m.Description.Valid ||
		m.ShortDesc.Valid ||
		m.ListPrice.Valid ||
		m.RetailPrice.Valid ||
		m.Note.Valid ||
		m.Status.Valid ||
		m.CostPrice.Valid
}
