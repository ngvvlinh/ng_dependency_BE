package address

import (
	"strings"
	"time"

	orderv1types "o.o/api/main/ordering/types"
	addresstype "o.o/api/top/types/etc/address_type"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

// +gen:event:topic=event/defaultaddress

type Address struct {
	ID           dot.ID
	FullName     string
	FirstName    string
	LastName     string
	Phone        string
	Position     string
	Email        string
	Country      string
	City         string
	Province     string
	District     string
	Ward         string
	Zip          string
	DistrictCode string
	ProvinceCode string
	WardCode     string
	Company      string
	Address1     string
	Address2     string
	IsDefault    bool
	Type         addresstype.AddressType
	AccountID    dot.ID
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	Coordinates  *orderv1types.Coordinates
	Notes        *AddressNote
}

func (m *Address) String() string { return jsonx.MustMarshalToString(m) }

type AddressNote struct {
	Note       string `json:"note"`
	OpenTime   string `json:"open_time"`
	LunchBreak string `json:"lunch_break"`
	Other      string `json:"other"`
}

func (m *AddressNote) String() string { return jsonx.MustMarshalToString(m) }

func (a *Address) ToOrderAddress() *orderv1types.Address {
	if a == nil {
		return nil
	}
	return &orderv1types.Address{
		FullName: a.FullName,
		Phone:    a.Phone,
		Email:    a.Email,
		Company:  a.Company,
		Address1: a.Address1,
		Address2: a.Address2,
		Location: orderv1types.Location{
			ProvinceCode: a.ProvinceCode,
			DistrictCode: a.DistrictCode,
			WardCode:     a.WardCode,
			Coordinates:  a.Coordinates,
		},
	}
}

func (a *Address) GetFullAddress() string {
	b := strings.Builder{}
	if a.Address1 != "" {
		b.WriteString(a.Address1)
		b.WriteByte('\n')
	}
	if a.Address2 != "" {
		b.WriteString(a.Address2)
		b.WriteByte('\n')
	}
	flag := false
	if a.Ward != "" {
		b.WriteString(a.Ward)
		flag = true
	}
	if a.District != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(a.District)
		flag = true
	}
	if a.Province != "" {
		if flag {
			b.WriteString(", ")
		}
		b.WriteString(a.Province)
	}
	return b.String()
}

type AddressCreatedEvent struct {
	ID        dot.ID
	AccountID dot.ID
	Type      addresstype.AddressType
}

type AddressDefaultUpdatedEvent struct {
	ID                dot.ID
	ShipFromAddressID dot.ID
}
