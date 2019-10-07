package model

import (
	"time"

	"etop.vn/api/main/etop"
	ordermodel "etop.vn/backend/com/main/ordering/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopTrader(&ShopTrader{})

type ShopTrader struct {
	ID     int64
	ShopID int64
	Type   string
}

var _ = sqlgenShopCustomer(&ShopCustomer{})

type ShopCustomer struct {
	ID        int64
	ShopID    int64
	Code      string
	FullName  string
	Gender    string
	Type      string
	Birthday  string
	Note      string
	Phone     string
	Email     string
	Status    int32
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time
}

var _ = sqlgenShopTraderAddress(&ShopTraderAddress{})

type ShopTraderAddress struct {
	ID           int64
	ShopID       int64
	TraderID     int64
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
	Status etop.Status3
}
