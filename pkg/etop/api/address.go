package api

import (
	"context"

	apietop "o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	addressmodelx "o.o/backend/com/main/address/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
)

type AddressService struct {
}

func (s *AddressService) Clone() *AddressService {
	res := *s
	return &res
}

func (s *AddressService) CreateAddress(ctx context.Context, q *CreateAddressEndpoint) error {
	address, err := convertpb.PbCreateAddressToModel(q.Context.AccountID, q.CreateAddressRequest)
	if err != nil {
		return err
	}
	cmd := &addressmodelx.CreateAddressCommand{
		Address: address,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbAddress(cmd.Result)
	return nil
}

func (s *AddressService) GetAddresses(ctx context.Context, q *GetAddressesEndpoint) error {
	accountID := q.Context.AccountID
	query := &addressmodelx.GetAddressesQuery{
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil
	}
	q.Result = &apietop.GetAddressResponse{
		Addresses: convertpb.PbAddresses(query.Result.Addresses),
	}
	return nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, q *UpdateAddressEndpoint) error {
	accountID := q.Context.AccountID
	address, err := convertpb.PbUpdateAddressToModel(accountID, q.UpdateAddressRequest)
	if err != nil {
		return err
	}
	cmd := &addressmodelx.UpdateAddressCommand{
		Address: address,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbAddress(cmd.Result)
	return nil
}

func (s *AddressService) RemoveAddress(ctx context.Context, q *RemoveAddressEndpoint) error {
	accountID := q.Context.AccountID
	cmd := &addressmodelx.DeleteAddressCommand{
		ID:        q.Id,
		AccountID: accountID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.Empty{}
	return nil
}
