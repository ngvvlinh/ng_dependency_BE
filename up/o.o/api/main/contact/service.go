package contact

import (
	"context"

	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateContact(context.Context, *CreateContactArgs) (*Contact, error)
	UpdateContact(context.Context, *UpdateContactArgs) (*Contact, error)
	DeleteContact(context.Context, *DeleteContactArgs) (deleted int, _ error)
}

type QueryService interface {
	GetContactByID(context.Context, *GetContactByIDArgs) (*Contact, error)
	GetContacts(context.Context, *GetContactsArgs) (*GetContactsResponse, error)
	GetContactsByPhone(context.Context, *GetContactsByPhoneArgs) ([]*Contact, error)
}

//-- queries --//
type GetContactByIDArgs struct {
	ID     dot.ID
	ShopID dot.ID
}

type GetContactsArgs struct {
	ShopID dot.ID
	IDs    []dot.ID
	Phone  string

	Paging meta.Paging
}

type GetContactsByPhoneArgs struct {
	ShopID dot.ID
	Phone  string
}

type GetContactsResponse struct {
	Contacts []*Contact
	Paging   meta.PageInfo
}

//-- commands --//

// +convert:create=Contact
type CreateContactArgs struct {
	ShopID   dot.ID
	FullName string
	Phone    string
}

// +convert:update=Contact
type UpdateContactArgs struct {
	ID       dot.ID
	ShopID   dot.ID
	FullName dot.NullString
	Phone    dot.NullString
}

type DeleteContactArgs struct {
	ID     dot.ID
	ShopID dot.ID
}
