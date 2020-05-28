package fabo

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/shopping/customering"
	"o.o/api/top/int/fabo"
	"o.o/backend/pkg/etop/authorize/session"
	convertpbfabo "o.o/backend/pkg/fabo/convertpb"
)

type CustomerService struct {
	session.Session

	customerQuery  customering.QueryBus
	fbUseringQuery fbusering.QueryBus
	fbUseringAggr  fbusering.CommandBus
}

func NewCustomerService(
	customerQ customering.QueryBus,
	fbUseringQ fbusering.QueryBus,
	fbUseringA fbusering.CommandBus,
	ssParam session.Session,
) *CustomerService {
	s := &CustomerService{
		customerQuery:  customerQ,
		fbUseringQuery: fbUseringQ,
		fbUseringAggr:  fbUseringA,
		Session:        ssParam,
	}
	return s
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
	err := s.fbUseringAggr.Dispatch(ctx, cmd)
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
	if err := s.fbUseringQuery.Dispatch(ctx, query); err != nil {
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
	if err := s.fbUseringQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	var result = &fabo.ListFbUsersResponse{}
	result.FbUsers = convertpbfabo.PbFbUsers(query.Result)
	return result, nil
}
