package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type CustomerService struct{}

func (s *CustomerService) Clone() *CustomerService { res := *s; return &res }

func (s *CustomerService) GetCustomer(ctx context.Context, r *GetCustomerEndpoint) error {
	resp, err := shopping.GetCustomer(ctx, r.Context.Shop.ID, r.GetCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerService) ListCustomers(ctx context.Context, r *ListCustomersEndpoint) error {
	resp, err := shopping.ListCustomers(ctx, r.Context.Shop.ID, r.ListCustomersRequest)
	r.Result = resp
	return err
}

func (s *CustomerService) CreateCustomer(ctx context.Context, r *CreateCustomerEndpoint) error {
	resp, err := shopping.CreateCustomer(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, r *UpdateCustomerEndpoint) error {
	resp, err := shopping.UpdateCustomer(ctx, r.Context.Shop.ID, r.UpdateCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, r *DeleteCustomerEndpoint) error {
	resp, err := shopping.DeleteCustomer(ctx, r.Context.Shop.ID, r.DeleteCustomerRequest)
	r.Result = resp
	return err
}
