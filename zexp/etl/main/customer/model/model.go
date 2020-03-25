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
	Gender    gender.Gender
	Type      customer_type.CustomerType
	Birthday  string
	Note      string
	Phone     string
	Email     string
	Status    int
	GroupIDs  []dot.ID `sq:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
