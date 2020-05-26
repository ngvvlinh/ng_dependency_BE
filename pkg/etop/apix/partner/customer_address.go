package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type CustomerAddressService struct{}

func (s *CustomerAddressService) Clone() *CustomerAddressService { res := *s; return &res }

func (s *CustomerAddressService) ListAddresses(ctx context.Context, r *ListAddressesEndpoint) error {
	resp, err := shopping.ListAddresses(ctx, r.Context.Shop.ID, r.ListCustomerAddressesRequest)
	r.Result = resp
	return err
}

func (s *CustomerAddressService) GetAddress(ctx context.Context, r *GetAddressEndpoint) error {
	resp, err := shopping.GetAddress(ctx, r.Context.Shop.ID, r.OrderIDRequest)
	r.Result = resp
	return err
}

func (s *CustomerAddressService) CreateAddress(ctx context.Context, r *CreateAddressEndpoint) error {
	resp, err := shopping.CreateAddress(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateCustomerAddressRequest)
	r.Result = resp
	return err
}

func (s *CustomerAddressService) UpdateAddress(ctx context.Context, r *UpdateAddressEndpoint) error {
	resp, err := shopping.UpdateAddress(ctx, r.Context.Shop.ID, r.UpdateCustomerAddressRequest)
	r.Result = resp
	return err
}

func (s *CustomerAddressService) DeleteAddress(ctx context.Context, r *DeleteAddressEndpoint) error {
	resp, err := shopping.DeleteAddress(ctx, r.Context.Shop.ID, r.IDRequest)
	r.Result = resp
	return err
}
