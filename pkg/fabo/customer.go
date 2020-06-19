package fabo

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/fabo"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	convertpbfabo "o.o/backend/pkg/fabo/convertpb"
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

func (s *CustomerService) ListCustomersWithFbUsers(ctx context.Context, request *fabo.ListCustomersWithFbUsersRequest) (*fabo.ListCustomersWithFbUsersResponse, error) {
	var result = &fabo.ListCustomersWithFbUsersResponse{}
	paging := cmapi.CMPaging(request.Paging)
	shopID := s.SS.Shop().ID
	if paging.Offset == 0 && request.GetAll && len(request.Filters) == 0 {
		paging.Limit = paging.Limit - 1
		queryCustomerIndenpendent := &customering.GetCustomerIndependentQuery{}
		if err := s.CustomerQuery.Dispatch(ctx, queryCustomerIndenpendent); err != nil {
			return nil, err
		}
		result.Customers = append(result.Customers,
			convertpbfabo.PbCustomerWithFbUser(
				&fbusering.ShopCustomerWithFbExternalUser{
					ShopCustomer: queryCustomerIndenpendent.Result,
				},
			))
	}
	if paging.Limit > 0 {
		query := &fbusering.ListShopCustomerWithFbExternalUserQuery{
			ShopID:  shopID,
			Paging:  *paging,
			Filters: cmapi.ToFilters(request.Filters),
		}
		if err := s.FBUseringQuery.Dispatch(ctx, query); err != nil {
			return nil, err
		}
		result.Customers = append(result.Customers, convertpbfabo.PbCustomersWithFbUsers(query.Result.Customers)...)
		result.Paging = cmapi.PbPageInfo(paging)
		result.Paging.Total = len(result.Customers)
		// Nếu có chứa customer anonymous
		if len(result.Customers) > 0 && result.Customers[0].Id == 1 {
			result.Paging.Limit++
		}
	}
	return result, nil
}
