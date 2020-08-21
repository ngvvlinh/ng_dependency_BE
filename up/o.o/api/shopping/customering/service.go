package customering

import (
	"context"

	"o.o/api/meta"
	"o.o/api/shopping"
	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/gender"
	. "o.o/capi/dot"
	dot "o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateCustomer(ctx context.Context, _ *CreateCustomerArgs) (*ShopCustomer, error)

	UpdateCustomer(ctx context.Context, _ *UpdateCustomerArgs) (*ShopCustomer, error)

	DeleteCustomer(ctx context.Context, _ *DeleteCustomerArgs) (deleted int, _ error)

	BatchSetCustomersStatus(ctx context.Context, IDs []dot.ID, shopID dot.ID, status int) (*meta.UpdatedResponse, error)

	AddCustomersToGroup(ctx context.Context, _ *AddCustomerToGroupArgs) (updateed int, _ error)

	RemoveCustomersFromGroup(ctx context.Context, _ *RemoveCustomerOutOfGroupArgs) (deleted int, _ error)

	CreateCustomerGroup(ctx context.Context, _ *CreateCustomerGroupArgs) (*ShopCustomerGroup, error)

	UpdateCustomerGroup(ctx context.Context, _ *UpdateCustomerGroupArgs) (*ShopCustomerGroup, error)

	DeleteGroup(ctx context.Context, _ *DeleteGroupArgs) (deleted int, _ error)
}

type QueryService interface {
	GetCustomer(context.Context, *GetCustomerArgs) (*ShopCustomer, error)

	GetCustomerByID(context.Context, *shopping.IDQueryShopArg) (*ShopCustomer, error)

	// unused
	GetCustomerByCode(ctx context.Context, code string, shopID dot.ID) (*ShopCustomer, error)

	GetCustomerByPhone(ctx context.Context, phone string, shopID dot.ID) (*ShopCustomer, error)

	GetCustomerByEmail(ctx context.Context, email string, shopID dot.ID) (*ShopCustomer, error)

	ListCustomers(context.Context, *shopping.ListQueryShopArgs) (*CustomersResponse, error)

	ListCustomersByIDs(context.Context, *ListCustomerByIDsArgs) (*CustomersResponse, error)

	ListCustomersByPhones(ctx context.Context, shopID dot.ID, phones []string) ([]*ShopCustomer, error)

	GetCustomerGroup(ctx context.Context, _ *GetCustomerGroupArgs) (*ShopCustomerGroup, error)

	ListCustomerGroups(ctx context.Context, _ *ListCustomerGroupArgs) (*CustomerGroupsResponse, error)

	ListCustomerGroupsCustomers(ctx context.Context, _ *ListCustomerGroupsCustomersArgs) (*CustomerGroupsCustomersResponse, error)

	GetCustomerIndependent(context.Context, *meta.Empty) (*ShopCustomer, error)
}

//-- queries --//

type ListCustomerGroupsCustomersArgs struct {
	CustomerIDs []dot.ID
	GroupIDs    []dot.ID

	Paging meta.Paging

	IncludeDeleted bool
}

type GetCustomerArgs struct {
	ID         dot.ID
	ShopID     dot.ID
	Code       string
	ExternalID string
}

type GetCustomerGroupArgs struct {
	ID dot.ID
}

type ListCustomerByIDsArgs struct {
	IDs     []dot.ID
	ShopIDs []dot.ID
	ShopID  dot.ID
	Paging  meta.Paging

	IncludeDeleted bool
}

type ListCustomerGroupArgs struct {
	Paging  meta.Paging
	Filters meta.Filters

	IncludeDeleted bool
}

type CustomerGroupsCustomersResponse struct {
	CustomerGroupsCustomers []*CustomerGroupCustomer
	Paging                  meta.PageInfo
}

type CustomerGroupCustomer struct {
	CustomerID dot.ID
	GroupID    dot.ID
	Deleted    bool
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

type DeleteCustomerArgs struct {
	ID     dot.ID
	ShopID dot.ID

	ExternalID string
	Code       string
}

type CreateCustomerGroupArgs struct {
	Name      string
	PartnerID dot.ID
	ShopID    dot.ID
}

// +convert:create=ShopCustomer
type CreateCustomerArgs struct {
	// @Optional
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID

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

type DeleteGroupArgs struct {
	ShopID  dot.ID
	GroupID dot.ID
}
