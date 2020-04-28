package model

import (
	"time"

	"o.o/capi/dot"
)

// +sqlgen
type ShopCustomerGroupCustomer struct {
	GroupID    dot.ID
	CustomerID dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
