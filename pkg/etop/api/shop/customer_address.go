package shop

import (
	"context"

	"etop.vn/api/shopping/addressing"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbshop "etop.vn/backend/pb/etop/shop"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("api",
		CreateCustomerAddress,
		DeleteCustomerAddress,
		GetCustomerAddresses,
		UpdateCustomerAddress,
	)
}

func CreateCustomerAddress(ctx context.Context, r *wrapshop.CreateCustomerAddressEndpoint) error {
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
		Coordinates:  pbetop.PbCoordinatesToModel(r.Coordinates),
	}
	if err := traderAddressAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	pbAddr, err := pbshop.PbShopAddress(ctx, cmd.Result, locationQuery)
	if err != nil {
		return err
	}
	r.Result = pbAddr
	return nil
}

func DeleteCustomerAddress(ctx context.Context, r *wrapshop.DeleteCustomerAddressEndpoint) error {
	cmd := &addressing.DeleteAddressCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := traderAddressAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}

func GetCustomerAddresses(ctx context.Context, r *wrapshop.GetCustomerAddressesEndpoint) error {
	query := &addressing.ListAddressesByTraderIDQuery{
		ShopID:   r.Context.Shop.ID,
		TraderID: r.CustomerId,
	}
	if err := traderAddressQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	addrs, err := pbshop.PbShopAddresses(ctx, query.Result, locationQuery)
	if err != nil {
		return err
	}
	r.Result = &pbshop.CustomerAddressesResponse{Addresses: addrs}
	return nil
}

func UpdateCustomerAddress(ctx context.Context, r *wrapshop.UpdateCustomerAddressEndpoint) error {
	// TODO: verify district & ward code
	cmd := &addressing.UpdateAddressCommand{
		ID:           r.Id,
		ShopID:       r.Context.Shop.ID,
		FullName:     String(r.FullName),
		Phone:        String(r.Phone),
		Email:        String(r.Email),
		Company:      String(r.Company),
		Address1:     String(r.Address1),
		Address2:     String(r.Address2),
		DistrictCode: String(r.DistrictCode),
		WardCode:     String(r.WardCode),
		Coordinates:  pbetop.PbCoordinatesToModel(r.Coordinates),
	}
	if err := traderAddressAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	addr, err := pbshop.PbShopAddress(ctx, cmd.Result, locationQuery)
	if err != nil {
		return err
	}
	r.Result = addr
	return nil
}
