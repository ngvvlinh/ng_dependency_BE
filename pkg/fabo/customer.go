package fabo

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/fabo"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	convertpbfabo "o.o/backend/pkg/fabo/convertpb"
	"o.o/capi/dot"
)

type CustomerService struct {
	session.Session

	CustomerQuery  customering.QueryBus
	FBUseringQuery fbusering.QueryBus
	FBUseringAggr  fbusering.CommandBus
}

func (s *CustomerService) Clone() fabo.CustomerService {
	res := *s
	return &res
}

func (s *CustomerService) CreateFbUserCustomer(ctx context.Context, request *fabo.CreateFbUserCustomerRequest) (*fabo.FbUserWithCustomer, error) {
	shopID := s.SS.Shop().ID
	cmd := &fbusering.CreateFbExternalUserShopCustomerCommand{
		ShopID:     shopID,
		ExternalID: request.ExternalID,
		CustomerID: request.CustomerID,
	}
	err := s.FBUseringAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return convertpbfabo.PbFbUserWithCustomer(cmd.Result.FbExternalUser, cmd.Result.ShopCustomer), nil
}

func (s *CustomerService) GetFbUser(ctx context.Context, request *fabo.GetFbUserRequest) (*fabo.FbUserWithCustomer, error) {
	shopID := s.SS.Shop().ID
	query := &fbusering.GetFbExternalUserWithCustomerByExternalIDQuery{
		ExternalID: request.ExternalID,
		ShopID:     shopID,
	}
	if err := s.FBUseringQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpbfabo.PbFbUserWithCustomer(query.Result.FbExternalUser, query.Result.ShopCustomer), nil
}

func (s *CustomerService) ListFbUsers(ctx context.Context, request *fabo.ListFbUsersRequest) (*fabo.ListFbUsersResponse, error) {
	shopID := s.SS.Shop().ID
	query := &fbusering.ListFbExternalUsersQuery{
		CustomerID: request.CustomerID,
		ShopID:     shopID,
	}
	if err := s.FBUseringQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var result = &fabo.ListFbUsersResponse{}
	result.FbUsers = convertpbfabo.PbExternalUsersWithCustomer(query.Result)
	return result, nil
}

func (s *CustomerService) ListCustomersWithFbUsers(ctx context.Context, request *fabo.ListCustomersWithFbUsersRequest) (resp *fabo.ListCustomersWithFbUsersResponse, err error) {
	paging := cmapi.CMPaging(request.Paging)

	if !request.GetAll.Valid {
		request.GetAll = dot.Bool(true)
	}

	if request.GetAll.Apply(false) {
		resp, err = s.getAllCustomers(ctx, paging, request)
	} else {
		resp, err = s.getCustomers(ctx, paging, request)
	}
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (s *CustomerService) UpdateTags(ctx context.Context, request *fabo.UpdateUserTagsRequest) (resp *fabo.UpdateUserTagResponse, err error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	cmdUpdateTag := &fbusering.UpdateShopUserTagsCommand{
		ShopID:           s.SS.Shop().ID,
		TagIDs:           request.TagIDs,
		FbExternalUserID: request.FbExternalUserID,
	}
	if err := s.FBUseringAggr.Dispatch(ctx, cmdUpdateTag); err != nil {
		return nil, err
	}

	return &fabo.UpdateUserTagResponse{TagIDs: cmdUpdateTag.Result.TagIDs}, nil
}

func (s *CustomerService) getAllCustomers(ctx context.Context, paging *cm.Paging, request *fabo.ListCustomersWithFbUsersRequest) (*fabo.ListCustomersWithFbUsersResponse, error) {
	queryCustomerIndependent := &customering.GetCustomerIndependentQuery{}
	if err := s.CustomerQuery.Dispatch(ctx, queryCustomerIndependent); err != nil {
		return nil, err
	}

	var customers []*fabo.CustomerWithFbUserAvatars
	customers = append(customers, convertpbfabo.PbCustomerWithFbUser(
		&fbusering.ShopCustomerWithFbExternalUser{
			ShopCustomer: queryCustomerIndependent.Result,
		},
	))

	result := &fabo.ListCustomersWithFbUsersResponse{
		Paging: cmapi.PbPageInfo(paging),
	}
	if paging.Limit == 1 && paging.Offset == 0 {
		result.Customers = customers
		return result, nil
	}
	if paging.Offset == 0 {
		paging.Limit--
		resp, err := s.getCustomers(ctx, paging, request)
		if err != nil {
			return nil, err
		}
		customers = append(customers, resp.Customers...)
		result.Customers = customers
	} else {
		paging.Offset--
		return s.getCustomers(ctx, paging, request)
	}
	return result, nil
}

func (s *CustomerService) getCustomers(ctx context.Context, paging *cm.Paging, request *fabo.ListCustomersWithFbUsersRequest) (*fabo.ListCustomersWithFbUsersResponse, error) {
	result := &fabo.ListCustomersWithFbUsersResponse{}
	query := &fbusering.ListShopCustomerWithFbExternalUserQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(request.Filters),
	}
	if err := s.FBUseringQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result.Customers = append(result.Customers, convertpbfabo.PbCustomersWithFbUsers(query.Result.Customers)...)
	result.Paging = cmapi.PbPageInfo(paging)
	result.Paging.Total = len(result.Customers)
	return result, nil
}
