package model

import (
	"time"

	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/types/etc/gender"
	"etop.vn/api/top/types/etc/status3"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopTrader(&ShopTrader{})

// +convert:type=tradering.ShopTrader
type ShopTrader struct {
	ID     dot.ID
	ShopID dot.ID
	Type   string
}

var _ = sqlgenShopCustomer(&ShopCustomer{})

// +convert:type=customering.ShopCustomer
type ShopCustomer struct {
	ExternalID   string
	ExternalCode string

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
}

var _ = sqlgenShopTraderAddress(&ShopTraderAddress{})

// +convert:type=addressing.ShopTraderAddress
type ShopTraderAddress struct {
	ID           dot.ID
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
	IsDefault    bool
	Coordinates  *ordermodel.Coordinates
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	DeletedAt    time.Time

	//Default status = 1
	Status status3.Status
}

var _ = sqlgenShopCustomerGroupCustomer(&ShopCustomerGroupCustomer{})

// +convert:type=customering.ShopCustomerGroupCustomer
type ShopCustomerGroupCustomer struct {
	GroupID    dot.ID
	CustomerID dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenShopCustomerGroup(&ShopCustomerGroup{})

// +convert:type=customering.ShopCustomerGroup
type ShopCustomerGroup struct {
	ID     dot.ID
	Name   string
	ShopID dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
