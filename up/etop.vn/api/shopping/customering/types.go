package customering

import (
	"time"
)

var (
	CustomerType = "customer"
	VendorType   = "vendor"
	CarrierType  = "carrier"
)

type ShopTrader struct {
	ID     int64
	ShopID int64
	Type   string
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
