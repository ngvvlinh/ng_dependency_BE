package customering

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

type Aggregate interface {
	CreateCustomer(context.Context, *CreateCustomerArgs) (*ShopCustomer, error)

	UpdateCustomer(context.Context, *UpdateCustomerArgs) (*ShopCustomer, error)

	DeleteCustomer(ctx context.Context, ID int64, shopID int64) error

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
