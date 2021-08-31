package types

import (
	"o.o/api/top/types/common"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type ListContactsRequest struct {
	Paging         *common.CursorPaging `json:"paging"`
	Filter         *FilterGetContacts   `json:"filter"`
	IncludeDeleted bool                 `json:"include_deleted"`
}

func (m *ListContactsRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetContacts struct {
	IDs   []dot.ID `json:"ids"`
	Phone string   `json:"phone"`
	Name  string   `json:"name"`
}

func (m *FilterGetContacts) String() string { return jsonx.MustMarshalToString(m) }

type ContactsResponse struct {
	Contacts []*Contact             `json:"contacts"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *ContactsResponse) String() string { return jsonx.MustMarshalToString(m) }

type Contact struct {
	ID        dot.ID   `json:"id"`
	ShopID    dot.ID   `json:"shop_id"`
	FullName  string   `json:"full_name"`
	Phone     string   `json:"phone"`
	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *Contact) String() string { return jsonx.MustMarshalToString(m) }

type CreateContactRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

func (m *CreateContactRequest) String() string { return jsonx.MustMarshalToString(m) }
