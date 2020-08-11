package api

import (
	"context"

	api "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	addressmodelx "o.o/backend/com/main/address/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type AddressService struct {
	session.Session
}

func (s *AddressService) Clone() api.AddressService {
	res := *s
	return &res
}

func (s *AddressService) CreateAddress(ctx context.Context, q *api.CreateAddressRequest) (*api.Address, error) {
	if _, ok := validate.NormalizePhone(q.Phone); !ok {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
	}

	address, err := convertpb.PbCreateAddressToModel(s.SS.Claim().AccountID, q)
	if err != nil {
		return nil, err
	}
	cmd := &addressmodelx.CreateAddressCommand{
		Address: address,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbAddress(cmd.Result)
	return result, nil
}

func (s *AddressService) GetAddresses(ctx context.Context, q *pbcm.Empty) (*api.GetAddressResponse, error) {
	accountID := s.SS.Claim().AccountID
	query := &addressmodelx.GetAddressesQuery{
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetAddressResponse{
		Addresses: convertpb.PbAddresses(query.Result.Addresses),
	}
	return result, nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, q *api.UpdateAddressRequest) (*api.Address, error) {
	accountID := s.SS.Claim().AccountID
	address, err := convertpb.PbUpdateAddressToModel(accountID, q)
	if err != nil {
		return nil, err
	}
	cmd := &addressmodelx.UpdateAddressCommand{
		Address: address,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbAddress(cmd.Result)
	return result, nil
}

func (s *AddressService) RemoveAddress(ctx context.Context, q *pbcm.IDRequest) (*pbcm.Empty, error) {
	accountID := s.SS.Claim().AccountID
	cmd := &addressmodelx.DeleteAddressCommand{
		ID:        q.Id,
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.Empty{}
	return result, nil
}
