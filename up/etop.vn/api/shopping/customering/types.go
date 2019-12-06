package customering

import (
	"time"

	"etop.vn/api/top/types/etc/gender"
	dot "etop.vn/capi/dot"
)

type CustomerType string

const (
	CustomerTypeIndividual   CustomerType = "individual"
	CustomerTypeOrganization CustomerType = "organization"
	CustomerTypeIndependent  CustomerType = "independent"
)

type ShopCustomer struct {
	ID        dot.ID
	ShopID    dot.ID
	GroupIDs  []dot.ID
	Code      string
	CodeNorm  int
	FullName  string
	Gender    gender.Gender
	Type      CustomerType
	Birthday  string
	Note      string
	Phone     string
	Email     string
	Status    int
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
