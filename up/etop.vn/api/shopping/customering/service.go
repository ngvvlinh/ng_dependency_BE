package customering

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/types/etc/gender"
	. "etop.vn/capi/dot"
	dot "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCustomer(ctx context.Context, _ *CreateCustomerArgs) (*ShopCustomer, error)

	UpdateCustomer(ctx context.Context, _ *UpdateCustomerArgs) (*ShopCustomer, error)

	DeleteCustomer(ctx context.Context, ID dot.ID, shopID dot.ID) (deleted int, _ error)

	BatchSetCustomersStatus(ctx context.Context, IDs []dot.ID, shopID dot.ID, status int) (*meta.UpdatedResponse, error)

	AddCustomersToGroup(ctx context.Context, _ *AddCustomerToGroupArgs) (updateed int, _ error)

	RemoveCustomersFromGroup(ctx context.Context, _ *RemoveCustomerOutOfGroupArgs) (deleted int, _ error)

	CreateCustomerGroup(ctx context.Context, _ *CreateCustomerGroupArgs) (*ShopCustomerGroup, error)

	UpdateCustomerGroup(ctx context.Context, _ *UpdateCustomerGroupArgs) (*ShopCustomerGroup, error)
}

type QueryService interface {
	GetCustomerByID(context.Context, *shopping.IDQueryShopArg) (*ShopCustomer, error)

	// unused
	GetCustomerByCode(ctx context.Context, code string, shopID dot.ID) (*ShopCustomer, error)

	GetCustomerByPhone(ctx context.Context, phone string, shopID dot.ID) (*ShopCustomer, error)

	GetCustomerByEmail(ctx context.Context, email string, shopID dot.ID) (*ShopCustomer, error)

	GetCustomerIndependentByShop(ctx context.Context, _ *GetCustomerIndependentByShop) (*ShopCustomer, error)

	ListCustomers(context.Context, *shopping.ListQueryShopArgs) (*CustomersResponse, error)

	ListCustomersByIDs(context.Context, *ListCustomerByIDsArgs) (*CustomersResponse, error)

	GetCustomerGroup(ctx context.Context, _ *GetCustomerGroupArgs) (*ShopCustomerGroup, error)

	ListCustomerGroups(ctx context.Context, _ *ListCustomerGroupArgs) (*CustomerGroupsResponse, error)
}

//-- queries --//
type GetCustomerGroupArgs struct {
	ID dot.ID
}

type ListCustomerByIDsArgs struct {
	IDs     []dot.ID
	ShopIDs []dot.ID
	ShopID  dot.ID
	Paging  meta.Paging
}

type ListCustomerGroupArgs struct {
	Paging  meta.Paging
	Filters meta.Filters
}

type CustomerGroupsResponse struct {
	CustomerGroups []*ShopCustomerGroup
	Paging         meta.PageInfo
}

type CustomersResponse struct {
	Customers []*ShopCustomer
	Paging    meta.PageInfo
}

//-- commands --//

type CreateCustomerGroupArgs struct {
	Name   string
	ShopID dot.ID
}

// +convert:create=ShopCustomer
type CreateCustomerArgs struct {
	ShopID   dot.ID
	FullName string
	Gender   gender.Gender
	Type     customer_type.CustomerType
	Birthday string
	Note     string
	Phone    string
	Email    string
}

// +convert:update=ShopCustomer(ID,ShopID)
type UpdateCustomerArgs struct {
	ID       dot.ID
	ShopID   dot.ID
	FullName NullString
	Gender   gender.NullGender
	Type     customer_type.CustomerType
	Birthday NullString
	Note     NullString
	Phone    NullString
	Email    NullString
}

type BatchSetCustomersStatusArgs struct {
	IDs    []dot.ID
	ShopID dot.ID
	Status int
}

type AddCustomerToGroupArgs struct {
	GroupID     dot.ID
	CustomerIDs []dot.ID
	ShopID      dot.ID
}

type RemoveCustomerOutOfGroupArgs struct {
	GroupID     dot.ID
	CustomerIDs []dot.ID
	ShopID      dot.ID
}

type UpdateCustomerGroupArgs struct {
	ID   dot.ID
	Name string
}

type GetCustomerIndependentByShop struct {
	ShopID dot.ID
}
