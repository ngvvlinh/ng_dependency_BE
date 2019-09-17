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
}

type QueryService interface {
	GetCustomerByID(context.Context, *shopping.IDQueryShopArg) (*ShopCustomer, error)

	ListCustomers(context.Context, *shopping.ListQueryShopArgs) (*CustomersResponse, error)

	ListCustomersByIDs(context.Context, *shopping.IDsQueryShopArgs) (*CustomersResponse, error)
}

//-- queries --//

type CustomersResponse struct {
	Customers []*ShopCustomer
	Count     int32
	Paging    meta.PageInfo
}

//-- commands --//

// +convert:create=ShopCustomer
type CreateCustomerArgs struct {
	ShopID   int64
	Code     string
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
	Code     NullString
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
