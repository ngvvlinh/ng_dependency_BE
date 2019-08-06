package customering

import (
	"time"

	orderingv1types "etop.vn/api/main/ordering/v1/types"
)

type ShopTrader struct {
	ID     int64
	ShopID int64
}

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
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ShopVendor struct {
	ID        int64
	ShopID    int64
	Name      string
	Note      string
	Status    int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
	Coordinates  *orderingv1types.Coordinates
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
