package shopping

import (
	"context"

	"etop.vn/api/main/ordering/types"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/customering/customer_type"
	externaltypes "etop.vn/api/top/external/types"
	cm "etop.vn/api/top/types/common"
	common "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
)

func GetCustomer(ctx context.Context, shopID dot.ID, request *externaltypes.GetCustomerRequest) (*externaltypes.Customer, error) {
	query := &customering.GetCustomerQuery{
		ID:         request.Id,
		ShopID:     shopID,
		Code:       request.Code,
		ExternalID: request.ExternalId,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopCustomer(query.Result), nil
}

func ListCustomers(ctx context.Context, shopID dot.ID, request *externaltypes.ListCustomersRequest) (*externaltypes.CustomersResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	var IDs []dot.ID
	if len(request.Filter.ID) != 0 {
		IDs = request.Filter.ID
	}

	query := &customering.ListCustomersByIDsQuery{
		ShopID:         shopID,
		IDs:            IDs,
		Paging:         *paging,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	return &externaltypes.CustomersResponse{
		Customers: convertpb.PbShopCustomers(query.Result.Customers),
		Paging:    convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func CreateCustomer(ctx context.Context, shopID dot.ID, partnerID dot.ID, request *externaltypes.CreateCustomerRequest) (*externaltypes.Customer, error) {
	cmd := &customering.CreateCustomerCommand{
		ExternalID:   request.ExternalId,
		ExternalCode: request.ExternalCode,
		PartnerID:    partnerID,
		ShopID:       shopID,
		FullName:     request.FullName,
		Gender:       request.Gender,
		Type:         request.Type,
		Birthday:     request.Birthday,
		Note:         request.Note,
		Phone:        request.Phone,
		Email:        request.Email,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopCustomer(cmd.Result), nil
}

func UpdateCustomer(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateCustomerRequest) (*externaltypes.Customer, error) {
	cmd := &customering.UpdateCustomerCommand{
		ID:       request.Id,
		ShopID:   shopID,
		FullName: request.FullName,
		Gender:   request.Gender,
		Birthday: request.Birthday,
		Note:     request.Note,
		Phone:    request.Phone,
		Email:    request.Email,
	}
	if request.Type.Valid && request.Type.String != "" {
		customerType, ok := customer_type.ParseCustomerType(request.Type.String)
		if !ok {
			return nil, common.Errorf(common.InvalidArgument, nil, "type không hợp lệ")
		}
		cmd.Type = customerType
	}

	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopCustomer(cmd.Result), nil
}

func DeleteCustomer(ctx context.Context, shopID dot.ID, request *externaltypes.DeleteCustomerRequest) (*cm.Empty, error) {
	cmd := &customering.DeleteCustomerCommand{
		ID:         request.Id,
		ShopID:     shopID,
		ExternalID: request.ExternalId,
		Code:       request.Code,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &cm.Empty{}, nil
}

func ListAddresses(ctx context.Context, shopID dot.ID, request *externaltypes.ListCustomerAddressesRequest) (*externaltypes.CustomerAddressesResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	var IDs []dot.ID
	if len(request.Filter.CustomerId) != 0 {
		IDs = request.Filter.CustomerId
	}
	query := &addressing.ListAddressesByTraderIDsQuery{
		TraderIDs:      IDs,
		ShopID:         shopID,
		Paging:         *paging,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := addressQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.CustomerAddressesResponse{
		CustomerAddresses: convertpb.PbShopTraderAddresses(ctx, query.Result.ShopTraderAddresses, locationQuery),
	}, nil
}

func GetAddress(ctx context.Context, shopID dot.ID, request *externaltypes.OrderIDRequest) (*externaltypes.CustomerAddress, error) {
	query := &addressing.GetAddressByIDQuery{
		ID:     request.Id,
		ShopID: shopID,
	}
	if err := addressQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopTraderAddress(ctx, query.Result, locationQuery), nil
}

func CreateAddress(ctx context.Context, shopID dot.ID, request *externaltypes.CreateCustomerAddressRequest) (*externaltypes.CustomerAddress, error) {
	var coordinates *types.Coordinates
	if request.Coordinates != nil {
		coordinates = &types.Coordinates{
			Latitude:  request.Coordinates.Latitude,
			Longitude: request.Coordinates.Longitude,
		}
	}
	cmd := &addressing.CreateAddressCommand{
		ShopID:       shopID,
		TraderID:     request.CustomerId,
		FullName:     request.FullName,
		Phone:        request.Phone,
		Email:        request.Email,
		Company:      request.Company,
		Address1:     request.Address1,
		Address2:     request.Address2,
		DistrictCode: request.DistrictCode,
		WardCode:     request.WardCode,
		Position:     request.Position,
		IsDefault:    false,
		Coordinates:  coordinates,
	}
	if err := addressAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopTraderAddress(ctx, cmd.Result, locationQuery), nil
}

func UpdateAddress(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateCustomerAddressRequest) (*externaltypes.CustomerAddress, error) {
	var coordinates *types.Coordinates
	if request.Coordinates != nil {
		coordinates = &types.Coordinates{
			Latitude:  request.Coordinates.Latitude,
			Longitude: request.Coordinates.Longitude,
		}
	}
	cmd := &addressing.UpdateAddressCommand{
		ID:           request.Id,
		ShopID:       shopID,
		FullName:     request.FullName,
		Phone:        request.Phone,
		Email:        request.Email,
		Company:      request.Company,
		Address1:     request.Address1,
		Address2:     request.Address2,
		DistrictCode: request.DistrictCode,
		WardCode:     request.WardCode,
		Position:     request.Position,
		Coordinates:  coordinates,
	}
	if err := addressAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopTraderAddress(ctx, cmd.Result, locationQuery), nil
}

func DeleteAddress(ctx context.Context, shopID dot.ID, request *cm.IDRequest) (*cm.Empty, error) {
	cmd := &addressing.DeleteAddressCommand{
		ID:     request.Id,
		ShopID: shopID,
	}
	if err := addressAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &cm.Empty{}, nil
}

func ListRelationshipsGroupCustomer(ctx context.Context, shopID dot.ID, request *externaltypes.ListCustomerGroupRelationshipsRequest) (*externaltypes.CustomerGroupRelationshipsResponse, error) {
	// TODO: add cursor paging
	query := &customering.ListCustomerGroupsCustomersQuery{
		CustomerIDs:    request.Filter.CustomerID,
		GroupIDs:       request.Filter.GroupID,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.CustomerGroupRelationshipsResponse{
		Relationships: convertpb.PbRelationships(query.Result.CustomerGroupsCustomers),
	}, nil
}

func CreateRelationshipGroupCustomer(ctx context.Context, shopID dot.ID, request *externaltypes.AddCustomerRequest) (*cm.Empty, error) {
	cmd := &customering.AddCustomersToGroupCommand{
		ShopID:      shopID,
		GroupID:     request.GroupID,
		CustomerIDs: []dot.ID{request.CustomerID},
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &cm.Empty{}, nil
}

func DeleteRelationshipGroupCustomer(ctx context.Context, shopID dot.ID, request *externaltypes.RemoveCustomerRequest) (*cm.Empty, error) {
	cmd := &customering.RemoveCustomersFromGroupCommand{
		ShopID:      shopID,
		GroupID:     request.GroupID,
		CustomerIDs: []dot.ID{request.CustomerID},
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &cm.Empty{}, nil
}

func GetGroup(ctx context.Context, shopID dot.ID, request *cm.IDRequest) (*externaltypes.CustomerGroup, error) {
	query := &customering.GetCustomerGroupQuery{
		ID: request.Id,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbCustomerGroup(query.Result), nil
}

func ListGroups(ctx context.Context, shopID dot.ID, request *externaltypes.ListCustomerGroupsRequest) (*externaltypes.CustomerGroupsResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}

	query := &customering.ListCustomerGroupsQuery{
		Paging:         *paging,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.CustomerGroupsResponse{
		CustomerGroups: convertpb.PbCustomerGroups(query.Result.CustomerGroups),
		Paging:         convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func CreateGroup(ctx context.Context, shopID dot.ID, request *externaltypes.CreateCustomerGroupRequest) (*externaltypes.CustomerGroup, error) {
	if request.Name == "" {
		return nil, common.Errorf(common.InvalidArgument, nil, "Tên không được rỗng.")
	}
	cmd := &customering.CreateCustomerGroupCommand{
		ShopID: shopID,
		Name:   request.Name,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbCustomerGroup(cmd.Result), nil
}

func UpdateGroup(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateCustomerGroupRequest) (*externaltypes.CustomerGroup, error) {
	if request.Name.String == "" {
		return nil, common.Errorf(common.InvalidArgument, nil, "Tên không được rỗng.")
	}
	cmd := &customering.UpdateCustomerGroupCommand{
		ID:   request.GroupId,
		Name: request.Name.String,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbCustomerGroup(cmd.Result), nil
}

func DeleteGroup(ctx context.Context, shopID dot.ID, request *cm.IDRequest) (*cm.Empty, error) {
	cmd := &customering.DeleteGroupCommand{
		GroupID: request.Id,
		ShopID:  shopID,
	}
	if err := customerAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &cm.Empty{}, nil
}
