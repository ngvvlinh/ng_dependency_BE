package shop

import (
	"context"

	"etop.vn/api/main/location"
	pbcm "etop.vn/api/pb/common"
	pbshop "etop.vn/api/pb/etop/shop"
	"etop.vn/api/shopping/addressing"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		customerService.CreateCustomerAddress,
		customerService.DeleteCustomerAddress,
		customerService.GetCustomerAddresses,
		customerService.UpdateCustomerAddress,
		customerService.SetDefaultCustomerAddress,
	)
}

func (s *CustomerService) CreateCustomerAddress(ctx context.Context, r *CreateCustomerAddressEndpoint) error {
	query := &location.GetLocationQuery{
		DistrictCode: r.DistrictCode,
		WardCode:     r.WardCode,
	}
	if err := locationQuery.Dispatch(ctx, query); err != nil {
		return err
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
		Coordinates:  convertpb.PbCoordinatesToModel(r.Coordinates),
		IsDefault:    true,
	}
	if err := traderAddressAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	pbAddr, err := convertpb.PbShopAddress(ctx, cmd.Result, locationQuery)
	if err != nil {
		return err
	}
	r.Result = pbAddr
	return nil
}

func (s *CustomerService) DeleteCustomerAddress(ctx context.Context, r *DeleteCustomerAddressEndpoint) error {
	cmd := &addressing.DeleteAddressCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := traderAddressAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: int(cmd.Result)}
	return nil
}

func (s *CustomerService) GetCustomerAddresses(ctx context.Context, r *GetCustomerAddressesEndpoint) error {
	query := &addressing.ListAddressesByTraderIDQuery{
		ShopID:   r.Context.Shop.ID,
		TraderID: r.CustomerId,
	}
	if err := traderAddressQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	addrs, err := convertpb.PbShopAddresses(ctx, query.Result, locationQuery)
	if err != nil {
		return err
	}
	r.Result = &pbshop.CustomerAddressesResponse{Addresses: addrs}
	return nil
}

func (s *CustomerService) UpdateCustomerAddress(ctx context.Context, r *UpdateCustomerAddressEndpoint) error {
	if r.DistrictCode.Valid && r.WardCode.Valid {
		query := &location.GetLocationQuery{
			DistrictCode: r.DistrictCode.String,
			WardCode:     r.WardCode.String,
		}
		if err := locationQuery.Dispatch(ctx, query); err != nil {
			return err
		}
	}

	// TODO: verify district & ward code
	cmd := &addressing.UpdateAddressCommand{
		ID:           r.Id,
		ShopID:       r.Context.Shop.ID,
		FullName:     PString(r.FullName),
		Phone:        PString(r.Phone),
		Email:        PString(r.Email),
		Company:      PString(r.Company),
		Address1:     PString(r.Address1),
		Address2:     PString(r.Address2),
		DistrictCode: PString(r.DistrictCode),
		WardCode:     PString(r.WardCode),
		Coordinates:  convertpb.PbCoordinatesToModel(r.Coordinates),
	}
	if err := traderAddressAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	addr, err := convertpb.PbShopAddress(ctx, cmd.Result, locationQuery)
	if err != nil {
		return err
	}
	r.Result = addr
	return nil
}

func (s *CustomerService) SetDefaultCustomerAddress(ctx context.Context, r *SetDefaultCustomerAddressEndpoint) error {
	query := &addressing.GetAddressByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := traderAddressQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrap(cm.NotFound, "traderAddress not found").
			Throw()
	}

	setDefaultAddressCmd := &addressing.SetDefaultAddressCommand{
		ID:       r.Id,
		TraderID: query.Result.TraderID,
		ShopID:   r.Context.Shop.ID,
	}
	if err := traderAddressAggr.Dispatch(ctx, setDefaultAddressCmd); err != nil {
		return nil
	}
	r.Result = &pbcm.UpdatedResponse{Updated: setDefaultAddressCmd.Result.Updated}
	return nil
}
