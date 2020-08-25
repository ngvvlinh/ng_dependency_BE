package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type CustomerService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *CustomerService) Clone() api.CustomerService { res := *s; return &res }

func (s *CustomerService) GetCustomer(ctx context.Context, r *externaltypes.GetCustomerRequest) (*externaltypes.Customer, error) {
	resp, err := s.Shopping.GetCustomer(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerService) ListCustomers(ctx context.Context, r *externaltypes.ListCustomersRequest) (*externaltypes.CustomersResponse, error) {
	resp, err := s.Shopping.ListCustomers(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerService) CreateCustomer(ctx context.Context, r *externaltypes.CreateCustomerRequest) (*externaltypes.Customer, error) {
	resp, err := s.Shopping.CreateCustomer(ctx, s.SS.Shop().ID, s.SS.Claim().AuthPartnerID, r)
	return resp, err
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, r *externaltypes.UpdateCustomerRequest) (*externaltypes.Customer, error) {
	resp, err := s.Shopping.UpdateCustomer(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, r *externaltypes.DeleteCustomerRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteCustomer(ctx, s.SS.Shop().ID, r)
	return resp, err
}
