package customering

import (
	"time"

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
	CodeNorm  int32
	FullName  string
	Gender    string
	Type      CustomerType
	Birthday  string
	Note      string
	Phone     string
	Email     string
	Status    int32
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
