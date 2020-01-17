package partner

import (
	"context"

	"etop.vn/backend/pkg/etop/apix/shopping"
)

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
	resp, err := shopping.CreateAddress(ctx, r.Context.Shop.ID, r.CreateCustomerAddressRequest)
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

func (s *CustomerGroupRelationshipService) ListRelationships(ctx context.Context, r *CustomerGroupListRelationshipsEndpoint) error {
	panic("TODO")
}

func (s *CustomerGroupRelationshipService) CreateRelationship(ctx context.Context, r *CustomerGroupCreateRelationshipEndpoint) error {
	resp, err := shopping.CreateRelationshipGroupCustomer(ctx, r.Context.Shop.ID, r.AddCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupRelationshipService) DeleteRelationship(ctx context.Context, r *CustomerGroupDeleteRelationshipEndpoint) error {
	resp, err := shopping.DeleteRelationshipGroupCustomer(ctx, r.Context.Shop.ID, r.RemoveCustomerRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) GetGroup(ctx context.Context, r *GetGroupEndpoint) error {
	resp, err := shopping.GetGroup(ctx, r.Context.Shop.ID, r.IDRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) ListGroups(ctx context.Context, r *ListGroupsEndpoint) error {
	resp, err := shopping.ListGroups(ctx, r.Context.Shop.ID, r.ListCustomerGroupsRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) CreateGroup(ctx context.Context, r *CreateGroupEndpoint) error {
	resp, err := shopping.CreateGroup(ctx, r.Context.Shop.ID, r.CreateCustomerGroupRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) UpdateGroup(ctx context.Context, r *UpdateGroupEndpoint) error {
	resp, err := shopping.UpdateGroup(ctx, r.Context.Shop.ID, r.UpdateCustomerGroupRequest)
	r.Result = resp
	return err
}

func (s *CustomerGroupService) DeleteGroup(ctx context.Context, r *DeleteGroupEndpoint) error {
	resp, err := shopping.DeleteGroup(ctx, r.Context.Shop.ID, r.IDRequest)
	r.Result = resp
	return err
}
