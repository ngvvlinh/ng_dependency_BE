package model

import (
	"time"

	"etop.vn/api/shopping/customering/customer_type"
	"etop.vn/api/top/types/etc/gender"
	"etop.vn/capi/dot"
)

// +sqlgen
type ShopCustomer struct {
	ExternalID   string
	ExternalCode string
	PartnerID    dot.ID

	ID        dot.ID `paging:"id"`
	ShopID    dot.ID
	Code      string
	FullName  string
	Gender    gender.Gender              `sql_type:"enum(gender_type)"`
	Type      customer_type.CustomerType `sql_type:"enum(customer_type)"`
	Birthday  string                     `sql_type:"date"`
	Note      string
	Phone     string
	Email     string
	Status    int      `sql_type:"int2"`
	GroupIDs  []dot.ID `sq:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
