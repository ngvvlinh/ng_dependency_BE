package modelx

import (
	"o.o/backend/com/main/address/model"
	"o.o/capi/dot"
)

type GetAddressQuery struct {
	AddressID dot.ID

	Result *model.Address
}

type GetAddressesQuery struct {
	AccountID dot.ID

	Result struct {
		Addresses []*model.Address
	}
}

type CreateAddressCommand struct {
	Address *model.Address
	Result  *model.Address
}

type UpdateAddressCommand struct {
	Address *model.Address
	Result  *model.Address
}

type DeleteAddressCommand struct {
	ID        dot.ID
	AccountID dot.ID
}
