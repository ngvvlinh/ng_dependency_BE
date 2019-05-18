package pm

import (
	"context"

	"etop.vn/api/main/address"
)

type ProcessManager struct {
	addressQuery address.QueryService
}

func NewProcessManager(addressQuery address.QueryService) *ProcessManager {
	return &ProcessManager{
		addressQuery: addressQuery,
	}
}

func (pm *ProcessManager) GetAddressByID(ctx context.Context, args *address.GetAddressByIDQueryArgs) (*address.Address, error) {
	return pm.addressQuery.GetAddressByID(ctx, args)
}
