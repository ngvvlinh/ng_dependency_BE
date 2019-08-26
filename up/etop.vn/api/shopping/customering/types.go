package customering

import (
	"time"
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
