package model

type GetAddressQuery struct {
	AddressID int64

	Result *Address
}

type GetAddressesQuery struct {
	AccountID int64

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
	ID        int64
	AccountID int64
}
