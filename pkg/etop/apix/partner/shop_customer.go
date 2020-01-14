package partner

import (
	"context"

	"etop.vn/api/main/ordering/types"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/customering/customer_type"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/api/top/types/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
)

func (s *CustomerService) GetCustomer(ctx context.Context, r *GetCustomerEndpoint) error {
	query := &customering.GetCustomerQuery{
		ID:         r.Id,
		ShopID:     r.Context.Shop.ID,
		Code:       r.Code,
		ExternalID: r.ExternalId,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShopCustomer(query.Result)
	return nil
}

func (s *CustomerService) ListCustomers(ctx context.Context, r *ListCustomersEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	var IDs []dot.ID
	if len(r.Filter.ID) != 0 {
		IDs = r.Filter.ID
	}

	query := &customering.ListCustomersByIDsQuery{
		ShopID: r.Context.Shop.ID,
		IDs:    IDs,
		Paging: *paging,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &externaltypes.CustomersResponse{
		Customers: convertpb.PbShopCustomers(query.Result.Customers),
		Paging:    convertpb.PbPageInfo(r.Paging, &query.Result.Paging),
	}
	return nil
}

func (s *CustomerService) CreateCustomer(ctx context.Context, r *CreateCustomerEndpoint) error {
	cmd := &customering.CreateCustomerCommand{
		ExternalID:   r.ExternalId,
		ExternalCode: r.ExternalCode,
		PartnerID:    r.Context.AuthPartnerID,
		ShopID:       r.Context.Shop.ID,
		FullName:     r.FullName,
		Gender:       r.Gender,
		Type:         r.Type,
		Birthday:     r.Birthday,
		Note:         r.Note,
		Phone:        r.Phone,
		Email:        r.Email,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopCustomer(cmd.Result)
	return nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, r *UpdateCustomerEndpoint) error {
	cmd := &customering.UpdateCustomerCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Gender:   r.Gender,
		Birthday: r.Birthday,
		Note:     r.Note,
		Phone:    r.Phone,
		Email:    r.Email,
	}
	if r.Type.Valid && r.Type.String != "" {
		customerType, ok := customer_type.ParseCustomerType(r.Type.String)
		if !ok {
			return cm.Errorf(cm.InvalidArgument, nil, "type không hợp lệ")
		}
		cmd.Type = customerType
	}

	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopCustomer(cmd.Result)
	return nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, r *DeleteCustomerEndpoint) error {
	cmd := &customering.DeleteCustomerCommand{
		ID:         r.Id,
		ShopID:     r.Context.Shop.ID,
		ExternalID: r.ExternalId,
		Code:       r.Code,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}

func (s *CustomerAddressService) ListAddresses(ctx context.Context, r *ListAddressesEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	var IDs []dot.ID
	if len(r.Filter.CustomerId) != 0 {
		IDs = r.Filter.CustomerId
	}
	query := &addressing.ListAddressesByTraderIDsQuery{
		TraderIDs: IDs,
		ShopID:    r.Context.Shop.ID,
		Paging:    *paging,
	}
	if err := addressQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &externaltypes.CustomerAddressesResponse{
		CustomerAddresses: convertpb.PbShopTraderAddresses(ctx, query.Result.ShopTraderAddresses, locationQuery),
	}
	return nil
}

func (s *CustomerAddressService) GetAddress(ctx context.Context, r *GetAddressEndpoint) error {
	query := &addressing.GetAddressByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := addressQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShopTraderAddress(ctx, query.Result, locationQuery)
	return nil
}

func (s *CustomerAddressService) CreateAddress(ctx context.Context, r *CreateAddressEndpoint) error {
	var coordinates *types.Coordinates
	if r.Coordinates != nil {
		coordinates = &types.Coordinates{
			Latitude:  r.Coordinates.Latitude,
			Longitude: r.Coordinates.Longitude,
		}
	}
	cmd := &addressing.CreateAddressCommand{
		ShopID:       r.Context.Shop.ID,
		TraderID:     r.CustomerId,
		FullName:     r.FullName,
		Phone:        r.Phone,
		Email:        r.Email,
		Company:      r.Company,
		Address1:     r.Address1,
		Address2:     r.Address2,
		DistrictCode: r.DistrictCode,
		WardCode:     r.WardCode,
		Position:     r.Position,
		IsDefault:    false,
		Coordinates:  coordinates,
	}
	if err := addressAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopTraderAddress(ctx, cmd.Result, locationQuery)
	return nil
}

func (s *CustomerAddressService) UpdateAddress(ctx context.Context, r *UpdateAddressEndpoint) error {
	var coordinates *types.Coordinates
	if r.Coordinates != nil {
		coordinates = &types.Coordinates{
			Latitude:  r.Coordinates.Latitude,
			Longitude: r.Coordinates.Longitude,
		}
	}
	cmd := &addressing.UpdateAddressCommand{
		ID:           r.Id,
		ShopID:       r.Context.Shop.ID,
		FullName:     r.FullName,
		Phone:        r.Phone,
		Email:        r.Email,
		Company:      r.Company,
		Address1:     r.Address1,
		Address2:     r.Address2,
		DistrictCode: r.DistrictCode,
		WardCode:     r.WardCode,
		Position:     r.Position,
		Coordinates:  coordinates,
	}
	if err := addressAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopTraderAddress(ctx, cmd.Result, locationQuery)
	return nil
}

func (s *CustomerAddressService) DeleteAddress(ctx context.Context, r *DeleteAddressEndpoint) error {
	cmd := &addressing.DeleteAddressCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := addressAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}

func (s *CustomerGroupRelationshipService) ListRelationships(ctx context.Context, r *CustomerGroupListRelationshipsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	query := &customering.ListCustomerGroupsCustomersQuery{
		Paging:      *paging,
		CustomerIDs: r.Filter.CustomerID,
		GroupIDs:    r.Filter.GroupID,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &externaltypes.CustomerGroupRelationshipsResponse{
		Relationships: convertpb.PbRelationships(query.Result.CustomerGroupsCustomers),
		Paging:        convertpb.PbPageInfo(r.Paging, &query.Result.Paging),
	}
	return nil
}

func (s *CustomerGroupRelationshipService) CreateRelationship(ctx context.Context, r *CustomerGroupCreateRelationshipEndpoint) error {
	cmd := &customering.AddCustomersToGroupCommand{
		ShopID:      r.Context.Shop.ID,
		GroupID:     r.GroupID,
		CustomerIDs: []dot.ID{r.CustomerID},
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}

func (s *CustomerGroupRelationshipService) DeleteRelationship(ctx context.Context, r *CustomerGroupDeleteRelationshipEndpoint) error {
	cmd := &customering.RemoveCustomersFromGroupCommand{
		ShopID:      r.Context.Shop.ID,
		GroupID:     r.GroupID,
		CustomerIDs: []dot.ID{r.CustomerID},
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}

func (s *CustomerGroupService) GetGroup(ctx context.Context, r *GetGroupEndpoint) error {
	query := &customering.GetCustomerGroupQuery{
		ID: r.Id,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbCustomerGroup(query.Result)
	return nil
}

func (s *CustomerGroupService) ListGroups(ctx context.Context, r *ListGroupsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	query := &customering.ListCustomerGroupsQuery{
		Paging: *paging,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &externaltypes.CustomerGroupsResponse{
		CustomerGroups: convertpb.PbCustomerGroups(query.Result.CustomerGroups),
		Paging:         convertpb.PbPageInfo(r.Paging, &query.Result.Paging),
	}
	return nil
}

func (s *CustomerGroupService) CreateGroup(ctx context.Context, r *CreateGroupEndpoint) error {
	if r.Name == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tên không được rỗng.")
	}
	cmd := &customering.CreateCustomerGroupCommand{
		ShopID: r.Context.Shop.ID,
		Name:   r.Name,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustomerGroup(cmd.Result)
	return nil
}

func (s *CustomerGroupService) UpdateGroup(ctx context.Context, r *UpdateGroupEndpoint) error {
	if r.Name.String == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tên không được rỗng.")
	}
	cmd := &customering.UpdateCustomerGroupCommand{
		ID:   r.GroupId,
		Name: r.Name.String,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbCustomerGroup(cmd.Result)
	return nil
}

func (s *CustomerGroupService) DeleteGroup(ctx context.Context, r *DeleteGroupEndpoint) error {
	cmd := &customering.DeleteGroupCommand{
		GroupID: r.Id,
		ShopID:  r.Context.Shop.ID,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}
