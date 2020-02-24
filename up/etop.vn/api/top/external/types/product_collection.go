package types

import (
	"etop.vn/api/top/types/common"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
	"etop.vn/common/jsonx"
)

type ProductCollection struct {
	ID          dot.ID   `json:"id"`
	ShopID      dot.ID   `json:"shop_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ShortDesc   string   `json:"short_desc"`
	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
	Deleted     bool     `json:"deleted"`
}

func (m *ProductCollection) String() string { return jsonx.MustMarshalToString(m) }

type ListProductCollectionRelationshipsFilter struct {
	ProductID    filter.IDs `json:"product_id"`
	CollectionID filter.IDs `json:"collection_id"`
}

func (m *ListProductCollectionRelationshipsFilter) String() string {
	return jsonx.MustMarshalToString(m)
}

type ListProductCollectionRelationshipsRequest struct {
	Filter ListProductCollectionRelationshipsFilter `json:"filter"`
	// TODO: add cursor paging
	//Paging         *common.CursorPaging                     `json:"paging"`
	IncludeDeleted bool `json:"include_deleted"`
}

func (m *ListProductCollectionRelationshipsRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type ProductCollectionRelationshipsResponse struct {
	Relationships []*ProductCollectionRelationship `json:"relationships"`
	Paging        *common.CursorPageInfo           `json:"paging"`
}

func (m *ProductCollectionRelationshipsResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type CreateCollectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ShortDesc   string `json:"short_desc"`
}

func (m *CreateCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateCollectionRequest struct {
	ID          dot.ID         `json:"id"`
	Name        dot.NullString `json:"name"`
	Description dot.NullString `json:"description"`
	ShortDesc   dot.NullString `json:"short_desc"`
}

func (m *UpdateCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductCollectionRelationship struct {
	ProductId    dot.ID `json:"product_id"`
	CollectionId dot.ID `json:"collection_id"`
	Deleted      bool   `json:"deleted"`
}

func (m *ProductCollectionRelationship) String() string { return jsonx.MustMarshalToString(m) }

type CreateProductCollectionRelationshipRequest struct {
	ProductId    dot.ID `json:"product_id"`
	CollectionId dot.ID `json:"collection_id"`
}

func (m *CreateProductCollectionRelationshipRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type RemoveProductCollectionRequest struct {
	ProductId    dot.ID `json:"product_id"`
	CollectionId dot.ID `json:"collection_id"`
}

func (m *RemoveProductCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ProductCollectionsResponse struct {
	Collections []*ProductCollection   `json:"collections"`
	Paging      *common.CursorPageInfo `json:"paging"`
}

func (m *ProductCollectionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetCollectionRequest struct {
	ID dot.ID `json:"id"`
}

func (m *GetCollectionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ListCollectionsFilter struct {
	ID filter.IDs `json:"id"`
}

func (m *ListCollectionsFilter) String() string { return jsonx.MustMarshalToString(m) }

type ListCollectionsRequest struct {
	Filter         ListCollectionsFilter `json:"filter"`
	Paging         *common.CursorPaging  `json:"paging"`
	IncludeDeleted bool                  `json:"include_deleted"`
}

func (m *ListCollectionsRequest) String() string { return jsonx.MustMarshalToString(m) }
