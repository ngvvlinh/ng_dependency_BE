package customering

import (
	"time"
)

type ShopCustomer struct {
	ID        int64
	ShopID    int64
	GroupIDs  []int64
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

type ShopCustomerGroup struct {
	ID   int64
	Name string
}

type ShopCustomerGroupCustomer struct {
	GroupID    int64
	CustomerID int64

	CreatedAt time.Time
	UpdatedAt time.Time
}
