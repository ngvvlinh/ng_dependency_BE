package customering

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCustomer(ctx context.Context, _ *CreateCustomerArgs) (*ShopCustomer, error)

	UpdateCustomer(ctx context.Context, _ *UpdateCustomerArgs) (*ShopCustomer, error)

	DeleteCustomer(ctx context.Context, ID int64, shopID int64) (deleted int, _ error)

	BatchSetCustomersStatus(ctx context.Context, IDs []int64, shopID int64, status int32) (*meta.UpdatedResponse, error)

	AddCustomersToGroup(ctx context.Context, _ *AddCustomerToGroupArgs) (updateed int, _ error)

	RemoveCustomersFromGroup(ctx context.Context, _ *RemoveCustomerOutOfGroupArgs) (deleted int, _ error)

	CreateCustomerGroup(ctx context.Context, _ *CreateCustomerGroupArgs) (*ShopCustomerGroup, error)

	UpdateCustomerGroup(ctx context.Context, _ *UpdateCustomerGroupArgs) (*ShopCustomerGroup, error)
}

type QueryService interface {
	GetCustomerByID(context.Context, *shopping.IDQueryShopArg) (*ShopCustomer, error)

	// unused
	GetCustomerByCode(ctx context.Context, code string, shopID int64) (*ShopCustomer, error)

	GetCustomerByPhone(ctx context.Context, phone string, shopID int64) (*ShopCustomer, error)

	GetCustomerByEmail(ctx context.Context, email string, shopID int64) (*ShopCustomer, error)

	ListCustomers(context.Context, *shopping.ListQueryShopArgs) (*CustomersResponse, error)

	ListCustomersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*CustomersResponse, error)

	GetCustomerGroup(ctx context.Context, _ *GetCustomerGroupArgs) (*ShopCustomerGroup, error)

	ListCustomerGroups(ctx context.Context, _ *ListCustomerGroupArgs) (*CustomerGroupsResponse, error)
}

//-- queries --//
type GetCustomerGroupArgs struct {
	ID int64
}

type ListCustomerGroupArgs struct {
	Paging  meta.Paging
	Filters meta.Filters
}

type CustomerGroupsResponse struct {
	CustomerGroups []*ShopCustomerGroup
	Count          int32
	Paging         meta.PageInfo
}

type CustomersResponse struct {
	Customers []*ShopCustomer
	Count     int32
	Paging    meta.PageInfo
}

//-- commands --//

type CreateCustomerGroupArgs struct {
	Name string
}

// +convert:create=ShopCustomer
type CreateCustomerArgs struct {
	ShopID   int64
	FullName string
	Gender   string
	Type     string
	Birthday string
	Note     string
	Phone    string
	Email    string
}

// +convert:update=ShopCustomer(ID,ShopID)
type UpdateCustomerArgs struct {
	ID       int64
	ShopID   int64
	FullName NullString
	Gender   NullString
	Type     NullString
	Birthday NullString
	Note     NullString
	Phone    NullString
	Email    NullString
}

type BatchSetCustomersStatusArgs struct {
	IDs    []int64
	ShopID int64
	Status int32
}

type AddCustomerToGroupArgs struct {
	GroupID     int64
	CustomerIDs []int64
}

type RemoveCustomerOutOfGroupArgs struct {
	GroupID     int64
	CustomerIDs []int64
}

type UpdateCustomerGroupArgs struct {
	ID   int64
	Name string
}
