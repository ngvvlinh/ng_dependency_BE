package model

import "etop.vn/capi/dot"

type GetAddressQuery struct {
	AddressID dot.ID

	Result *Address
}

type GetAddressesQuery struct {
	AccountID dot.ID

	Result struct {
		Addresses []*Address
	}
}

type CreateAddressCommand struct {
	Address *Address
	Result  *Address
}

type UpdateAddressCommand struct {
	Address *Address
	Result  *Address
}

type DeleteAddressCommand struct {
	ID        dot.ID
	AccountID dot.ID
}
