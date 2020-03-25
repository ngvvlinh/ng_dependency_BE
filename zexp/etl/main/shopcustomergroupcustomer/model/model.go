package model

import (
	"time"

	"etop.vn/capi/dot"
)

// +sqlgen
type ShopCustomerGroupCustomer struct {
	GroupID    dot.ID
	CustomerID dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time

	Rid dot.ID
}
