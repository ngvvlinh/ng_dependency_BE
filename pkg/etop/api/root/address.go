package root

import (
	"context"

	address "o.o/api/main/address"
	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type AddressService struct {
	session.Session

	AddressAggr address.CommandBus
	AddressQS   address.QueryBus
}

func (s *AddressService) Clone() api.AddressService {
	res := *s
	return &res
}

func (s *AddressService) CreateAddress(ctx context.Context, q *api.CreateAddressRequest) (*api.Address, error) {
	cmd := convertpb.SetCreateAddressArgs(q, s.SS.Claim().AccountID)

	if err := s.AddressAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpb.Convert_core_Address_To_api_Address(cmd.Result), nil
}

func (s *AddressService) GetAddresses(ctx context.Context, q *pbcm.Empty) (*api.GetAddressResponse, error) {
	cmd := &address.ListAddressesQuery{
		ID: s.SS.Claim().AccountID,
	}
	if err := s.AddressQS.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	if cmd.Result != nil && len(cmd.Result.Addresses) > 0 {
		return &api.GetAddressResponse{
			Addresses: convertpb.Convert_core_Addresses_To_api_Addresses(cmd.Result.Addresses),
		}, nil
	}

	return &api.GetAddressResponse{
		Addresses: []*api.Address{},
	}, nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, q *api.UpdateAddressRequest) (*api.Address, error) {
	cmd := convertpb.SetUpdateAddressArgs(q, s.SS.Claim().AccountID)

	if err := s.AddressAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpb.Convert_core_Address_To_api_Address(cmd.Result), nil
}

func (s *AddressService) RemoveAddress(ctx context.Context, q *pbcm.IDRequest) (*pbcm.Empty, error) {
	cmd := &address.RemoveAddressCommand{
		AccountID: s.SS.Claim().AccountID,
		ID:        q.Id,
	}

	if err := s.AddressAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.Empty{}
	return result, nil
}
