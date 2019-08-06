package shop

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
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
	return cm.ErrTODO
}

func DeleteCustomerAddress(ctx context.Context, r *wrapshop.DeleteCustomerAddressEndpoint) error {
	return cm.ErrTODO
}

func GetCustomerAddresses(ctx context.Context, r *wrapshop.GetCustomerAddressesEndpoint) error {
	return cm.ErrTODO
}

func UpdateCustomerAddress(ctx context.Context, r *wrapshop.UpdateCustomerAddressEndpoint) error {
	return cm.ErrTODO
}
