package customer

import (
	"context"

	"o.o/api/main/location"
	"o.o/api/shopping/addressing"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
)

func (s *CustomerService) CreateCustomerAddress(ctx context.Context, r *api.CreateCustomerAddressRequest) (*api.CustomerAddress, error) {
	query := &location.GetLocationQuery{
		DistrictCode: r.DistrictCode,
		WardCode:     r.WardCode,
	}
	if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	cmd := &addressing.CreateAddressCommand{
		ShopID:       s.SS.Shop().ID,
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
		Coordinates:  convertpb.PbCoordinatesToModel(r.Coordinates),
		IsDefault:    true,
	}
	if err := s.AddressAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	pbAddr, err := convertpb.PbShopAddress(ctx, cmd.Result, s.LocationQuery)
	if err != nil {
		return nil, err
	}
	result := pbAddr
	return result, nil
}

func (s *CustomerService) DeleteCustomerAddress(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &addressing.DeleteAddressCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.AddressAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: cmd.Result}
	return result, nil
}

func (s *CustomerService) GetCustomerAddresses(ctx context.Context, r *api.GetCustomerAddressesRequest) (*api.CustomerAddressesResponse, error) {
	var phone = ""
	if r.Filter != nil {
		phone = r.Filter.Phone
	}
	query := &addressing.ListAddressesQuery{
		ShopID:   s.SS.Shop().ID,
		TraderID: r.CustomerId,
		Phone:    phone,
	}
	if err := s.AddressQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	addrs, err := convertpb.PbShopAddresses(ctx, query.Result.ShopTraderAddresses, s.LocationQuery)
	if err != nil {
		return nil, err
	}
	result := &api.CustomerAddressesResponse{Addresses: addrs}
	return result, nil
}

func (s *CustomerService) UpdateCustomerAddress(ctx context.Context, r *api.UpdateCustomerAddressRequest) (*api.CustomerAddress, error) {
	if r.DistrictCode.Valid && r.WardCode.Valid {
		query := &location.GetLocationQuery{
			DistrictCode: r.DistrictCode.String,
			WardCode:     r.WardCode.String,
		}
		if err := s.LocationQuery.Dispatch(ctx, query); err != nil {
			return nil, err
		}
	}

	// TODO: verify district & ward code
	cmd := &addressing.UpdateAddressCommand{
		ID:           r.Id,
		ShopID:       s.SS.Shop().ID,
		FullName:     r.FullName,
		Phone:        r.Phone,
		Email:        r.Email,
		Company:      r.Company,
		Address1:     r.Address1,
		Address2:     r.Address2,
		DistrictCode: r.DistrictCode,
		WardCode:     r.WardCode,
		Position:     r.Position,
		Coordinates:  convertpb.PbCoordinatesToModel(r.Coordinates),
	}
	if err := s.AddressAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	addr, err := convertpb.PbShopAddress(ctx, cmd.Result, s.LocationQuery)
	if err != nil {
		return nil, err
	}
	result := addr
	return result, nil
}

func (s *CustomerService) SetDefaultCustomerAddress(ctx context.Context, r *pbcm.IDRequest) (*pbcm.UpdatedResponse, error) {
	query := &addressing.GetAddressByIDQuery{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.AddressQuery.Dispatch(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "traderAddress not found").
			Throw()
	}

	setDefaultAddressCmd := &addressing.SetDefaultAddressCommand{
		ID:       r.Id,
		TraderID: query.Result.TraderID,
		ShopID:   s.SS.Shop().ID,
	}
	if err := s.AddressAggr.Dispatch(ctx, setDefaultAddressCmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: setDefaultAddressCmd.Result.Updated}
	return result, nil
}
