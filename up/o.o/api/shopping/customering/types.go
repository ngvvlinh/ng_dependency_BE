package customering

import (
	"time"

	"o.o/api/shopping/customering/customer_type"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/status3"
	dot "o.o/capi/dot"
)

const CustomerAnonymous dot.ID = 1

type ShopCustomer struct {
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID

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

	Deleted bool
}

type ShopCustomerGroup struct {
	ID        dot.ID
	PartnerID dot.ID
	ShopID    dot.ID
	Name      string
	Deleted   bool
}

type ShopCustomerGroupCustomer struct {
	GroupID    dot.ID
	CustomerID dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}