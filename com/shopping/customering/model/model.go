package model

import (
	"time"

	"o.o/api/top/types/etc/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/status3"
	addressmodel "o.o/backend/com/main/address/model"
	"o.o/capi/dot"
)

// +convert:type=tradering.ShopTrader
// +sqlgen
type ShopTrader struct {
	ID     dot.ID
	ShopID dot.ID
	Type   string
}

// +convert:type=customering.ShopCustomer
// +sqlgen
type ShopCustomer struct {
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID

	ID           dot.ID `paging:"id"`
	ShopID       dot.ID
	Code         string
	CodeNorm     int
	FullName     string
	Gender       gender.Gender
	Type         customer_type.CustomerType
	Birthday     string
	Note         string
	Phone        string
	Email        string
	Status       int
	FullNameNorm string
	PhoneNorm    string
	GroupIDs     []dot.ID  `sq:"-"`
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update" paging:"updated_at"`
	DeletedAt    time.Time

	Rid dot.ID
}

// +convert:type=addressing.ShopTraderAddress
// +sqlgen
type ShopTraderAddress struct {
	ID           dot.ID `paging:"id"`
	PartnerID    dot.ID
	ShopID       dot.ID
	TraderID     dot.ID
	FullName     string
	Phone        string
	Email        string
	Company      string
	Address1     string
	Address2     string
	DistrictCode string
	WardCode     string
	Position     string
	IsDefault    bool
	Coordinates  *addressmodel.Coordinates
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update" paging:"updated_at"`
	DeletedAt    time.Time

	//Default status = 1
	Status status3.Status

	Rid dot.ID
}

// +sqlgen
// +sqlgen:           ShopTraderAddress as sta
// +sqlgen:left-join: ShopTraderAddressSearch as stas on sta.id = stas.id
type ShopTraderAddressExtended struct {
	ShopTraderAddress       *ShopTraderAddress
	ShopTraderAddressSearch *ShopTraderAddressSearch
}

// +sqlgen
type ShopTraderAddressSearch struct {
	ID        dot.ID
	PhoneNorm string
}

// +convert:type=customering.ShopCustomerGroupCustomer
// +sqlgen
type ShopCustomerGroupCustomer struct {
	GroupID    dot.ID
	CustomerID dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	Rid dot.ID
}

// +convert:type=customering.ShopCustomerGroup
// +sqlgen
type ShopCustomerGroup struct {
	ID        dot.ID `paging:"id"`
	PartnerID dot.ID
	Name      string
	ShopID    dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update" paging:"updated_at"`
	DeletedAt time.Time

	Rid dot.ID
}
