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
	ID           dot.ID                    `json:"id"`
	FullName     string                    `json:"full_name"`
	FirstName    string                    `json:"first_name"`
	LastName     string                    `json:"last_name"`
	Phone        string                    `json:"phone"`
	Position     string                    `json:"position"`
	Email        string                    `json:"email"`
	Country      string                    `json:"country"`
	City         string                    `json:"city"`
	Province     string                    `json:"province"`
	District     string                    `json:"district"`
	Ward         string                    `json:"ward"`
	Zip          string                    `json:"zip"`
	DistrictCode string                    `json:"district_code"`
	ProvinceCode string                    `json:"province_code"`
	WardCode     string                    `json:"ward_code"`
	Company      string                    `json:"company"`
	Address1     string                    `json:"address_1"`
	Address2     string                    `json:"address_2"`
	IsDefault    bool                      `json:"is_default"`
	Type         addresstype.AddressType   `json:"type"`
	AccountID    dot.ID                    `json:"account_id"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	Coordinates  *orderv1types.Coordinates `json:"coordinates"`
	Notes        *AddressNote              `json:"notes"`
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
