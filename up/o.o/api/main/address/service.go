package address

import (
	"context"

	orderv1types "o.o/api/main/ordering/types"
	addresstype "o.o/api/top/types/etc/address_type"
	dot "o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateAddress(context.Context, *CreateAddressArgs) (*Address, error)
	UpdateAddress(context.Context, *UpdateAddressArgs) (*Address, error)
	RemoveAddress(context.Context, *DeleteAddressArgs) error
}

type QueryService interface {
	GetAddressByID(context.Context, *GetAddressByIDQueryArgs) (*Address, error)
	ListAddresses(ctx context.Context, ID dot.ID) (*GetAddressResponse, error)
}

type GetAddressByIDQueryArgs struct {
	ID dot.ID
}

type DeleteAddressArgs struct {
	ID        dot.ID
	AccountID dot.ID
}

func (m *DeleteAddressArgs) Validate() error {
	if m.ID == 0 {
		return EditErrorMsg("ID")
	}
	if m.AccountID == 0 {
		return EditErrorMsg("AccountID")
	}

	return nil
}

type GetAddressResponse struct {
	Addresses []*Address
}

// +convert:create=Address
type CreateAddressArgs struct {
	Province     string
	ProvinceCode string
	District     string
	DistrictCode string
	Ward         string
	WardCode     string
	Address1     string
	Address2     string
	Zip          string
	Company      string
	City         string
	Country      string
	FullName     string
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	Position     string
	Type         addresstype.AddressType
	AccountID    dot.ID
	Notes        *AddressNote
	Coordinates  *orderv1types.Coordinates
}

func (m *CreateAddressArgs) Validate() error {
	if m.FullName == "" {
		return EditErrorMsg("Tên")
	}
	if m.Address1 == "" {
		return EditErrorMsg("Địa chỉ cụ thể")
	}
	if m.Phone == "" {
		return EditErrorMsg("Số điện thoại")
	}
	if m.District == "" || m.DistrictCode == "" {
		return EditErrorMsg("Quận/Huyện")
	}
	if m.Ward == "" || m.WardCode == "" {
		return EditErrorMsg("Phường/Xã")
	}
	if m.Province == "" || m.ProvinceCode == "" {
		return EditErrorMsg("Tỉnh")
	}
	if m.AccountID == 0 {
		return EditErrorMsg("Tên")
	}
	return nil
}

// +convert:update=Address
type UpdateAddressArgs struct {
	ID           dot.ID
	Province     string
	ProvinceCode string
	District     string
	DistrictCode string
	Ward         string
	WardCode     string
	Address1     string
	Address2     string
	Zip          string
	Company      string
	City         string
	Country      string
	FullName     string
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	Position     string
	Type         addresstype.AddressType
	AccountID    dot.ID
	Notes        *AddressNote
	Coordinates  *orderv1types.Coordinates
}

func (m *UpdateAddressArgs) Validate() (isUpdateLocation int, err error) {
	if m.ID == 0 {
		return 0, EditErrorMsg("ID")
	}

	if m.District != "" || m.DistrictCode != "" || m.Ward != "" || m.WardCode != "" || m.Address1 != "" || m.ProvinceCode != "" || m.Province != "" {
		if m.District == "" || m.DistrictCode == "" {
			return 0, EditErrorMsg("Quận/Huyện")
		}
		if m.Ward == "" || m.WardCode == "" {
			return 0, EditErrorMsg("Phường/Xã")
		}
		if m.Address1 == "" {
			return 0, EditErrorMsg("Địa chỉ cụ thể")
		}

		if m.Province == "" || m.ProvinceCode == "" {
			return 0, EditErrorMsg("Tỉnh")
		}
		return 1, nil
	}

	return 0, nil
}

func EditErrorMsg(str string) error {
	return xerrors.Errorf(xerrors.InvalidArgument, nil, "Vui lòng nhập thông tin bắt buộc, thiếu %v", str)
}
