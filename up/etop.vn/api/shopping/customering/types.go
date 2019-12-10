package customering

import (
	"time"

	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/types/etc/gender"
	"etop.vn/api/top/types/etc/status3"
	dot "etop.vn/capi/dot"
)

type ShopCustomer struct {
	ID        dot.ID
	ShopID    dot.ID
	GroupIDs  []dot.ID
	Code      string
	CodeNorm  int
	FullName  string
	Gender    gender.Gender
	Type      customer_type.CustomerType
	Birthday  string
	Note      string
	Phone     string
	Email     string
	Status    status3.Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ShopCustomerGroup struct {
	ID   dot.ID
	Name string
}

type ShopCustomerGroupCustomer struct {
	GroupID    dot.ID
	CustomerID dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}
